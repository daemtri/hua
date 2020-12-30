package huamock

import (
	"reflect"

	"github.com/bxcodec/faker/v3"
)

// Stub 伪造未实现的方法用于测试
func Stub(s interface{}) interface{} {
	v := reflect.Indirect(reflect.ValueOf(s))
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		if !v.Field(i).IsNil() {
			continue
		}
		mt := t.Field(i).Type
		fn := reflect.MakeFunc(mt, (&fakerHandler{
			replyType: mt.Out(0),
			errType:   mt.Out(1),
		}).stub)
		v.Field(i).Set(fn)
	}

	return s
}

type fakerHandler struct {
	replyType reflect.Type
	errType   reflect.Type
}

func (f *fakerHandler) stub(_ []reflect.Value) []reflect.Value {
	reply := reflect.New(f.replyType.Elem())
	err := reflect.Zero(f.errType)
	if err2 := faker.FakeData(reply.Interface()); err2 != nil {
		err = reflect.ValueOf(err)
	}

	return []reflect.Value{reply, err}
}
