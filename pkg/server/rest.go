package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type searchResult struct {
	Data []*searchData `json:"data"`
}

type searchData struct {
	Show    string      `json:"show"`
	Episode string      `json:"episode"`
	Subs    *searchSubs `json:"subs"`
}

type searchSubs struct {
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

	hitsIds := index.Search(searchText)
	if hitsIds == nil {
		logMessage += " (No results)"
		fmt.Println(logMessage)
		hitsIds = make([][]int, 0)
	}

	hitsIds = hitsIds[0:Min(len(hitsIds), 5)]
	result := searchResult{make([]*searchData, len(hitsIds))}
	for i, ids := range hitsIds {
		result.Data[i] = retreiveResults(&index.Data, ids)
	}

	data, _ := json.Marshal(result)
	fmt.Fprintf(w, string(data))

	elapsed := time.Since(start)
	logMessage += fmt.Sprintf(" (%s)", elapsed)

	fmt.Println(logMessage)
}

func retreiveResults(data *Data, ids []int) *searchData {
	if ids == nil {
		return nil
	}

	showId, fileId, id := ids[0], ids[1], ids[2]
	show := data.Shows[showId]
	file := show.Files[fileId]
	subs := retreiveRecordContext(&file, id)

	return &searchData{show.Title, file.Name, subs}
}

func retreiveRecordContext(file *Showfile, id int) *searchSubs {
	pre := GetRecord(file, id-1)
	sub := GetRecord(file, id)
	post := GetRecord(file, id+1)

	return &searchSubs{Pre: pre, Sub: sub, Post: post}
}

func GetRecord(file *Showfile, id int) *Record {
	if id < 0 || id >= len(file.Records) {
		return nil
	}
	return &file.Records[id]
}
