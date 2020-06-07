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

	port := config.GetPort()
	url := fmt.Sprintf("%s:%d", config.Hostname, port)
	fmt.Printf("Starting server on %v\n", url)
	http.ListenAndServe(url, nil)
}
