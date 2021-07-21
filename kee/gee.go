package kee

import (
	"encoding/json"
	"net/http"
)

type Engine struct {
	Router Router
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

func (e *Engine) GET(url string, handlerFunc func(*Context)) {
	e.Router.AddRoute("GET", url, handlerFunc)
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

func (c *Context) Json(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")

	jsonByte, err := json.Marshal(obj)
	if err != nil {
		c.Error(http.StatusInternalServerError, err.Error())
	}
	c.Status(code)
	c.Writer.Write(jsonByte)
}
