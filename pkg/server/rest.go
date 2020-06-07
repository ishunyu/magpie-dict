package main

import (
	"encoding/json"
	"fmt"
	"net/http"
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

func GetSearchHandler(index *Index) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		handleSearch(w, req, index)
	}
}

func handleSearch(w http.ResponseWriter, req *http.Request, index *Index) {
	start := time.Now()

	searchText := req.FormValue("searchText")
	logMessage := searchText

	records := index.Search(searchText)
	if records == nil {
		logMessage += " (No results)"
		fmt.Println(logMessage)
		records = make([]*Record, 0)
	}

	records = records[0:Min(len(records), 5)]
	result := searchResult{make([]*searchRecord, len(records))}
	for i, rec := range records {
		result.Data[i] = index.getRecordContext(rec)
	}

	data, _ := json.Marshal(result)
	fmt.Fprintf(w, string(data))

	elapsed := time.Since(start)
	logMessage += fmt.Sprintf(" (%s)", elapsed)

	fmt.Println(logMessage)
}

func (index *Index) getRecordContext(record *Record) *searchRecord {
	pre := index.GetRecord(record.ID - 1)
	post := index.GetRecord(record.ID + 1)

	return &searchRecord{Pre: pre, Sub: record, Post: post}
}
