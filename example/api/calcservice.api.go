package api

import (
	"context"
)

// CalcService 算数计算服务
type CalcService struct {
	Add       func(context.Context, AddArg) (*AddReply, error)                    `http:"GET /add" help:"加法"`
	Div       func(context.Context, DivArg) (*DivReply, error)                    `http:"GET /div/{left}/{right}" help:"除法"`
	Mul       func(context.Context, MulArg) (*MulReply, error)                    `http:"GET /mul" help:"除法"`
	Fibonacci func(context.Context, FibonacciArg) (<-chan *FibonacciReply, error) `http:"GET /fibonacci" sse-retry:"100" help:"斐波那契数列"`
}

// FibonacciArg 请求参数
type FibonacciArg struct {
	First  int `json:"first"`
	Second int `json:"second"`
}

// FibonacciReply 返回参数
type FibonacciReply struct {
	ID     string `json:"-" sse:"id"`
	Number int    `json:"number"`
}

// MulArg 请求参数
type MulArg struct {
	Left  float64 `json:"left" form:"left" valid:"int" help:"左值"`
	Right float64 `json:"right" form:"right" valid:"int" help:"右值"`
}

// MulReply 返回参数
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
