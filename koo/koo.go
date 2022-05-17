package koo

import (
	"log"
	"net/http"
)

// HandlerFunc 是路由处理函数类型的一个缩写
//type HandlerFunc func(http.ResponseWriter, *http.Request)
type HandlerFunc func(c *Context) // 修改为使用上下文

// Engine 实现了 ServerHTTP 方法，可以作为 ListenAndServer 的第二个参数使用
// - Engine.router 将 string 类型的路由映射到一个路由处理函数
type Engine struct {
	router *router
}

// New 是 koo.Engine 的构造函数
func New() *Engine {
	return &Engine{
		router: newRouter(),
	}
}

// addRoute 将一个路由的信息存储到 engine 的 router 中
// 第一个参数是使用的方法，第二个参数是具体的路由，第三个参数是路由处理函数
func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	log.Printf("Route %4s - %s", method, pattern)
	engine.router.addRoute(method, pattern, handler)
}

// GET 使用 GET 方法调用 addRoute 不用在参数列表中声明 GET 而是使用 GET 方法
func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}

// POST 使用 POST 方法调用 addRoute 不用在参数列表中声明 POST 而是使用 POST 方法
func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}

// Run 使用 run 接口，将传入的 addr 使用  http ListenAndServe 运行，第二个参数的 engine 已经实现 ServeHTTP 方法
// 所有的路由请求交给 engine 处理
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	// 创建一个上下文
	engine.router.handle(c)
	// 使用 router.handle 处理这个上下文
}
