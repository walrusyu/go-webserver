package main

import (
	"github.com/walrusyu/go-webserve/webserver"
	"net/http"
)

func main() {
	server := webserver.NewServer()
	server.Route(http.MethodGet, "/", home)
	err := server.Start(":8080")
	if err != nil {
		panic(err)
	}
}

func home(c webserver.Context) {
	c.WriteStr("hello, world")
}
