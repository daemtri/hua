package huarpc

import (
	"context"
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
}

func parseServiceMethod(field reflect.StructField, value reflect.Value) (*ServiceMethod, error) {
	if !token.IsExported(field.Name) {
		return nil, errors.New("service的所有属性必须是可导出的")
	}
	sft := field.Type

	// TODO: 允许没有请求参数
	if sft.NumIn() != 2 || sft.NumOut() != 2 || sft.In(0) != contextType || sft.Out(1) != errType {
		return nil, fmt.Errorf("函数%s签名错误（正确: func(context.Context,*xxArg)(*xxReply,error),xxArg为入参，xxReply为出参)", field.Name)
	}

	if !strings.HasSuffix(sft.In(1).Elem().Name(), "Arg") {
		return nil, fmt.Errorf("%s must end with 'Arg'", sft.In(1).Elem().Name())
	}

	// TODO: 限制返回必须是结构体或者golang 标准变量类型
	if sft.Out(0).Kind() == reflect.Chan {
		if !strings.HasSuffix(sft.Out(0).Elem().Elem().Name(), "Reply") {
			return nil, fmt.Errorf("%s must end with 'Reply'", sft.Out(0).Elem().Elem().Name())
		}
	} else if !strings.HasSuffix(sft.Out(0).Elem().Name(), "Reply") {
		return nil, fmt.Errorf("%s must end with 'Reply'", sft.Out(0).Elem().Name())
	}

	sm := &ServiceMethod{
		Name:     field.Name,
		ArgType:  sft.In(1),
		Callable: value,
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
	argInterface := arg.Interface()

	defer func() {
		_ = r.Body.Close()
	}()
	if strings.HasPrefix(contentType, "application/json") {
		if err := json.NewDecoder(r.Body).Decode(argInterface); err != nil {
			http.Error(w, fmt.Sprintf("json decode error: %s", err), http.StatusBadRequest)
			return
		}
	} else {
		if err := r.ParseForm(); err != nil {
			http.Error(w, fmt.Sprintf("parse form error: %s", err), http.StatusBadRequest)
			return
		}
		if err := form.Decoder.Decode(argInterface, r.Form); err != nil {
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

	// validate
	if s.Belong.validator != nil {
		if err := s.Belong.validator.ValidateStruct(argInterface); err != nil {
			http.Error(w, fmt.Sprintf("validate: %s", err), http.StatusBadRequest)
			return
		}
	}

	reply := s.Callable.Call([]reflect.Value{reflect.ValueOf(r.Context()), arg})
	err, _ := reply[1].Interface().(error)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// chan
	if reply[0].Kind() == reflect.Chan {
		// Set the headers related to event streaming.
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Transfer-Encoding", "chunked")
		flusher := w.(http.Flusher)
		for {
			v, ok := reply[0].Recv()
			if !ok {
				break
			}

			_, _ = fmt.Fprintf(w, "id: %s\n", v.Elem().FieldByName("ID").Interface())
			_, _ = fmt.Fprintf(w, "event: %s\n", v.Elem().Type().Name())
			_, _ = fmt.Fprintf(w, "retry: %d\n", 100)
			_, _ = fmt.Fprint(w, "data: ")
			_ = json.NewEncoder(w).Encode(v.Interface())
			_, _ = fmt.Fprint(w, "\n")
			flusher.Flush()
		}
	} else if err := json.NewEncoder(w).Encode(reply[0].Interface()); err != nil {
		http.Error(w, fmt.Sprintf("返回数据失败: %s", err), http.StatusInternalServerError)
		return
	}
}

type Service struct {
	Name    string
	Methods []*ServiceMethod

	router          Router
	validator       Validator
	httpMiddlewares []func(handler http.Handler) http.Handler
}

// route 注册路由
func (s *Service) route(r Router) {
	mux := r
	if len(s.httpMiddlewares) > 0 {
		mux = r.With(s.httpMiddlewares...)
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

func (s *Service) Endpoint() http.Handler {
	mux := NewServeMux()
	s.route(mux)
	return mux
}

// NewService 创建服务
func NewService(s interface{}, opts ...Option) *Service {
	v := reflect.Indirect(reflect.ValueOf(s))
	t := v.Type()

	if !strings.HasSuffix(t.Name(), "Service") {
		panic(errors.New("service name必须包含Service后缀"))
	}

	srv := &Service{
		Name: strings.TrimSuffix(t.Name(), "Service"),
	}
	o := newOptions()
	for i := range opts {
		opts[i].apply(o)
	}
	srv.validator = o.validator
	srv.httpMiddlewares = o.httpMiddlewares

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
