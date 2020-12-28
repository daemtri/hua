package main

import (
	"context"
	"github.com/duanqy/hua/pkg/huamock"
	"log"
	"net/http"

	"github.com/duanqy/hua/example/api"
	"github.com/duanqy/hua/pkg/huarpc"
)

func main() {
	s := &api.CalcService{
		Add: func(ctx context.Context, arg *api.AddArg) (*api.AddReply, error) {
			return &api.AddReply{Result: arg.Left + arg.Right}, nil
		},
	}

	if err := huamock.Stub(s); err != nil {
		panic(err)
	}
	err := http.ListenAndServe(":80", huarpc.NewServer().Register(s))
	if err != nil {
		log.Fatalln(err)
	}
}
