package webserver

import "net/http"

type Handler interface {
	Route(method string, pattern string, handler func(Context))
	Handle(method string, pattern string, c Context)
}

var _ Handler = &WebContextHandler{}

type WebContextHandler struct {
	handlers map[string]func(Context)
}

func (h *WebContextHandler) Route(method string, pattern string, handler func(Context)) {
	key := getKey(method, pattern)
	h.handlers[key] = handler
}

func (h *WebContextHandler) Handle(method string, pattern string, c Context) {
	key := getKey(method, pattern)
	handler, success := h.handlers[key]
	if !success {
		c.WriteCode(http.StatusNotFound)
	} else {
		handler(c)
	}
}

func getKey(method string, pattern string) string {
	return method + "#" + pattern
}

func NewHandler() Handler {
	return &WebContextHandler{handlers: make(map[string]func(Context))}
}
