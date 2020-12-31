package huarpc

import "github.com/go-playground/validator/v10"

type Validator interface {
	ValidateStruct(x interface{}) error
}

func WithCustomValidator(v Validator) Option {
	return optionFunc(func(o *options) {
		o.validator = v
	})
}

func EnableValidator() Option {
	return optionFunc(func(o *options) {
		o.validator = &validate{validator.New()}
	})
}

type validate struct {
	*validator.Validate
}

func (v *validate) ValidateStruct(x interface{}) error {
	return v.Struct(x)
}
