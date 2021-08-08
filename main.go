package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/", index)
	http.HandleFunc("/foo", foo)
	log.Fatal(http.ListenAndServe(":8090", nil))
}

func index(resp http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(resp, "URL.path = %q\n", req.URL.Path)
}

func foo(resp http.ResponseWriter, req *http.Request) {
	for k, v := range req.Header {
		fmt.Fprintf(resp, "Header[%q] = %q\n", k, v)
	}
}
