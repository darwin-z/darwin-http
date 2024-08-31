package middlewares

import (
	"flux"
	"log"
	"time"
)

// Logger 中间件,用于记录请求执行时间
func Logger() flux.HandlerFunc {
	return func(c *flux.Context) {
		t := time.Now()
		c.Next()
		log.Printf("[Logger] %d %s in %v", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}
