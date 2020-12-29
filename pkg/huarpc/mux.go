package huarpc

import (
	"net/http"

	"github.com/go-chi/chi"
)

type Requester interface {
	Do(req *http.Request) (*http.Response, error)
}

// NewRequester create a Requester
var NewRequester = func() Requester { return http.DefaultClient }

type Router interface {
	Method(method, pattern string, h http.Handler)
	// Handle and HandleFunc adds routes for `pattern` that matches
	// all HTTP methods.
	Handle(pattern string, h http.Handler)
}

type ServeMux interface {
	http.Handler
	Router

	// Use appends one or more middlewares onto the Router stack.
	Use(middlewares ...func(http.Handler) http.Handler)
}

// NewServeMux create a ServeMux
var NewServeMux = func() ServeMux { return chi.NewRouter() }
