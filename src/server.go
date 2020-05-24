package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type searchResult struct {
	Data []*searchRecord `json:"data"`
}

type searchRecord struct {
	Sub  *Record `json:"sub"`
	Pre  *Record `json:"pre"`
	Post *Record `json:"post"`
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (index *Index) getRecordContext(record *Record) *searchRecord {
	pre := index.GetRecord(record.ID - 1)
	post := index.GetRecord(record.ID + 1)

	return &searchRecord{Pre: pre, Sub: record, Post: post}
}

func handleSearch(w http.ResponseWriter, req *http.Request, index *Index) {
	start := time.Now()

	searchText := req.FormValue("searchText")
	logMessage := searchText

	records := index.Search(searchText)
	if records == nil {
		logMessage += " (No results)"
		fmt.Println(logMessage)
		return
	}

	records = records[0:min(len(records), 5)]
	result := searchResult{make([]*searchRecord, len(records))}
	for i, rec := range records {
		result.Data[i] = index.getRecordContext(rec)
	}

	data, _ := json.Marshal(result)
	fmt.Fprintf(w, string(data))

	elapsed := time.Since(start)
	logMessage += fmt.Sprintf(" %s", elapsed)

	fmt.Println(logMessage)
}

func getSearchHandler(index *Index) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		handleSearch(w, req, index)
	}
}

func getConfig() *Config {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("Missing configuration file argument.")
		os.Exit(1)
	}

	configPath := args[0]
	return GetConfig(configPath)
}

func main() {
	config := getConfig()
	index := GetIndex(config)

	http.Handle("/", http.FileServer(http.Dir(config.GetHtmlDir())))
	http.HandleFunc("/search", getSearchHandler(index))

	port := 8090
	fmt.Printf("Starting server on localhost:%v\n", port)
	http.ListenAndServe(fmt.Sprint("localhost:", port), nil)
}
