package huarpc

import (
	"reflect"
)

var (
	argValueAllowedType = []reflect.Kind{
		reflect.Struct,
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Float32,
		reflect.Float64,
		reflect.Bool,
		reflect.String,
		reflect.Slice,
		reflect.Array,
	}
)
