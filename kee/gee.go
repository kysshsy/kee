package kee

import (
	"fmt"
	"log"
	"net/http"
	"path"
	"runtime"
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

func Default() *Engine {
	engine := New()
	engine.Use(Logger(), Recovery())
	return engine
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

// print stack trace for debug
func trace(message string) string {
	var pcs [32]uintptr
	n := runtime.Callers(3, pcs[:]) // skip first 3 caller

	var str strings.Builder
	str.WriteString(message + "\nTraceback:")
	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return str.String()
}

func Recovery() HandlerFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				message := fmt.Sprintf("%s", err)
				log.Printf("%s\n\n", trace(message))
				c.Fail(http.StatusInternalServerError, "Internal Server Error")
			}
		}()
		c.Next() // 很重要 不然defer回提前执行 不在middlerware的下半部分 而是上半部分
	}
}

func createStaticHandler(prefix string, rootPath string) HandlerFunc {
	// fs
	filesystem := http.Dir(rootPath)
	handler := http.StripPrefix(prefix, http.FileServer(filesystem))

	return func(c *Context) {

		filename := c.Param("filename")

		if _, err := filesystem.Open(filename); err != nil {
			c.Fail(500, "error")
			return
		}

		handler.ServeHTTP(c.Writer, c.Req)
	}
}

func (group *RouterGroup) Static(relativePath, rootPath string) {

	urlPrefix := path.Join(group.prefix, relativePath)

	handler := createStaticHandler(urlPrefix, rootPath)
	// 注册动态路由
	pattern := urlPrefix + "/*filename"
	group.GET(pattern, handler)

}
