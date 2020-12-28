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
	"text/template"
)

var (
	errType     = reflect.TypeOf(struct{ error }{}).Field(0).Type
	contextType = reflect.TypeOf(struct{ context.Context }{}).Field(0).Type
)

type serviceMethod struct {
	name      string
	argType   reflect.Type
	replyType reflect.Type
	httpTags  struct {
		method string
		path   string
	}

	handler reflect.Value

	formDecoder FormDecoder
	formEncoder FormEncoder
}

func parseServiceMethod(field reflect.StructField, value reflect.Value) (*serviceMethod, error) {
	if !token.IsExported(field.Name) {
		panic("service的所有属性必须是可导出的")
	}
	sft := field.Type
	if sft.NumIn() != 2 || sft.NumOut() != 2 || sft.In(0) != contextType || sft.Out(1) != errType {
		return nil, fmt.Errorf("函数%s签名错误（正确: func(context.Context,T)(K,error),T为入参，K为出参)", field.Name)
	}

	if !strings.HasSuffix(indirect(sft.In(1)).Name(), "Arg") {
		panic(fmt.Errorf("%s must end with 'Arg'", indirect(sft.In(1)).Name()))
	}
	if !strings.HasSuffix(indirect(sft.Out(0)).Name(), "Reply") {
		panic(fmt.Errorf("%s must end with 'Reply'", indirect(sft.Out(0)).Name()))
	}

	sm := &serviceMethod{
		name:        field.Name,
		argType:     sft.In(1),
		replyType:   sft.Out(0),
		handler:     value,
		formDecoder: newFormDecoder(),
		formEncoder: newFromEncoder(),
	}
	if httpTag, ok := field.Tag.Lookup("http"); ok {
		httpTag := strings.SplitN(httpTag, ",", 2)
		sm.httpTags.method = httpTag[0]
		sm.httpTags.path = httpTag[1]
	}

	return sm, nil
}

func (s *serviceMethod) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	arg := reflect.New(s.argType.Elem())
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
		if err := s.formDecoder.Decode(arg.Interface(), r.Form); err != nil {
			http.Error(w, fmt.Sprintf("decode form error: %s", err), http.StatusBadRequest)
			return
		}
	}

	reply := s.handler.Call([]reflect.Value{reflect.ValueOf(r.Context()), arg})
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

type service struct {
	patternTemplate template.Template
	structName      string

	name    string
	version string

	methods []*serviceMethod
}

func (s service) route(mux Router) {
	var pattern string
	for i := range s.methods {
		m := s.methods[i]
		if s.version != "" {
			pattern = fmt.Sprintf("/%s/%s/%s", strings.ToLower(s.name), strings.ToLower(s.version), strings.ToLower(m.name))
		} else {
			pattern = fmt.Sprintf("%s/%s", strings.ToLower(s.name), strings.ToLower(m.name))
		}
		if m.httpTags.method == "" {
			mux.Handle(pattern, m)
		} else {
			mux.Method(m.httpTags.method, pattern, m)
		}
	}
}

var serviceNamePattern = regexp.MustCompile(`^(?:([a-zA-Z0-9]*)(V[0-9]+)|([a-zA-Z0-9]*))Service$`)

func parseService(s interface{}) (*service, error) {
	v := reflect.Indirect(reflect.ValueOf(s))
	t := v.Type()
	ret := serviceNamePattern.FindStringSubmatch(t.Name())
	if ret == nil {
		return nil, errors.New("service name必须包含Service后缀")
	}

	srv := &service{}
	if ret[3] == "" {
		srv.name, srv.version = ret[1], ret[2]
	} else {
		srv.name = ret[3]
	}

	for i := 0; i < t.NumField(); i++ {
		sm, err := parseServiceMethod(t.Field(i), v.Field(i))
		if err != nil {
			return nil, fmt.Errorf("解析方法出错: %w", err)
		}
		srv.methods = append(srv.methods, sm)
	}

	return srv, nil
}
