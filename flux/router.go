package flux

import (
	"net/http"
	"strings"
)

// 表示一个路由处理器
type router struct {
	routes   map[string]*node       // 路由树 map[请求方法]node
	handlers map[string]HandlerFunc // 处理函数 map[请求方法-请求路径]处理函数
}

func newRouter() *router {
	return &router{
		routes:   make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// 添加一条路由规则
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	parts := parseUrlPattern(pattern)
	key := method + "-" + pattern
	if _, ok := r.routes[method]; !ok {
		r.routes[method] = &node{}
	}
	r.routes[method].insert(pattern, parts, 0)
	r.handlers[key] = handler
}

// 解析url路径,返回路径中的各个部分
func (r *router) getRoute(method string, pattern string) (*node, map[string]string) {
	// 解析url路径
	parsedPattern := parseUrlPattern(pattern)
	//解析url中param参数
	params := make(map[string]string)
	//获取路由树
	route, ok := r.routes[method]
	if !ok {
		return nil, nil
	}
	//查找匹配的节点
	matchedNode := route.search(parsedPattern, 0)
	if matchedNode == nil {
		return nil, nil
	}
	//匹配到了节点,解析params参数
	parts := parseUrlPattern(matchedNode.urlPattern)
	for i, part := range parts {
		if part[0] == ':' {
			params[part[1:]] = parsedPattern[i]
		}
		if part[0] == '*' && len(part) > 1 {
			params[part[1:]] = strings.Join(parsedPattern[i:], "/")
			break
		}
	}
	return matchedNode, params
}

// 根据上下文处理请求
func (r *router) handle(c *Context) {
	node, params := r.getRoute(c.Method, c.Path)
	if node == nil {
		c.SetString(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
		return
	}
	c.Params = params
	key := c.Method + "-" + node.urlPattern
	c.middlewares = append(c.middlewares, r.handlers[key]) // 添加处理函数
	c.Next()                                               // 执行中间件
}
