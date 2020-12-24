package huarpc

import (
	"context"
	"net/http"
)

type Router interface {
	http.Handler
	// Method and MethodFunc adds routes for `pattern` that matches
	// the `method` HTTP method.
	HandleFunc(method, pattern string, handler http.Handler)
}

func Register(x interface{}) {

}

type Context = context.Context
