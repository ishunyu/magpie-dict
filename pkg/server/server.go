package main

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func main() {
	setupLogger()
	config := GetConfig()
	index := GetIndex(config)

	setupRequestLogger(config)
	http.Handle("/", http.FileServer(http.Dir(config.GetHtmlDir())))
	http.HandleFunc("/shows", RequestLogHandler(ShowsHandler(index), config))
	http.HandleFunc("/search", RequestLogHandler(GetSearchHandler(index), config))
	http.HandleFunc("/subs", RequestLogHandler(SubsHandler(index), config))
	http.HandleFunc("/comparefiles", RequestLogHandler(CompareHandler(config.TempPath, config.ComparePath, config.CompareVenvPath), config))

	port := config.GetPort()
	url := fmt.Sprintf("%s:%d", config.Hostname, port)
	log.Infof("Starting server on %v", url)
	err := http.ListenAndServe(url, nil)

	log.Error(err)
}
