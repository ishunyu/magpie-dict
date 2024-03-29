package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

type showsResponse struct {
	Shows []*showResponse `json:"shows"`
}

type showResponse struct {
	Name    string `json:"name"`
	Episode string `json:"episode"`
}

type searchResponse struct {
	Data []*searchResponseData `json:"data"`
}

type searchResponseData struct {
	Show    string    `json:"show"`
	Episode string    `json:"episode"`
	Subs    []*Record `json:"subs"`
}

type compareResponse struct {
	Status string `json:"status"`
	Output string `json:"output"`
}

func ShowsHandler(index *Index) HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request, stats *RequestStats) {
		handleShows(w, req, index, stats)
	}
}

func GetSearchHandler(index *Index) HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request, stats *RequestStats) {
		handleSearch(w, req, index, stats)
	}
}

func SubsHandler(index *Index) HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request, stats *RequestStats) {
		handleSubs(w, req, index, stats)
	}
}

func CompareHandler(tmpPath string, comparePath string, compareVenvPath string) HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request, stats *RequestStats) {
		handleCompare(w, req, tmpPath, comparePath, compareVenvPath, stats)
	}
}

func handleShows(w http.ResponseWriter, req *http.Request, index *Index, stats *RequestStats) {
	shows := make([]*showResponse, 0, len(index.Data.Shows))
	for _, show := range index.Data.Shows {
		file := ""
		for filename := range show.Files {
			if filename > file {
				file = filename
			}
		}

		showRes := showResponse{show.Title, file}
		shows = append(shows, &showRes)
	}
	showsRes := showsResponse{shows}

	data, _ := json.Marshal(showsRes)
	fmt.Fprintf(w, string(data))

	stats.Add("shows", shows)
}

func handleSearch(w http.ResponseWriter, req *http.Request, index *Index, stats *RequestStats) {
	searchText := req.FormValue("searchText")
	showID := req.FormValue("showID")

	stats.Add("searchText", searchText)
	stats.Add("showID", showID)

	searchResults := index.Search(searchText, showID)
	if searchResults == nil {
		searchResults = make([]*recordID, 0)
	}

	stats.Add("numResults", len(searchResults))

	searchResults = searchResults[0:Min(len(searchResults), 10)]
	response := searchResponse{make([]*searchResponseData, len(searchResults))}
	for i, result := range searchResults {
		response.Data[i] = retreiveResponse(index.Data, result)
	}

	data, _ := json.Marshal(response)
	fmt.Fprintf(w, string(data))
}

func handleSubs(w http.ResponseWriter, req *http.Request, index *Index, stats *RequestStats) {
	id := req.FormValue("id")
	expandType, _ := strconv.ParseBool(req.FormValue("type"))

	stats.Add("recordID", id)
	stats.Add("expandType", expandType)

	rID := parseRecordID(id)
	show := index.Data.Shows[rID.showID]
	file := show.Files[rID.filename]

	var record *Record
	response := &searchResponseData{show.Title, file.Name, make([]*Record, 0)}

	if expandType {
		record = GetRecord(&file, rID.subID-1)
	} else {
		record = GetRecord(&file, rID.subID+1)
	}

	if record != nil {
		response.Subs = append(response.Subs, record)
	}

	data, _ := json.Marshal(response)
	fmt.Fprintf(w, string(data))
}

func handleCompare(
	w http.ResponseWriter,
	req *http.Request,
	tmpPath string,
	comparePath string,
	compareVenvPath string,
	stats *RequestStats) {

	logger := log.WithField(KEY_REQUEST_ID, stats.RequestID())

	dir := getCompareTmpPath(tmpPath)
	defer os.RemoveAll(dir)

	stats.Add("dir", dir)

	chinese_file := saveFileFromRequest(req, "CHINESE_FILE", dir, "chinese_file.sbv", logger)
	original_file := saveFileFromRequest(req, "ORIGINAL_FILE", dir, "original_file.sbv", logger)
	revised_file := saveFileFromRequest(req, "REVISED_FILE", dir, "revised_file.sbv", logger)

	if original_file == "" || revised_file == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	output_file := filepath.Join(dir, "output.xlsx")
	var shCmd string
	if chinese_file == "" {
		shCmd = fmt.Sprintf("%s %s %s %s %s %s %s %s %s", ".", compareVenvPath, ";", "python", comparePath, "-o", output_file, original_file, revised_file)
	} else {
		shCmd = fmt.Sprintf("%s %s %s %s %s %s %s %s %s %s", ".", compareVenvPath, ";", "python", comparePath, "-o", output_file, original_file, revised_file, chinese_file)
	}
	logger.Infof("executing: \"%s\"", shCmd)
	cmd := exec.Command("/bin/sh", "-c", shCmd)

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		logger.Error(err)
		return
	}

	bytes, err := ioutil.ReadFile(output_file)
	if err != nil {
		logger.Error(err)
		return
	}

	w.Write(bytes)
}

func getCompareTmpPath(tmpPath string) string {
	compareTmpPath := filepath.Join(tmpPath, "compare")
	os.MkdirAll(compareTmpPath, 0700)
	ms := time.Now().UnixNano() / int64(time.Millisecond)
	dir, _ := ioutil.TempDir(compareTmpPath, strconv.FormatInt(ms, 10)+"_")
	return dir
}

func saveFileFromRequest(req *http.Request, key string, dir string, filename string, logger *log.Entry) string {
	formFile, _, err := req.FormFile(key)
	if err != nil {
		logger.Warn(err, " (", filename, ")")
		return ""
	}

	bytes, _ := ioutil.ReadAll(formFile)
	file := filepath.Join(dir, filename)
	err = ioutil.WriteFile(file, bytes, 0700)
	if err != nil {
		logger.Warn(err, " (", filename, ")")
		return ""
	}

	return file
}

func retreiveResponse(data *Data, result *recordID) *searchResponseData {
	if result == nil {
		return nil
	}

	show := data.Shows[result.showID]
	file := show.Files[result.filename]
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
