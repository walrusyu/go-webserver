package main

import (
	"fmt"
	"github.com/walrusyu/go-webserve/webserver"
	"net/http"
)

func main() {
	shutdown := webserver.NewGracefulShutdown()

	fmt.Println("starting")
	server := webserver.NewServer(webserver.TimeFilterBuilder, shutdown.ShutdownFilterBuilder)
	server.Route(http.MethodGet, "/", home)
	server.Route(http.MethodGet, "/user", user)
	server.Route(http.MethodGet, "/*", whatever)
	go func() {
		err := server.Start(":8080")
		if err != nil {
			panic(err)
		}
	}()

	webserver.WaitForShutdown(webserver.BuildShutdownServerHook(server), shutdown.RejectNewRequestAndWaiting)
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
