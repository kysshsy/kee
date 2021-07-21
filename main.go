package main

import (
	"kee"
	"net/http"
)

func main() {
	engine := kee.New()

	engine.GET("/hello", indexHandler)

	http.ListenAndServe(":7812", engine)
}

// handler echoes r.URL.Path
func indexHandler(c *kee.Context) {
	c.Json(http.StatusOK, kee.H{
		"user": "kyss",
		"hsy":  "douy",
	})
}
