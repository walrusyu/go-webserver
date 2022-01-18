package webserver

import (
	"net/http"
	"strings"
)

type Handler interface {
	Route(method string, pattern string, handler func(Context))
	ServeHTTP(c Context)
}

var _ Handler = &Handler_BaseOnMap{}

type Handler_BaseOnMap struct {
	handlers map[string]func(Context)
}

func (h *Handler_BaseOnMap) Route(method string, pattern string, handler func(Context)) {
	key := getKey(method, pattern)
	h.handlers[key] = handler
}

func (h *Handler_BaseOnMap) ServeHTTP(c Context) {
	key := getKey(c.R().Method, c.R().URL.Path)
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
	return &Handler_BaseOnMap{handlers: make(map[string]func(Context))}
}

func NewHandler2() Handler {
	return &Handler_BaseOnTree{
		root: &Node{
			children: map[string]*Node{},
		},
	}
}

var _ Handler = &Handler_BaseOnTree{}

type Handler_BaseOnTree struct {
	root *Node
}

func (h *Handler_BaseOnTree) Route(method string, pattern string, handler func(Context)) {
	paths := getPaths(method, pattern)
	h.root.CreateNode(paths, handler)
}

func (h *Handler_BaseOnTree) ServeHTTP(c Context) {
	paths := getPaths(c.R().Method, c.R().URL.Path)
	n := h.root.FindNode(paths)
	if n == nil {
		c.WriteCode(http.StatusNotFound)
	} else {
		n.handler(c)
	}
}

type Node struct {
	path     string
	children map[string]*Node
	handler  func(Context)
}

func (n *Node) FindNode(paths []string) *Node {
	curr := n
	for _, path := range paths {
		if curr != nil && curr.children[path] != nil {
			curr = curr.children[path]
		} else {
			return curr.children["*"]
		}
	}

	return curr
}

func (n *Node) CreateNode(paths []string, handler func(Context)) *Node {
	if n == nil {
		return nil
	}
	curr := n
	for _, path := range paths {
		if curr.children[path] == nil {
			curr.children[path] = &Node{
				path:     path,
				children: make(map[string]*Node),
			}
		}
		curr = curr.children[path]
	}
	curr.handler = handler
	return curr
}

func getPaths(method string, pattern string) []string {
	paths := []string{method}
	pattern = strings.Trim(pattern, "/")
	paths = append(paths, strings.Split(pattern, "/")...)
	return paths
}
