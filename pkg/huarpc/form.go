package huarpc

import (
	"net/url"

	"github.com/go-playground/form/v4"
)

type FormDecoder interface {
	Decode(v interface{}, values url.Values) (err error)
}

type FormEncoder interface {
	Encode(v interface{}) (values url.Values, err error)
}

func newFormDecoder() FormDecoder {
	dec := form.NewDecoder()
	return dec
}

func newFromEncoder() FormEncoder {
	enc := form.NewEncoder()
	return enc
}
