package main

import (
	"github.com/walrusyu/go-webserve/webserver"
	"net/http"
)

func main() {
	server := webserver.NewServer(webserver.TimeFilterBuilder)
	server.Route(http.MethodGet, "/", home)
	err := server.Start(":8080")
	if err != nil {
		panic(err)
	}
}

func home(c webserver.Context) {
	c.WriteStr("hello, world")
}
