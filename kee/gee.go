package kee

import (
	"log"
	"net/http"
	"strings"
	"time"
)

type H map[string]interface{}
type HandlerFunc func(*Context)

type RouterGroup struct {
	prefix     string
	middleware []HandlerFunc
	engine     *Engine
}

type Engine struct {
	*RouterGroup
	router *router
	groups []*RouterGroup
}

func New() *Engine {
	engine := Engine{
		router: newRouter(),
	}
	engine.RouterGroup = &RouterGroup{engine: &engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return &engine
}

func (e *Engine) Run(addr string) {
	err := http.ListenAndServe(addr, e)
	if err != nil {
		return
	}
}

func (e *Engine) ServeHTTP(respWriter http.ResponseWriter, req *http.Request) {
	context := NewContext(respWriter, req)
	path := context.Path

	context.handlers = append(context.handlers, e.middleware...)
	for _, group := range e.groups {
		strings.HasPrefix(path, group.prefix)
		context.handlers = append(context.handlers, group.middleware...)
	}

	e.router.handle(context)

}

func (group *RouterGroup) Group(path string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + path,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middleware = append(group.middleware, middlewares...)
}

func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

func Logger() HandlerFunc {
	return func(c *Context) {
		// Start timer
		t := time.Now()
		// Process request
		c.Next()
		// Calculate resolution time
		log.Printf("[%d] %s in %v", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}
