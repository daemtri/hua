package huarpc

import (
	"net/http"
	"reflect"
)

func isInvalidKind(kind reflect.Kind) bool {
	return kind == reflect.Int8 ||
		kind == reflect.Int16 ||
		kind == reflect.Int32 ||
		kind == reflect.Int64 ||
		kind == reflect.Int ||
		kind == reflect.Uint8 ||
		kind == reflect.Uint16 ||
		kind == reflect.Uint32 ||
		kind == reflect.Uint64 ||
		kind == reflect.Uint ||
		kind == reflect.Float32 ||
		kind == reflect.Float64 ||
		kind == reflect.Bool ||
		kind == reflect.String
}

type argField struct {
	form string
	kind reflect.Kind
}

type argFields struct {
	fields []argField
	typ    reflect.Type
}

func formDecode(fields *argFields, r *http.Request) (*reflect.Value, error) {
	return nil, nil
}
