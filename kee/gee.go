package kee

import (
	"fmt"
	"net/http"
)

type Engine struct {
	router map[string]http.HandlerFunc
}

func New() *Engine {

	ret := new(Engine)
	ret.router = make(map[string]http.HandlerFunc)
	return ret
}

func (e *Engine) ServeHTTP(respWriter http.ResponseWriter, req *http.Request) {
	//
	key := req.Method + "-" + req.URL.Path

	if value, ok := e.router[key]; ok {
		value(respWriter, req)
	} else {
		fmt.Fprintf(respWriter, "404 NOT FOUND: %s\n", req.URL)
	}
}

func (e *Engine) GET(url string, handlerFunc http.HandlerFunc) {
	pattern := "GET"
	key := pattern + "-" + url

	e.router[key] = handlerFunc
}
