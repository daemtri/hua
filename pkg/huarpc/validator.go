package huarpc

import "github.com/go-playground/validator/v10"

type Validator interface {
	ValidateStruct(x interface{}) error
}

type validate struct {
	*validator.Validate
}

func (v *validate) ValidateStruct(x interface{}) error {
	return v.Struct(x)
}
