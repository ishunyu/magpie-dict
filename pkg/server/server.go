package main

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func main() {
	config := GetConfig()
	SetupLogger(config)
	index := GetIndex(config)

	http.Handle("/", http.FileServer(http.Dir(config.GetHtmlDir())))
	http.HandleFunc("/shows", RequestLogHandler(ShowsHandler(index)))
	http.HandleFunc("/search", RequestLogHandler(GetSearchHandler(index)))
	http.HandleFunc("/subs", RequestLogHandler(SubsHandler(index)))
	http.HandleFunc("/comparefiles", RequestLogHandler(CompareHandler(config.TempPath, config.ComparePath, config.CompareVenvPath)))

	port := config.GetPort()
	url := fmt.Sprintf("%s:%d", config.Hostname, port)
	log.Infof("Starting server on %v", url)
	err := http.ListenAndServe(url, nil)

	log.Error(err)
}
