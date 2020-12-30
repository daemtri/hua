package huarpc

import (
	"io/fs"
	"net/http"
)

type Validator interface {
	ValidateStruct(x interface{}) error
}

func WithValidator() Option {
	return optionFunc(func(o *options) {

	})
}

func WithSchema(fs fs.ReadDirFS) Option {
	return optionFunc(func(o *options) {
		o.protocol = fs
	})
}

func WithServerHost(host string) Option {
	return optionFunc(func(o *options) {
		o.serverHost = host
	})
}

func WithMiddleware(f func(http.Handler) http.Handler) Option {
	return optionFunc(func(o *options) {
		o.httpMiddlewares = append(o.httpMiddlewares, f)
	})
}

type options struct {
	validator  Validator
	protocol   fs.ReadDirFS
	serverHost string

	httpMiddlewares []func(handler http.Handler) http.Handler
}

func newOptions() *options {
	return &options{
		validator:  nil,
		protocol:   nil,
		serverHost: "HTTP://127.0.0.1",
	}
}

type Option interface {
	apply(*options)
}

type optionFunc func(*options)

func (f optionFunc) apply(o *options) {
	f(o)
}
