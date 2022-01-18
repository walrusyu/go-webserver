package webserver

import (
	"encoding/json"
	"io"
	"net/http"
)

type Context interface {
	ReadJson(obj interface{}) error
	WriteJson(object interface{}) error
	WriteStr(str string) error
	WriteCode(code int)
	W() http.ResponseWriter
	R() *http.Request
}

var _ Context = &WebContext{}

type WebContext struct {
	w http.ResponseWriter
	r *http.Request
}

func (c *WebContext) ReadJson(obj interface{}) error {
	b, err := io.ReadAll(c.r.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, obj)
	return err
}

func (c *WebContext) WriteJson(object interface{}) error {
	b, err := json.Marshal(object)
	if err != nil {
		return err
	}

	_, err = c.w.Write(b)
	return err
}

func (c *WebContext) WriteStr(str string) error {
	_, err := c.w.Write([]byte(str))
	return err
}

func (c *WebContext) WriteCode(code int) {
	c.w.WriteHeader(code)
}

func (c *WebContext) W() http.ResponseWriter {
	return c.w
}
func (c *WebContext) R() *http.Request {
	return c.r
}

func NewContext(w http.ResponseWriter, r *http.Request) Context {
	return &WebContext{w: w, r: r}
}
