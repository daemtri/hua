package main

import (
	"context"
	"log"
	"net/http"

	"github.com/duanqy/hua/example/api"
	"github.com/duanqy/hua/pkg/huarpc"
)

func main() {

	err := http.ListenAndServe(":80", huarpc.BuildServer(&api.CalcService{
		Add: func(ctx context.Context, arg *api.AddArg) (*api.AddReply, error) {
			return &api.AddReply{Result: arg.Left + arg.Right}, nil
		},
	}))
	if err != nil {
		log.Fatalln(err)
	}
}
