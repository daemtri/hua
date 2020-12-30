package main

import (
	"net/http"

	"github.com/go-chi/chi"

	"github.com/duanqy/hua/example/api"
	"github.com/duanqy/hua/pkg/huamock"
	"github.com/duanqy/hua/pkg/huarpc"
)

func main() {
	r := chi.NewRouter()
	_, b := huarpc.NewService(huamock.Stub(&api.CalcService{})).Endpoint()
	if err := http.ListenAndServe(":80", r); err != nil {
		panic(err)
	}
}
