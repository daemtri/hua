package huarpc

import (
	"reflect"
)

var (
	argeAllowedKind = []reflect.Kind{
		reflect.Struct,
	}
	replyAllowedKind = []reflect.Kind{
		reflect.Struct,
		reflect.Ptr,
		reflect.Chan,
	}

	valueAllowedKind = []reflect.Kind{
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
		reflect.String,
		reflect.Struct,
		reflect.Bool,
		reflect.Slice,
		reflect.Array,
	}

	mapKeyAllowedKind = []reflect.Kind{
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
		reflect.String,
	}
)

func isAllowedKind(allowed []reflect.Kind, k reflect.Kind) bool {
	for i := range allowed {
		if allowed[i] == k {
			return true
		}
	}

	return false
}

type unaryHandler struct {
}

type streamHandler struct {
}
