package api

import "context"

type AddArg struct {
	Left  int `form json:"left" help:"左值"`
	Right int `form json:"right" help:"右值"`
}

type AddReply struct {
	Result int `json:"result"`
}

type CalcService struct {
	Add func(context.Context, *AddArg) (*AddReply, error) `http:"GET /calc/add" help:"获取用户信息"`
}
