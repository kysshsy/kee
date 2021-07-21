package kee

import "net/http"

type Router struct {
	handlers map[string]func(c *Context)
}

func NewRouter() Router {
	return Router{handlers: make(map[string]func(c *Context))}
}

func (r Router) AddRoute(method string, url string, handler func(c *Context)) {
	key := method + "-" + url

	r.handlers[key] = handler

}

func (r Router) Handle(c *Context) {
	// 搜索 handler 进行context转发
	key := c.Method + "-" + c.Url

	if handler, ok := r.handlers[key]; ok {
		handler(c)
		return
	}

	// 错误处理：没有找到handler
	c.Error(http.StatusNotFound, "")
}
