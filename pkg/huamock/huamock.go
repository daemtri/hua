package huamock

import (
	"fmt"
	"reflect"
	"time"

	"github.com/bxcodec/faker/v3"
)

// Stub 伪造未实现的方法用于测试
func Stub(s interface{}) error {
	if reflect.TypeOf(s).Kind() != reflect.Ptr {
		return fmt.Errorf("need s kind of ptr, got: %T", s)
	}
	v := reflect.Indirect(reflect.ValueOf(s))
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		if !v.Field(i).IsNil() {
			continue
		}
		mt := t.Field(i).Type
		fh := &fakerHandler{
			replyType: mt.Out(0),
			errType:   mt.Out(1),
		}
		var fn reflect.Value
		if mt.Out(0).Kind() == reflect.Chan {
			fn = reflect.MakeFunc(mt, fh.stubChan)
		} else {
			fn = reflect.MakeFunc(mt, fh.stub)
		}

		v.Field(i).Set(fn)
	}

	return nil
}

type fakerHandler struct {
	replyType reflect.Type
	errType   reflect.Type
}

func (f *fakerHandler) stubChan(_ []reflect.Value) []reflect.Value {
	replyChan := reflect.MakeChan(reflect.ChanOf(reflect.BothDir, f.replyType.Elem()), 0)
	err := reflect.Zero(f.errType)

	go func() {
		for {
			reply := reflect.New(f.replyType.Elem().Elem())
			_ = faker.FakeData(reply.Interface())
			replyChan.Send(reply)
			time.Sleep(1 * time.Second)
		}
	}()

	return []reflect.Value{replyChan, err}
}

func (f *fakerHandler) stub(_ []reflect.Value) []reflect.Value {
	reply := reflect.New(f.replyType.Elem())
	err := reflect.Zero(f.errType)

	if err2 := faker.FakeData(reply.Interface()); err2 != nil {
		err = reflect.ValueOf(err)
	}

	return []reflect.Value{reply, err}
}
