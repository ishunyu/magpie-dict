package main

import (
	"fmt"
	"net/http"
)

func main() {
	config := GetConfig()
	index := GetIndex(config)

	http.Handle("/", http.FileServer(http.Dir(config.GetHtmlDir())))
	http.HandleFunc("/shows", ShowsHandler(index))
	http.HandleFunc("/search", GetSearchHandler(index))
	http.HandleFunc("/subs", SubsHandler(index))
	http.HandleFunc("/comparefiles", CompareHandler(config.TempPath, config.ComparePath))

	port := config.GetPort()
	url := fmt.Sprintf("%s:%d", config.Hostname, port)
	fmt.Printf("Starting server on %v\n", url)
	http.ListenAndServe(url, nil)
}
