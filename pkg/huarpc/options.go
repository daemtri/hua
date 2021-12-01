package huarpc

import "github.com/go-playground/validator/v10"

func WithRouter(r Router) Option {
	return optionFunc(func(o *options) {
		o.router = r
	})
}

func WithValidator(v Validator) Option {
	return optionFunc(func(o *options) {
		o.validator = v
	})
}

func WithInstance(x interface{}) Option {
	return optionFunc(func(opt *options) {
		opt.instance = x
	})
}

type options struct {
	validator Validator
	router    Router
	instance  interface{}
}

func newOptions() *options {
	return &options{
		validator: &validate{validator.New()},
		router:    NewServeMux(),
	}
}

type Option interface {
	apply(*options)
}

type optionFunc func(*options)

func (f optionFunc) apply(o *options) {
	f(o)
}
