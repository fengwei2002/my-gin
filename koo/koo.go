package koo

import (
	"log"
	"net/http"
	"strings"
)

// HandlerFunc 是路由处理函数类型的一个缩写
// type HandlerFunc func(http.ResponseWriter, *http.Request)
type HandlerFunc func(c *Context) // 修改为使用上下文

// Engine 实现了 ServerHTTP 方法，可以作为 ListenAndServer 的第二个参数使用
// - Engine.router 将 string 类型的路由映射到一个路由处理函数
type Engine struct {
	*RouterGroup
	router *router
	groups []*RouterGroup // store all groups
}

// RouterGroup 支持对分组的路由应用一些规则
type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc // 支持中间件
	parent      *RouterGroup  // support nesting
	engine      *Engine       // 所有的 groups 共享同一个 engine 实例
}

// New 是 koo.Engine 的构造函数
func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

// Group 函数，接收一个 prefix 字符串，创建一个新的 RouterGroup
// 新的 prefix 是传入的 prefix 加上 group 本来就拥有的 prefix
// 新创建的 newGroup 和 group 使用的 engine 是同一个 engine，在 engine 的属性中可以访问到
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

// addRoute 将一个路由的信息存储到 engine 的 router 中
// 第一个参数是使用的方法，第二个参数是具体的路由，第三个参数是路由处理函数
//func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
//	log.Printf("Route %4s - %s", method, pattern)
//	engine.router.addRoute(method, pattern, handler)
//}
func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	// 使用 group 中的信息构建一个路由之后
	// 然后调用 router 中的 addRoute 将想要添加的 路由添加到 tire 中
	group.engine.router.addRoute(method, pattern, handler)
}

// GET 使用 GET 方法调用 addRoute 不用在参数列表中声明 GET 而是使用 GET 方法
//func (engine *Engine) GET(pattern string, handler HandlerFunc) {
//	engine.addRoute("GET", pattern, handler)
//}
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

// POST 使用 POST 方法调用 addRoute 不用在参数列表中声明 POST 而是使用 POST 方法
//func (engine *Engine) POST(pattern string, handler HandlerFunc) {
//	engine.addRoute("POST", pattern, handler)
//}
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

// Use 调用 Use，将传入的 HandlerFunc 全部添加到这个 group 的信息中
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

// Run 使用 run 接口，将传入的 addr 使用  http ListenAndServe 运行，第二个参数的 engine 已经实现 ServeHTTP 方法
// 所有的路由请求交给 engine 处理
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

//
//func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
//	c := newContext(w, req)
//	// 创建一个上下文
//	engine.router.handle(c)
//	// 使用 router.handle 处理这个上下文
//}

// ServeHTTP 当我们接收到一个具体请求时，要判断该请求适用于哪些中间件，
// 在这里简单通过 URL 的前缀来判断。
// 得到中间件列表后，赋值给 c.handlers。
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc
	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}

	c := newContext(w, req)
	c.handlers = middlewares
	engine.router.handle(c) // 使用 router.handle 处理这个上下文
}
