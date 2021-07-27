package main

/*
(1)
$ curl -i http://localhost:9999/
HTTP/1.1 200 OK
Date: Mon, 12 Aug 2019 16:52:52 GMT
Content-Length: 18
Content-Type: text/html; charset=utf-8
<h1>Hello kee</h1>
(2)
$ curl "http://localhost:9999/hello?name=keektutu"
hello keektutu, you're at /hello
(3)
$ curl "http://localhost:9999/hello/keektutu"
hello keektutu, you're at /hello/keektutu
(4)
$ curl "http://localhost:9999/assets/css/keektutu.css"
{"filepath":"css/keektutu.css"}
(5)
$ curl "http://localhost:9999/xxx"
404 NOT FOUND: /xxx
*/

import (
	"net/http"

	"kee"
)

func main() {
	r := kee.New()
	r.GET("/", func(c *kee.Context) {
		c.HTML(http.StatusOK, "<h1>Hello kee</h1>")
	})

	r.GET("/hello", func(c *kee.Context) {
		// expect /hello?name=keektutu
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})

	r.GET("/hello/:name", func(c *kee.Context) {
		// expect /hello/keektutu
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
	})

	r.GET("/assets/*filepath", func(c *kee.Context) {
		c.JSON(http.StatusOK, kee.H{"filepath": c.Param("filepath")})
	})

	r.Run(":9999")
}
