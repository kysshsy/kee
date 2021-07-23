package kee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Engine struct {
	Router *Router
}

func New() *Engine {

	ret := new(Engine)
	ret.Router = NewRouter()
	return ret
}

func (e *Engine) ServeHTTP(respWriter http.ResponseWriter, req *http.Request) {
	//
	context := NewContext(respWriter, req)

	e.Router.Handle(context)

}

func (e *Engine) GET(url string, handler func(*Context)) {
	e.Router.AddRoute("GET", url, handler)
}

func (e *Engine) POST(url string, handler func(*Context)) {
	e.Router.AddRoute("POST", url, handler)
}

type H map[string]interface{}

type Context struct {
	Writer http.ResponseWriter
	Req    *http.Request

	Url    string
	Method string

	StatusCode int
}

func NewContext(writer http.ResponseWriter, req *http.Request) *Context {
	return &Context{Writer: writer, Req: req, Url: req.URL.Path, Method: req.Method}
}

// getter

func (c *Context) Query(key string) string {
	values := c.Req.Header.Values(key)

	if len(values) == 0 {
		return ""
	}

	return values[0]
}

func (c *Context) QueryArray(key string) []string {
	return c.Req.Header.Values(key)
}

func (c *Context) GetForm(key string) string {
	return c.Req.Form.Get(key)
}

// setters

func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

func (c *Context) Error(code int, emsg string) {
	http.Error(c.Writer, emsg, code)
}

func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Contnt-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

func (c *Context) Json(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")

	jsonByte, err := json.Marshal(obj)
	if err != nil {
		c.Error(http.StatusInternalServerError, err.Error())
	}
	c.Status(code)
	c.Writer.Write(jsonByte)
}

func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	c.Writer.Write([]byte(html))
}

func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}
