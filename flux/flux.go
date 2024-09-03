package flux

import (
	"log"
	"net/http"
	"path"
	"strings"
	"text/template"
)

// 定义处理请求的回调函数
type HandlerFunc func(*Context)

// 代表组路由,拥有共同的前缀
type RouterGroup struct {
	prefix          string
	middlewares     []HandlerFunc //当前组的中间件
	parent          *RouterGroup
	engine          *Engine            // 所有的group共享一个engine实例
	renderTemplates *template.Template // 模板引擎
	renderFuncMap   template.FuncMap   // 模板函数
}

// 框架核心引擎,包含路由映射
type Engine struct {
	*RouterGroup
	router *router
	groups []*RouterGroup //所有路由组
}

// 实例化一个Engine
func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func (engine *Engine) SetRenderFuncMap(funcMap template.FuncMap) {
	engine.renderFuncMap = funcMap
}

func (engine *Engine) LoadTemplates(pattern string) {
	engine.renderTemplates = template.Must(template.New("").Funcs(engine.renderFuncMap).ParseGlob(pattern))
}

// 添加一个路由组
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

// 添加一条路由规则
func (group *RouterGroup) addRoute(method string, pattern string, handler HandlerFunc) {
	fullPattern := group.prefix + pattern
	log.Printf("Route %4s - %s", method, fullPattern)
	group.engine.router.addRoute(method, fullPattern, handler)
}

// 添加静态文件处理服务
func (group *RouterGroup) createStaticFileHandler(relativePath string, absoluteStream http.FileSystem) HandlerFunc {
	routePath := path.Join(group.prefix, relativePath)                         // 计算绝对路径,包含组前缀
	fileServer := http.StripPrefix(routePath, http.FileServer(absoluteStream)) // 创建文件服务
	return func(c *Context) {
		file := c.GetParam("filepath")                       // 获取文件名称
		if _, err := absoluteStream.Open(file); err != nil { //打开absolutePath/filepath的文件
			c.SetStatus(http.StatusNotFound)
			return
		}
		fileServer.ServeHTTP(c.Writer, c.Req) // 处理静态文件
	}
}

// 将静态文件服务注册到路由中
func (group *RouterGroup) MapFile(relativePath string, absolutePath string) {
	handler := group.createStaticFileHandler(relativePath, http.Dir(absolutePath)) // 创建文件服务处理函数
	urlPattern := path.Join(relativePath, "/*filepath")                            // 计算url路径
	group.GET(urlPattern, handler)                                                 // 注册GET请求
}

// 注册GET请求
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

// 注册POST请求
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

// 实现ServeHTTP接口,处理http请求,根据不同url路由到具体的处理函数
func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// 根据请求路径查找对应的中间件
	var middlewares []HandlerFunc
	for _, group := range e.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := newContext(w, req)
	c.engine = e
	c.middlewares = middlewares
	e.router.handle(c)
}

// 启动flux web框架
func (e *Engine) RunFlux(addr string) (err error) {
	return http.ListenAndServe(addr, e)
}
