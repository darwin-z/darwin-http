package flux

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

// 代表一个http上下文
type Context struct {
	Writer      http.ResponseWriter //响应控制器
	Req         *http.Request       //请求信息
	Path        string              //请求路径
	Method      string              //请求方法
	Params      map[string]string   //请求路径参数
	StatusCode  int                 //响应状态码
	middlewares []HandlerFunc       //中间件
	index       int                 //记录当前执行到第几个中间件
}

// 创建一个新的http上下文
func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
		index:  -1,
	}
}

// 执行中间件
func (c *Context) Next() {
	c.index++
	l := len(c.middlewares)
	// 依次执行中间件,避免中间件没有调用Next()导致后续中间件不执行
	for ; c.index < l; c.index++ {
		c.middlewares[c.index](c)
	}
}

// 获取表单参数
func (c *Context) GetForm(key string) string {
	return c.Req.FormValue(key)
}

// 获取query参数
func (c *Context) GetQuery(key string) string {
	return c.Req.URL.Query().Get(key)
}

// 获取请求路径参数
func (c *Context) GetParam(key string) string {
	return c.Params[key]
}

// 设置响应状态码
func (c *Context) SetStatus(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

// 设置响应头
func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

// 设置失败响应
func (c *Context) SetFail(code int, message string) {
	c.index = len(c.middlewares) //跳过后续中间件
	c.SetJSON(code, H{"message": message})
}

// 设置字符串类型的响应
func (c *Context) SetString(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.SetStatus(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

// 设置json类型的响应
func (c *Context) SetJSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.SetStatus(code)
	// 使用json.NewEncoder(c.Writer).Encode(obj)将obj序列化到c.Writer中
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

// 设置html类型的响应
func (c *Context) SetHTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.SetStatus(code)
	c.Writer.Write([]byte(html))
}

// 设置字节类型的响应
func (c *Context) SetData(code int, data []byte) {
	c.SetStatus(code)
	c.Writer.Write(data)
}
