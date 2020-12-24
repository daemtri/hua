package api

import (
	"context"
)

type AddArg struct {
	Left  int `json:"left" help:"左值"`
	Right int `json:"right" help:"右值"`
}

type AddReply struct {
	Result int `json:"result"`
}

type DivArg struct {
	Left  int `path:"left" help:"左值"`
	Right int `path:"right" help:"右值"`
}

type DivReply struct {
	Result int `json:"result"`
}

type CalcService struct {
	Add func(context.Context, *AddArg) (*AddReply, error) `http:"GET,/calc/add" help:"加法"`
	Div func(context.Context, *DivArg) (*DivReply, error) `http:"GET,/calc/div/{left}/{right}" help:"除法"`
}
