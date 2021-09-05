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

func ShowsHandler(index *Index) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		handleShows(w, req, index)
	}
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

func CompareHandler(tmpPath string, comparePath string, compareVenvPath string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		handleCompare(w, req, tmpPath, comparePath, compareVenvPath)
	}
}

func handleShows(w http.ResponseWriter, req *http.Request, index *Index) {
	start := time.Now()
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

	elapsed := time.Since(start)
	fmt.Println("/shows " + fmt.Sprintf(",%s", elapsed))
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

func handleCompare(w http.ResponseWriter, req *http.Request, tmpDir string, comparePath string, compareVenvPath string) {
	compareDir := filepath.Join(tmpDir, "compare")
	os.MkdirAll(compareDir, 0700)
	ms := time.Now().UnixNano() / int64(time.Millisecond)
	dir, _ := ioutil.TempDir(compareDir, strconv.FormatInt(ms, 10)+"_")
	fmt.Println(dir)
	defer os.RemoveAll(dir)

	chinese_file := getAndSaveFile(req, "CHINESE_FILE", dir, "chinese_file.sbv")
	original_file := getAndSaveFile(req, "ORIGINAL_FILE", dir, "original_file.sbv")
	revised_file := getAndSaveFile(req, "REVISED_FILE", dir, "revised_file.sbv")

	fmt.Println("chinese_file: ", chinese_file)
	fmt.Println("original_file: ", original_file)
	fmt.Println("revised_file", revised_file)

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
	fmt.Println(shCmd)
	cmd := exec.Command("/bin/sh", "-c", shCmd)

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(output_file)

	bytes, err := ioutil.ReadFile(output_file)
	if err != nil {
		fmt.Println(err)
		return
	}

	w.Write(bytes)
}

func getAndSaveFile(req *http.Request, key string, dir string, filename string) string {
	formFile, _, err := req.FormFile(key)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	bytes, _ := ioutil.ReadAll(formFile)
	file := filepath.Join(dir, filename)
	err = ioutil.WriteFile(file, bytes, 0700)
	if err != nil {
		fmt.Println(err)
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
