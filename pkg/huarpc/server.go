package huarpc

import (
	"net/http"
)

type Server struct {
	Router
	host string
}

func NewServer(opts ...Option) *Server {
	s := &Server{
		Router: NewServeMux(),
	}
	o := newOptions()
	for i := range opts {
		opts[i].apply(o)
	}
	s.host = o.serverHost
	for i := range o.httpMiddlewares {
		s.Router.Use(o.httpMiddlewares[i])
	}
	return s
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.Router.ServeHTTP(writer, request)
}

func (s *Server) Run(addr string) error {
	return http.ListenAndServe(addr, s)
}

func (s *Server) Register(service interface{}) *Server {
	srv := NewService(service)
	srv.route(s.Router)
	return s
}
