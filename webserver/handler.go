package webserver

import "net/http"

type Handler interface {
	Route(method string, pattern string, handler func(Context))
	ServeHTTP(c Context)
}

var _ Handler = &WebContextHandler{}

type WebContextHandler struct {
	handlers map[string]func(Context)
}

func (h *WebContextHandler) Route(method string, pattern string, handler func(Context)) {
	key := getKey(method, pattern)
	h.handlers[key] = handler
}

func (h *WebContextHandler) ServeHTTP(c Context) {
	key := c.Key()
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
