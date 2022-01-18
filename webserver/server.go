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
}

func (s *WebServer) Route(method string, pattern string, handler func(Context)) {
	s.handler.Route(method, pattern, handler)
}

func (s *WebServer) Start(address string) error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		c := &WebContext{W: w, R: r}
		method, path := c.R.Method, c.R.URL.Path
		s.handler.Handle(method, path, c)
	})
	err := http.ListenAndServe(address, nil)
	return err
}

func NewServer() Server {
	return &WebServer{handler: NewHandler()}
}
