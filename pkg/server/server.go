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

	http.Handle("/", http.FileServer(http.Dir(config.GetHtmlDir())))
	http.HandleFunc("/shows", ShowsHandler(index))
	http.HandleFunc("/search", GetSearchHandler(index))
	http.HandleFunc("/subs", SubsHandler(index))
	http.HandleFunc("/comparefiles", CompareHandler(config.TempPath, config.ComparePath, config.CompareVenvPath))

	port := config.GetPort()
	url := fmt.Sprintf("%s:%d", config.Hostname, port)
	log.Infof("Starting server on %v", url)
	err := http.ListenAndServe(url, nil)

	log.Error(err)
}
