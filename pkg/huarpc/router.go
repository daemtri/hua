package huarpc

import (
	"net/http"
)

type Router interface {
	http.Handler
	// Method and MethodFunc adds routes for `pattern` that matches
	// the `method` HTTP method.
	HandleFunc(method, pattern string, handler http.Handler)
}
