package main

import (
	"github.com/walrusyu/go-webserve/webserver"
	"net/http"
)

func main() {
	server := webserver.NewServer(webserver.TimeFilterBuilder)
	server.Route(http.MethodGet, "/", home)
	server.Route(http.MethodGet, "/user", user)
	server.Route(http.MethodGet, "/*", whatever)
	err := server.Start(":8080")
	if err != nil {
		panic(err)
	}
}

func home(c webserver.Context) {
	c.WriteStr("hello, world")
}

func user(c webserver.Context) {
	c.WriteJson(User{Name: "ywf", Age: 18})
}

type User struct {
	Name string `name`
	Age  int    `age`
}

func whatever(c webserver.Context) {
	c.WriteStr("whatever")
}
