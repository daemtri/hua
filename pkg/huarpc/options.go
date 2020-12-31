package huarpc

import (
	"net/http"
)

func WithMiddleware(f func(http.Handler) http.Handler) Option {
	return optionFunc(func(o *options) {
		o.httpMiddlewares = append(o.httpMiddlewares, f)
	})
}

type options struct {
	validator       Validator
	httpMiddlewares []func(handler http.Handler) http.Handler
}

func newOptions() *options {
	return &options{
		validator: nil,
	}
}

type Option interface {
	apply(*options)
}

type optionFunc func(*options)

func (f optionFunc) apply(o *options) {
	f(o)
}
