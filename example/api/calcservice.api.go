package api

import (
	"context"
)

// CalcService xxxxxx
type CalcService struct {
	// Add 函数xxx
	Add func(context.Context, *AddArg) (*AddReply, error) `http:"GET /calc/add" help:"加法"`
	// Div xxx
	Div func(context.Context, *DivArg) (*DivReply, error) `http:"GET /calc/div/{left}/{right}" help:"除法"`
	// Div xxx
	Mul func(context.Context, *MulArg) (*MulReply, error) `http:"GET /calc/mul" help:"除法"`
}

type MulArg struct {
	Left  float64 `json:"left" form:"left" valid:"int" help:"左值"`
	Right float64 `json:"right" form:"right" valid:"int" help:"右值"`
}

type MulReply struct {
	Result float64 `json:"result"`
}

// AddArg xxx
type AddArg struct {
	Left  int `json:"left" form:"left" valid:"int" help:"左值"`
	Right int `json:"right" form:"right" valid:"int" help:"右值"`
}

// AddReply xxx
type AddReply struct {
	Result int `json:"result"`
}

// DivArg xxx
type DivArg struct {
	Left  int `path:"left" help:"左值"`
	Right int `path:"right" help:"右值"`
}

// DivReply xxx
type DivReply struct {
	Result int `json:"result"`
}
