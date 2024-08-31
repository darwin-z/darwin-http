package main

import (
	"flux"
	"flux/middlewares"
	"fmt"
	"log"
	"time"
)

func main() {
	server := flux.New()
	server.Use(middlewares.Logger()) //注册全局中间件
	server.GET("/", func(c *flux.Context) {
		c.SetHTML(200, "<h1>Hello, Flux!</h1>")
	})

	v1 := server.Group("/v1")
	v1.Use(func(c *flux.Context) {
		t := time.Now()
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
	})
	v1.GET("/hello/:name", func(c *flux.Context) {
		c.SetHTML(200, fmt.Sprintf("<h1>Hello, %s!</h1>", c.GetParam("name")))
	})

	server.RunFlux(":9999")
}
