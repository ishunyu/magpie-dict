package main

import (
	"fmt"
	"net/http"
)

func hello(w http.ResponseWriter, req *http.Request) {

	fmt.Fprint(w, "Hello World!")
}

func main() {
	http.HandleFunc("/hello", hello)

	port := 8090
	fmt.Printf("Starting server on port %v...", port)
	http.ListenAndServe(fmt.Sprint(":", port), nil)
}
