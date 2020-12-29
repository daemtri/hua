package huarpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go/token"
	"net/http"
	"reflect"
	"regexp"
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
	Name    string
	ArgType reflect.Type

	Tags MethodTags

	Handler reflect.Value

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
		Handler:     value,
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

	reply := s.Handler.Call([]reflect.Value{reflect.ValueOf(r.Context()), arg})
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
	Name       string
	NamePrefix string
	Version    string
	Methods    []*ServiceMethod
}

func (s Service) Route(r Router) {
	var restPattern, rpcPattern string
	for i := range s.Methods {
		m := s.Methods[i]
		rpcPattern = fmt.Sprintf("%s%s.%s", s.Name, s.Version, m.Name)
		if s.Version != "" {
			restPattern = fmt.Sprintf("/%s/%s/%s", strings.ToLower(s.Name), strings.ToLower(s.Version), strings.ToLower(m.Name))
		} else {
			restPattern = fmt.Sprintf("/%s/%s", strings.ToLower(s.Name), strings.ToLower(m.Name))
		}
		if m.Tags.HTTP.Method == "" {
			r.Handle(restPattern, m)
			r.Handle(rpcPattern, m)
		} else {
			r.Method(m.Tags.HTTP.Method, restPattern, m)
			r.Method(m.Tags.HTTP.Method, rpcPattern, m)
		}
	}
}

var serviceNamePattern = regexp.MustCompile(`^(?:([a-zA-Z0-9]*)(V[0-9]+)|([a-zA-Z0-9]*))Service$`)

func MustNewService(s interface{}) *Service {
	service, err := NewService(s)
	if err != nil {
		panic(fmt.Errorf("new service error: %w", err))
	}
	return service
}

func NewService(s interface{}) (*Service, error) {
	v := reflect.Indirect(reflect.ValueOf(s))
	t := v.Type()
	ret := serviceNamePattern.FindStringSubmatch(t.Name())
	if ret == nil {
		return nil, errors.New("service name必须包含Service后缀")
	}

	srv := &Service{}
	if ret[3] == "" {
		srv.Name, srv.Version = ret[1], ret[2]
	} else {
		srv.Name = ret[3]
	}

	for i := 0; i < t.NumField(); i++ {
		sm, err := parseServiceMethod(t.Field(i), v.Field(i))
		if err != nil {
			return nil, fmt.Errorf("解析方法出错: %w", err)
		}
		srv.Methods = append(srv.Methods, sm)
	}

	return srv, nil
}
