package huamock

import (
	"errors"
	"github.com/bxcodec/faker/v3"
	"reflect"
)

// Stub 伪造未实现的方法用于测试
func Stub(s interface{}) error {
	v := reflect.ValueOf(s)
	t := v.Type()
	if t.Kind() != reflect.Ptr {
		return errors.New("service type must be a ptr")
	}

	t = t.Elem()
	v = v.Elem()
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

	return nil
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
