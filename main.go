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
	"kee"
	"log"
	"net/http"
	"time"
)

func onlyForV2() kee.HandlerFunc {
	return func(c *kee.Context) {
		// Start timer
		t := time.Now()
		// if a server error occurred
		c.Fail(500, "Internal Server Error")
		// Calculate resolution time
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

func main() {
	r := kee.New()
	r.Use(kee.Logger()) // global midlleware
	r.GET("/", func(c *kee.Context) {
		c.HTML(http.StatusOK, "<h1>Hello kee</h1>")
	})

	v2 := r.Group("/v2")
	v2.Use(onlyForV2()) // v2 group middleware
	{
		v2.GET("/hello/:name", func(c *kee.Context) {
			// expect /hello/keektutu
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})
	}

	r.Run(":9999")
}
