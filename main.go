package main

import (
	"fmt"
	"log"
	"net/http"
)

type Engine struct{}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
	case "/foo":
		for k, v := range r.Header {
			fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
		}
	default:
		fmt.Fprintf(w, "404 NOT FOUND: %s\n", r.URL)
	}
}

func main() {

	//http.HandleFunc("/", index)
	//http.HandleFunc("/foo", foo)
	engine := new(Engine)
	log.Fatal(http.ListenAndServe(":8090", engine))
}

func index(resp http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(resp, "URL.path = %q\n", req.URL.Path)
}

func foo(resp http.ResponseWriter, req *http.Request) {
	for k, v := range req.Header {
		fmt.Fprintf(resp, "Header[%q] = %q\n", k, v)
	}
}
