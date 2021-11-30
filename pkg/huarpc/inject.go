package huarpc

import (
	"fmt"
	"reflect"
)

// Inject 把h结构体的方法注入到s结构体的属性之中
func Inject(dst interface{}, src interface{}) error {
	dstVal := reflect.ValueOf(dst)
	if dstVal.Kind() != reflect.Ptr {
		return fmt.Errorf("need dst kind of ptr, got: %T", dst)
	}
	dstVal = reflect.Indirect(dstVal)
	if dstVal.Kind() != reflect.Struct {
		return fmt.Errorf("need dst indirect kind of struct, got: %T", dst)
	}
	dstTyp := dstVal.Type()

	srcTyp := reflect.TypeOf(src)
	if srcTyp.Kind() == reflect.Ptr {
		srcTyp = srcTyp.Elem()
	}
	srcVal := reflect.Indirect(reflect.ValueOf(src))

	for i := 0; i < srcVal.NumMethod(); i++ {
		name := srcTyp.Method(i).Name
		srcMethodVal := srcVal.Method(i)
		for j := 0; j < dstTyp.NumField(); j++ {
			dstFieldVal := dstVal.Field(j)
			if name == dstTyp.Field(j).Name &&
				reflect.DeepEqual(dstFieldVal.Type(), srcMethodVal.Type()) {
				dstFieldVal.Set(srcMethodVal)
			}
		}
	}

	return nil
}
