package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type searchResponse struct {
	Data []*searchResponseData `json:"data"`
}

type searchResponseData struct {
	Show    string    `json:"show"`
	Episode string    `json:"episode"`
	Subs    []*Record `json:"subs"`
}

func GetSearchHandler(index *Index) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		handleSearch(w, req, index)
	}
}

func SubsHandler(index *Index) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		handleSubs(w, req, index)
	}
}

func handleSearch(w http.ResponseWriter, req *http.Request, index *Index) {
	start := time.Now()

	searchText := req.FormValue("searchText")
	showID := req.FormValue("showID")

	logMessage := "/search " + searchText
	if showID != "" {
		logMessage += " (" + showID + ")"
	}

	searchResults := index.Search(searchText, showID)
	if searchResults == nil {
		searchResults = make([]*recordID, 0)
	}

	logMessage += fmt.Sprintf(",%d", len(searchResults))

	searchResults = searchResults[0:Min(len(searchResults), 10)]
	response := searchResponse{make([]*searchResponseData, len(searchResults))}
	for i, result := range searchResults {
		response.Data[i] = retreiveResponse(&index.Data, result)
	}

	data, _ := json.Marshal(response)
	fmt.Fprintf(w, string(data))

	elapsed := time.Since(start)
	logMessage += fmt.Sprintf(",%s", elapsed)

	fmt.Println(logMessage)
}

func handleSubs(w http.ResponseWriter, req *http.Request, index *Index) {
	id := req.FormValue("id")
	expandType, _ := strconv.ParseBool(req.FormValue("type"))

	fmt.Println("/subs id:", id, "type:", expandType)

	rID := parseRecordID(id)
	show := index.Data.Shows[rID.showID]
	file := &show.Files[rID.fileID]

	var record *Record
	response := &searchResponseData{show.Title, file.Name, make([]*Record, 0)}

	if expandType {
		record = GetRecord(file, rID.subID-1)
	} else {
		record = GetRecord(file, rID.subID+1)
	}

	if record != nil {
		response.Subs = append(response.Subs, record)
	}

	data, _ := json.Marshal(response)
	fmt.Fprintf(w, string(data))
}

func retreiveResponse(data *Data, result *recordID) *searchResponseData {
	if result == nil {
		return nil
	}

	show := data.Shows[result.showID]
	file := show.Files[result.fileID]
	records := retreiveRecordContext(&file, result.subID)

	return &searchResponseData{show.Title, file.Name, records}
}

func retreiveRecordContext(file *Showfile, id int) []*Record {
	records := make([]*Record, 0, 3)

	for _, i := range []int{-1, 0, 1} {
		sub := GetRecord(file, id+i)
		if sub != nil {
			records = append(records, sub)
		}
	}

	return records
}

func GetRecord(file *Showfile, id int) *Record {
	if id < 0 || id >= len(file.Records) {
		return nil
	}
	return &file.Records[id]
}
