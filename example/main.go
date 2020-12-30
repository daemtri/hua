package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/duanqy/hua/pkg/huamock"

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
	// 作为客户端调用stub方法
	ret, err := s.Mul(context.Background(), &api.MulArg{
		Left:  1,
		Right: 2,
	})
	fmt.Println(ret)
	err = http.ListenAndServe(":80", huarpc.NewServer().Register(s))
	if err != nil {
		log.Fatalln(err)
	}

}
