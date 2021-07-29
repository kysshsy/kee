package kee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Context struct {
	// origin object
	Writer http.ResponseWriter
	Req    *http.Request

	// request info
	Path   string
	Method string
	Params map[string]string

	// response info
	StatusCode int

	// middleware
	handlers []HandlerFunc
	index    int
}

func NewContext(writer http.ResponseWriter, req *http.Request) *Context {
	return &Context{Writer: writer, Req: req, Path: req.URL.Path, Method: req.Method, index: -1}
}

// getter

func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}

func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)

}

func (c *Context) GetForm(key string) string {
	return c.Req.FormValue(key)
}

func (c *Context) PostForm(key string) string {
	return c.Req.PostFormValue(key)
}

// setters

func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

func (c *Context) Error(code int, err interface{}) {
	if errs, ok := err.(error); ok {
		http.Error(c.Writer, errs.Error(), code)
	} else if msg, ok := err.(string); ok {
		http.Error(c.Writer, msg, code)
	}
}

func (c *Context) Fail(code int, err string) {
	c.index = len(c.handlers)
	c.JSON(code, H{"message": err})
}

func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

func (c *Context) JSON(code int, obj interface{}) {
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

func (c *Context) Next() {
	c.index++
	s := len(c.handlers)
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}
