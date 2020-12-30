package huarpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go/token"
	"net/http"
	"reflect"
	"strings"
)

var (
	errType     = reflect.TypeOf(struct{ error }{}).Field(0).Type
	contextType = reflect.TypeOf(struct{ context.Context }{}).Field(0).Type
)

type MethodTags struct {
	HTTP struct {
		Method string
		Path   string
	}
	Help string
}

type ServiceMethod struct {
	Belong *Service

	Name     string
	ArgType  reflect.Type
	Callable reflect.Value
	Tags     MethodTags

	FormDecoder FormDecoder
	FormEncoder FormEncoder
}

func parseServiceMethod(field reflect.StructField, value reflect.Value) (*ServiceMethod, error) {
	if !token.IsExported(field.Name) {
		return nil, errors.New("service的所有属性必须是可导出的")
	}
	sft := field.Type
	if sft.NumIn() != 2 || sft.NumOut() != 2 || sft.In(0) != contextType || sft.Out(1) != errType {
		return nil, fmt.Errorf("函数%s签名错误（正确: func(context.Context,*xxArg)(*xxReply,error),xxArg为入参，xxReply为出参)", field.Name)
	}

	if !strings.HasSuffix(sft.In(1).Elem().Name(), "Arg") {
		return nil, fmt.Errorf("%s must end with 'Arg'", sft.In(1).Elem().Name())
	}
	if !strings.HasSuffix(sft.Out(0).Elem().Name(), "Reply") {
		return nil, fmt.Errorf("%s must end with 'Reply'", sft.Out(0).Elem().Name())
	}

	sm := &ServiceMethod{
		Name:        field.Name,
		ArgType:     sft.In(1),
		Callable:    value,
		FormDecoder: newFormDecoder(),
		FormEncoder: newFromEncoder(),
		Tags: MethodTags{
			Help: field.Tag.Get("help"),
		},
	}
	if httpTag, ok := field.Tag.Lookup("http"); ok {
		httpTag := strings.SplitN(httpTag, " ", 2)
		sm.Tags.HTTP.Method = httpTag[0]
		sm.Tags.HTTP.Path = httpTag[1]
	}

	return sm, nil
}

func (s *ServiceMethod) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	arg := reflect.New(s.ArgType.Elem())
	contentType := r.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "application/json") {
		if err := json.NewDecoder(r.Body).Decode(arg.Interface()); err != nil {
			http.Error(w, fmt.Sprintf("json decode error: %s", err), http.StatusBadRequest)
			return
		}
	} else {
		if err := r.ParseForm(); err != nil {
			http.Error(w, fmt.Sprintf("parse form error: %s", err), http.StatusBadRequest)
			return
		}
		if err := s.FormDecoder.Decode(arg.Interface(), r.Form); err != nil {
			http.Error(w, fmt.Sprintf("decode form error: %s", err), http.StatusBadRequest)
			return
		}
	}

	// TODO: 优化性能
	n := arg.Type().Elem()
	ve := arg.Elem()
	for i := 0; i < n.NumField(); i++ {
		key := n.Field(i).Name
		if tag := n.Field(i).Tag.Get("path"); tag != "" {
			key = tag
		}
		val := s.Belong.router.URLParams(r, key)
		if val != "" {
			ve.Field(i).Set(reflect.ValueOf(val).Convert(ve.Field(i).Type()))
		}
	}

	reply := s.Callable.Call([]reflect.Value{reflect.ValueOf(r.Context()), arg})
	err, _ := reply[1].Interface().(error)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(reply[0].Interface()); err != nil {
		http.Error(w, fmt.Sprintf("返回数据失败: %s", err), http.StatusInternalServerError)
		return
	}
}

type Service struct {
	Name    string
	Methods []*ServiceMethod

	options *options
	router  Router
}

// route 注册路由
func (s *Service) route(r Router) {
	mux := r
	if s.options != nil && len(s.options.httpMiddlewares) > 0 {
		mux = r.With(s.options.httpMiddlewares...)
	}
	s.router = mux

	var restPattern string
	for i := range s.Methods {
		m := s.Methods[i]
		methodName := m.Tags.HTTP.Path
		if methodName == "" {
			methodName = "/" + strings.ToLower(m.Name)
		}

		restPattern = "/" + strings.ToLower(s.Name) + methodName
		if m.Tags.HTTP.Method == "" {
			mux.Handle(restPattern, m)
		} else {
			mux.Method(m.Tags.HTTP.Method, restPattern, m)
		}
	}
}

func (s *Service) With(opts ...Option) *Service {
	if s.options == nil {
		s.options = newOptions()
	}
	for i := range opts {
		opts[i].apply(s.options)
	}
	return s
}

func (s *Service) Endpoint() (pattern string, handler http.Handler) {
	mux := NewServeMux()
	s.route(mux)
	return fmt.Sprintf("/%s/", strings.ToLower(s.Name)), mux
}

// NewService 创建服务
func NewService(s interface{}) *Service {
	v := reflect.Indirect(reflect.ValueOf(s))
	t := v.Type()

	if !strings.HasSuffix(t.Name(), "Service") {
		panic(errors.New("service name必须包含Service后缀"))
	}

	srv := &Service{}
	srv.Name = strings.TrimSuffix(t.Name(), "Service")

	for i := 0; i < t.NumField(); i++ {
		sm, err := parseServiceMethod(t.Field(i), v.Field(i))
		if err != nil {
			panic(fmt.Errorf("解析方法出错: %w", err))
		}
		sm.Belong = srv
		srv.Methods = append(srv.Methods, sm)
	}

	return srv
}
