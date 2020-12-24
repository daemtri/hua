package huarpc

import (
	"encoding/json"
	"fmt"
	"go/token"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
)

var (
	errType = reflect.TypeOf(struct{ error }{}).Field(0).Type
)

func indirect(v reflect.Type) reflect.Type {
	if v.Kind() != reflect.Ptr {
		return v
	}
	return v.Elem()
}

func BuildServer(s interface{}) http.Handler {
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
			httpTags := strings.SplitN(httpTag, " ", 2)
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

	return mux
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
		argType := m.argType.Elem()
		for i := 0; i < argType.NumField(); i++ {
			var val string
			valType := argType.Field(i)
			formTag, ok := valType.Tag.Lookup("form")
			if ok {
				val = r.FormValue(formTag)
			} else {
				val = r.FormValue(argType.Field(i).Name)
			}

			switch valType.Type.Kind() {
			case reflect.Int8:
				fieldValue, err := strconv.ParseInt(val, 10, 8)
				if err != nil {
					http.Error(w, "解析参数出错"+err.Error(), http.StatusInternalServerError)
					return
				}
				arg.Elem().Field(i).SetInt(fieldValue)
			case reflect.Int16:
				fieldValue, err := strconv.ParseInt(val, 10, 16)
				if err != nil {
					http.Error(w, "解析参数出错"+err.Error(), http.StatusInternalServerError)
					return
				}
				arg.Elem().Field(i).SetInt(fieldValue)
			case reflect.Int32:
				fieldValue, err := strconv.ParseInt(val, 10, 32)
				if err != nil {
					http.Error(w, "解析参数出错"+err.Error(), http.StatusInternalServerError)
					return
				}
				arg.Elem().Field(i).SetInt(fieldValue)
			case reflect.Int:
				fieldValue, err := strconv.ParseInt(val, 10, strconv.IntSize)
				if err != nil {
					http.Error(w, "解析参数出错"+err.Error(), http.StatusInternalServerError)
					return
				}
				arg.Elem().Field(i).SetInt(fieldValue)
			case reflect.Int64:
				fieldValue, err := strconv.ParseInt(val, 10, 64)
				if err != nil {
					http.Error(w, "解析参数出错"+err.Error(), http.StatusInternalServerError)
					return
				}
				arg.Elem().Field(i).SetInt(fieldValue)
			case reflect.Uint8:
				fieldValue, err := strconv.ParseUint(val, 10, 6)
				if err != nil {
					http.Error(w, "解析参数出错"+err.Error(), http.StatusInternalServerError)
					return
				}
				arg.Elem().Field(i).SetUint(fieldValue)
			case reflect.Uint16:
				fieldValue, err := strconv.ParseUint(val, 10, 16)
				if err != nil {
					http.Error(w, "解析参数出错"+err.Error(), http.StatusInternalServerError)
					return
				}
				arg.Elem().Field(i).SetUint(fieldValue)
			case reflect.Uint32:
				fieldValue, err := strconv.ParseUint(val, 10, 32)
				if err != nil {
					http.Error(w, "解析参数出错"+err.Error(), http.StatusInternalServerError)
					return
				}
				arg.Elem().Field(i).SetUint(fieldValue)
			case reflect.Uint:
				fieldValue, err := strconv.ParseUint(val, 10, strconv.IntSize)
				if err != nil {
					http.Error(w, "解析参数出错"+err.Error(), http.StatusInternalServerError)
					return
				}
				arg.Elem().Field(i).SetUint(fieldValue)
			case reflect.Uint64:
				fieldValue, err := strconv.ParseUint(val, 10, 64)
				if err != nil {
					http.Error(w, "解析参数出错"+err.Error(), http.StatusInternalServerError)
					return
				}
				arg.Elem().Field(i).SetUint(fieldValue)
			case reflect.Float32:
				fieldValue, err := strconv.ParseFloat(val, 32)
				if err != nil {
					http.Error(w, "解析参数出错"+err.Error(), http.StatusInternalServerError)
					return
				}
				arg.Elem().Field(i).SetFloat(fieldValue)
			case reflect.Float64:
				fieldValue, err := strconv.ParseFloat(val, 64)
				if err != nil {
					http.Error(w, "解析参数出错"+err.Error(), http.StatusInternalServerError)
					return
				}
				arg.Elem().Field(i).SetFloat(fieldValue)
			case reflect.String:
				arg.Elem().Field(i).SetString(val)
			case reflect.Bool:
				fieldValue, err := strconv.ParseBool(val)
				if err != nil {
					http.Error(w, "解析参数出错"+err.Error(), http.StatusInternalServerError)
					return
				}
				arg.Elem().Field(i).SetBool(fieldValue)
			default:
				http.Error(w, "不支持的字段类型", http.StatusInternalServerError)
			}

		}
	}

	reply := m.callable.Call([]reflect.Value{reflect.ValueOf(r.Context()), arg})
	err, _ := reply[1].Interface().(error)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(reply[0].Interface())
}
