package huarpc

import (
	"encoding/json"
	"fmt"
	"go/token"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-chi/chi"
)

func indirect(v reflect.Type) reflect.Type {
	if v.Kind() != reflect.Ptr {
		return v
	}
	return v.Elem()
}

func BuildServer(s interface{}) Server {
	v := reflect.Indirect(reflect.ValueOf(s))
	t := v.Type()
	if !strings.HasSuffix(t.Name(), "Service") {
		panic("service name must end with 'Service'")
	}

	mux := chi.NewRouter()
	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		if !token.IsExported(sf.Name) {
			panic("service的所有属性必须是可导出的")
		}
		sft := sf.Type
		if sft.NumIn() != 2 || sft.NumOut() != 2 {
			panic(fmt.Errorf("service%s的所有方法必须要包含一个入参，2个出参", t.Name()))
		}

		if sft.Out(1) != errType {
			panic(fmt.Errorf("方法%s的第二个返回值必须是error", sf.Name))
		}

		if !strings.HasSuffix(indirect(sft.In(1)).Name(), "Arg") {
			panic(fmt.Errorf("%s must end with 'Arg'", indirect(sft.In(1)).Name()))
		}
		if !strings.HasSuffix(indirect(sft.Out(0)).Name(), "Reply") {
			panic(fmt.Errorf("%s must end with 'Reply'", indirect(sft.Out(0)).Name()))
		}

		m, p := http.MethodPost, fmt.Sprintf("/%s/%s", t.Name(), sf.Name)
		if httpTag, ok := sf.Tag.Lookup("http"); ok {
			httpTags := strings.SplitN(httpTag, ",", 2)
			m = httpTags[0]
			if len(httpTags) == 2 {
				p = httpTags[1]
			}
		}

		mux.Method(m, p, &method{
			argType:   sft.In(1),
			replyType: sft.Out(0),
			callable:  v.Field(i),
		})
	}

	return Server{mux: mux}
}

type method struct {
	argType   reflect.Type
	replyType reflect.Type
	callable  reflect.Value
}

func (m *method) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	arg := reflect.New(m.argType.Elem())
	contentType := r.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "application/json") {
		if err := json.NewDecoder(r.Body).Decode(arg.Interface()); err != nil {
			panic(fmt.Errorf("json decode error: %w", err))
		}
	} else {
		panic(fmt.Errorf("content-type error:accept only application/json now, recived: %s", contentType))
	}

	reply := m.callable.Call([]reflect.Value{reflect.ValueOf(r.Context()), arg})
	err, _ := reply[1].Interface().(error)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(reply[0].Interface())
}
