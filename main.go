package main

import (
	"fmt"
	"kee"
	"net/http"
)

func main() {
	kee := kee.New()

	kee.GET("/hello", indexHandler)
	kee.GET("/header", helloHandler)

	http.ListenAndServe(":7812", kee)
}

// handler echoes r.URL.Path
func indexHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
}

// handler echoes r.URL.Header
func helloHandler(w http.ResponseWriter, req *http.Request) {
	for k, v := range req.Header {
		fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
	}
}
