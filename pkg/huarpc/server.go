package huarpc

import (
	"io/fs"
	"net/http"

	"github.com/go-chi/chi"
)

type Validator interface {
	ValidateStruct(x interface{}) error
}

func WithValidator() Option {
	return optionFunc(func(o *options) {

	})
}

func WithProtocol(fs fs.ReadDirFS) Option {
	return optionFunc(func(o *options) {
		o.protocol = fs
	})
}

func WithServerHost(host string) Option {
	return optionFunc(func(o *options) {
		o.serverHost = host
	})
}

func WithHttpMiddleware(f func(http.Handler) http.Handler) Option {
	return optionFunc(func(o *options) {
		o.httpMiddlewares = append(o.httpMiddlewares, f)
	})
}

type options struct {
	validator  Validator
	protocol   fs.ReadDirFS
	serverHost string

	httpMiddlewares []func(handler http.Handler) http.Handler
}

func newOptions() *options {
	return &options{
		validator:       nil,
		protocol:        nil,
		serverHost:      "http://127.0.0.1",
		httpMiddlewares: nil,
	}
}

type Option interface {
	apply(*options)
}

type optionFunc func(*options)

func (f optionFunc) apply(o *options) {
	f(o)
}

type Server struct {
	host string
	mux  Router
}

func NewServer(opts ...Option) *Server {
	s := &Server{
		mux: chi.NewRouter(),
	}
	o := newOptions()
	for i := range opts {
		opts[i].apply(o)
	}
	s.host = o.serverHost
	for i := range o.httpMiddlewares {
		s.mux.Use(o.httpMiddlewares[i])
	}
	return nil
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.mux.ServeHTTP(writer, request)
}

func (s *Server) Register(service interface{}) *Server {
	srv, err := parseService(service)
	if err != nil {
		panic(err)
	}
	srv.route(s.mux)
	return nil
}