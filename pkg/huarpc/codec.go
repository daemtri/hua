package huarpc

import (
	formv4 "github.com/go-playground/form/v4"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

var form = struct {
	Decoder *formv4.Decoder
	Encoder *formv4.Encoder
}{
	Decoder: formv4.NewDecoder(),
	Encoder: formv4.NewEncoder(),
}
