package koo

import (
	"log"
	"net/http"
)

// 将 router 部分进行独立，方便在 router 中太你家功能
type router struct {
	handlers map[string]HandlerFunc
}


// newRouter 是 router 的构造函数，返回一个 router
func newRouter() *router {
	return &router{
		handlers: make(map[string]HandlerFunc),
	}
}

// addRoute 提供接口，method， pattern 和 handler 参数
// 将信息打印到 log 中，同时，将路由信息保存到 router 的 map[string]HandlerFunc 中
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	log.Printf("Route %4s - %s\n", method, pattern)
	key := method + "-" + pattern
	r.handlers[key] = handler
}

func (r *router) handle(c *Context) {
	key := c.Method + "-" + c.Path
	if handler, ok := r.handlers[key]; ok == true {
		handler(c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}

