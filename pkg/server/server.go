package main

import (
	"fmt"
	"net/http"
)

func main() {
	config := GetConfig()
	index := GetIndex(config)

	http.Handle("/", http.FileServer(http.Dir(config.GetHtmlDir())))
	http.HandleFunc("/search", GetSearchHandler(index))

	port := 8090
	fmt.Printf("Starting server on localhost:%v\n", port)
	http.ListenAndServe(fmt.Sprint("localhost:", port), nil)
}
