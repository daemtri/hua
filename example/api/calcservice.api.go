package api

import (
	"context"
)

// CalcService xxxxxx
// 哈哈哈
type CalcService struct {
	// Add 函数xxx
	Add func(context.Context, *AddArg) (*AddReply, error) `http:"GET,/calc/add" help:"加法"`
	// Div xxx
	Div func(context.Context, *DivArg) (*DivReply, error) `http:"GET,/calc/div/{left}/{right}" help:"除法"`
}

// AddArg xxx
type AddArg struct {
	Left  int `json:"left" help:"左值"`
	Right int `json:"right" help:"右值"`
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
