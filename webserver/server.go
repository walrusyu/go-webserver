package webserver

import (
	"net/http"
)

type Server interface {
	Route(method string, pattern string, handler func(Context))
	Start(address string) error
}

var _ Server = &WebServer{}

type WebServer struct {
	handler Handler
	root    Filter
}

func (s *WebServer) Route(method string, pattern string, handler func(Context)) {
	s.handler.Route(method, pattern, handler)
}

func (s *WebServer) Start(address string) error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		c := &WebContext{W: w, R: r}
		s.root(c)
	})
	err := http.ListenAndServe(address, nil)
	return err
}

func NewServer(builders ...FilterBuilder) Server {
	handler := NewHandler()
	root := handler.ServeHTTP

	for i := len(builders) - 1; i >= 0; i-- {
		b := builders[i]
		root = b(root)
	}

	return &WebServer{handler: handler, root: root}
}
