package webserver

import (
	"context"
	"net/http"
	"sync"
)

type Server interface {
	Route(method string, pattern string, handler func(Context))
	Start(address string) error
	Shutdown(ctx context.Context) error
}

var _ Server = &WebServer{}

type WebServer struct {
	handler     Handler
	root        Filter
	contextPool sync.Pool
}

func (s *WebServer) Route(method string, pattern string, handler func(Context)) {
	s.handler.Route(method, pattern, handler)
}

func (s *WebServer) Start(address string) error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		c := s.contextPool.Get().(Context)
		c.SetW(w)
		c.SetR(r)
		s.root(c)
	})
	err := http.ListenAndServe(address, nil)
	return err
}

func (s *WebServer) Shutdown(ctx context.Context) error {
	//do something
	return nil
}

func NewServer(builders ...FilterBuilder) Server {
	handler := NewHandler2()
	root := handler.ServeHTTP

	for i := len(builders) - 1; i >= 0; i-- {
		b := builders[i]
		root = b(root)
	}

	return &WebServer{
		handler: handler,
		root:    root,
		contextPool: sync.Pool{
			New: func() interface{} {
				return NewContext(nil, nil)
			},
		},
	}
}
