package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type searchResponse struct {
	Data []*searchResponseData `json:"data"`
}

type searchResponseData struct {
	Show    string              `json:"show"`
	Episode string              `json:"episode"`
	Subs    *searchResponseSubs `json:"subs"`
}

type searchResponseSubs struct {
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

	searchResults := index.Search(searchText)
	if searchResults == nil {
		searchResults = make([]searchResult, 0)
	}

	logMessage += fmt.Sprintf(",%d", len(searchResults))

	searchResults = searchResults[0:Min(len(searchResults), 10)]
	response := searchResponse{make([]*searchResponseData, len(searchResults))}
	for i, result := range searchResults {
		response.Data[i] = retreiveResponse(&index.Data, &result)
	}

	data, _ := json.Marshal(response)
	fmt.Fprintf(w, string(data))

	elapsed := time.Since(start)
	logMessage += fmt.Sprintf(",%s", elapsed)

	fmt.Println(logMessage)
}

func retreiveResponse(data *Data, result *searchResult) *searchResponseData {
	if result == nil {
		return nil
	}

	show := data.Shows[result.showID]
	file := show.Files[result.fileID]
	subs := retreiveRecordContext(&file, result.subID)

	return &searchResponseData{show.Title, file.Name, subs}
}

func retreiveRecordContext(file *Showfile, id int) *searchResponseSubs {
	pre := GetRecord(file, id-1)
	sub := GetRecord(file, id)
	post := GetRecord(file, id+1)

	return &searchResponseSubs{Pre: pre, Sub: sub, Post: post}
}

func GetRecord(file *Showfile, id int) *Record {
	if id < 0 || id >= len(file.Records) {
		return nil
	}
	return &file.Records[id]
}
