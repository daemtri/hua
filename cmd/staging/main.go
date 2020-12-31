package main

import (
	"net/http"

	"github.com/duanqy/hua/example/api"
	"github.com/duanqy/hua/pkg/huamock"
	"github.com/duanqy/hua/pkg/huarpc"
)

func main() {
	if err := http.ListenAndServe(":80", huarpc.NewService(huamock.Stub(&api.CalcService{})).Endpoint()); err != nil {
		panic(err)
	}
}
