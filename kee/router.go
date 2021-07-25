package kee

import (
	"net/http"
	"strings"
)

type router struct {
	roots    map[string]*node
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, part := range vs {
		if part != "" {
			parts = append(parts, part)
			if strings.HasPrefix(part, "*") {
				break
			}
		}
	}
	return parts
}

func (r router) addRoute(method string, pattern string, handler HandlerFunc) {

	if _, ok := r.roots[method]; !ok {
		r.roots[method] = &node{}
	}

	root := r.roots[method]

	parts := parsePattern(pattern)
	root.insert(pattern, parts, 0)

	key := method + "-" + pattern
	r.handlers[key] = handler

}

func (r router) getRoute(method string, path string) (*node, map[string]string) {

	searchParts := parsePattern(path)
	params := make(map[string]string)

	if _, ok := r.roots[method]; !ok {
		return nil, nil
	}
	root := r.roots[method]

	n := root.search(searchParts, 0)

	if n != nil {
		parts := parsePattern(n.pattern)

		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}
	return nil, nil
}

func (r router) getRoutes(method string) []*node {
	if _, ok := r.roots[method]; !ok {
		return nil
	}
	root := r.roots[method]
	nodes := []*node{}

	root.travel(&nodes)
	return nodes
}

func (r router) handle(c *Context) {

	n, params := r.getRoute(c.Method, c.Path)

	if n != nil {
		c.Params = params
		key := c.Method + "-" + n.pattern
		handler := r.handlers[key]
		handler(c)

	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s \n", c.Path)
	}
}
