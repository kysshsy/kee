package kee

import (
	"log"
	"net/http"
)

type HandlerFunc func(*Context)

type RouterGroup struct {
	prefix     string
	middleware []HandlerFunc
	engine     *Engine
}

type Engine struct {
	*RouterGroup
	router *router
	groups []RouterGroup
}

func New() *Engine {
	engine := Engine{
		router: newRouter(),
	}
	engine.RouterGroup = &RouterGroup{engine: &engine}
	// TODO: groups 置nil
	return &engine
}

func (e *Engine) Run(addr string) {
	err := http.ListenAndServe(addr, e)
	if err != nil {
		return
	}
}

func (e *Engine) ServeHTTP(respWriter http.ResponseWriter, req *http.Request) {
	//
	context := NewContext(respWriter, req)

	e.router.handle(context)

}

func (group *RouterGroup) Group(path string) *RouterGroup {
	engine := group.engine
	newGroup := RouterGroup{
		prefix: group.prefix + path,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return &newGroup
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

type H map[string]interface{}
