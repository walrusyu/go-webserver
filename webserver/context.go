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
	Key() string
}

var _ Context = &WebContext{}

type WebContext struct {
	W http.ResponseWriter
	R *http.Request
}

func (c *WebContext) ReadJson(obj interface{}) error {
	b, err := io.ReadAll(c.R.Body)
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

	_, err = c.W.Write(b)
	return err
}

func (c *WebContext) WriteStr(str string) error {
	_, err := c.W.Write([]byte(str))
	return err
}

func (c *WebContext) WriteCode(code int) {
	c.W.WriteHeader(code)
}

func (c *WebContext) Key() string {
	return getKey(c.R.Method, c.R.URL.Path)

}
