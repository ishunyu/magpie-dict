package main

import (
	"fmt"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
)

func main() {
	config := GetConfig()
	SetupLogger(config)
	index, err := NewIndex(config.DataPath, config.IndexPath)
	if err != nil {
		log.Error("Initializing index was unsuccessful.", err)
		os.Exit(1)
	}

	http.Handle("/", http.FileServer(http.Dir(config.GetHtmlDir())))
	http.HandleFunc("/shows", RequestLogHandler(ShowsHandler(index)))
	http.HandleFunc("/search", RequestLogHandler(GetSearchHandler(index)))
	http.HandleFunc("/subs", RequestLogHandler(SubsHandler(index)))
	http.HandleFunc("/comparefiles", RequestLogHandler(CompareHandler(config.TempPath, config.ComparePath, config.CompareVenvPath)))

	port := config.GetPort()
	url := fmt.Sprintf("%s:%d", config.Hostname, port)
	log.Infof("Starting server on %v", url)
	err = http.ListenAndServe(url, nil)

	log.Error(err)
}
