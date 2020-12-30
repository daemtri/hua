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
	http.Handler
	// Handle and HandleFunc adds routes for `pattern` that matches
	// all HTTP methods.
	Handle(pattern string, h http.Handler)
	// Method and MethodFunc adds routes for `pattern` that matches
	// the `method` HTTP method.
	Method(method, pattern string, h http.Handler)
	// URLParam returns the url parameter from a http.Request object.
	URLParams(r *http.Request, key string) string
	// With adds inline middlewares for an endpoint handler.
	With(middlewares ...func(http.Handler) http.Handler) Router
	// Use appends one or more middlewares onto the Router stack.
	Use(middlewares ...func(http.Handler) http.Handler)
}

type mux struct {
	chi.Router
}

func (m *mux) URLParams(r *http.Request, key string) string {
	return chi.URLParam(r, key)
}

func (m *mux) With(middlewares ...func(http.Handler) http.Handler) Router {
	return &mux{m.Router.With(middlewares...)}
}

// NewServeMux create a ServeMux
var NewServeMux = func() Router { return &mux{chi.NewRouter()} }
