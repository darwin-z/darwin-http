package main

import (
	"flux"
	"flux/middlewares"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	server := flux.New()
	server.Use(middlewares.Logger()) //注册全局中间件
	server.Use(middlewares.Recovery())
	//设置模版渲染函数
	server.LoadTemplates("./template/*") //加载模版文件
	server.MapFile("/file", "./static/") //映射静态文件

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
	v1.GET("/time", func(c *flux.Context) {
		c.SetHTMLTemplate(200, "hello.html", flux.Object{
			"date": time.Now(),
		})
	})
	v1.GET("/panic", func(c *flux.Context) {
		names := []string{"geektutu"}
		c.SetString(http.StatusOK, names[100])
	})

	server.RunFlux(":9999")
}
