package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Line struct {
	Start string `json:"start"`
	End   string `json:"end"`
	Text  string `json:"text"`
}

type Record struct {
	ID string `json:"id"`
	A  Line   `json:"a"`
	B  Line   `json:"b"`
}

type recordID struct {
	showID string
	fileID int
	subID  int
}

type Showfile struct {
	Name    string
	Records []Record
}

type Show struct {
	ID    string
	Title string
	Files []Showfile
}

type Data struct {
	Shows map[string]Show
}

type manifest struct {
	Title string `json:"title"`
}

type WalkFunc func(showID string, fileID int, record Record)

func GetData(dataPath string) Data {
	shows := make(map[string]Show)
	filepath.Walk(dataPath, func(showPath string, info os.FileInfo, err error) error {
		if showPath == dataPath || !info.IsDir() || filepath.Dir(showPath) != dataPath {
			return nil
		}

		show := getShow(showPath)
		shows[show.ID] = show
		return nil
	})
	return Data{shows}
}

func getShow(showPath string) Show {
	fmt.Print("Loading show data: " + showPath)

	manifestPath := filepath.Join(showPath, "manifest.json")
	manifestFile, err := os.Open(manifestPath)
	if err != nil {
		fmt.Println()
		fmt.Println(err)
		os.Exit(1)
	}
	defer manifestFile.Close()
	bytes, _ := ioutil.ReadAll(manifestFile)

	var manifestData manifest
	json.Unmarshal(bytes, &manifestData)
	id := filepath.Base(showPath)
	title := manifestData.Title

	fmt.Print(", id: " + id)
	fmt.Println(", title: " + title)

	return Show{id, title, getRecordFiles(filepath.Join(showPath, "data"), id)}
}

func getRecordFiles(filesPath string, showID string) []Showfile {
	files := make([]Showfile, 0, 10)
	filepath.Walk(filesPath, func(filePath string, info os.FileInfo, err error) error {
		if filePath == filesPath || strings.Split(filepath.Base(filePath), ".")[0] == "" {
			return nil
		}

		files = append(files, getRecordFile(filePath, showID, len(files)))
		return nil
	})
	return files
}

func getRecordFile(filePath string, showID string, fileID int) Showfile {
	filename := filepath.Base(filePath)
	name := strings.Split(filename, ".")[0]
	return Showfile{name, getRecords(filePath, showID, fileID)}
}

func getRecords(fileCSV string, showID string, fileID int) []Record {
	csvfile, _ := os.Open(fileCSV)
	defer csvfile.Close()
	data := csv.NewReader(csvfile)

	recordsData, _ := data.ReadAll()
	records := make([]Record, len(recordsData))
	for i, d := range recordsData {
		id := fmt.Sprintf("%s.%d.%d", showID, fileID, i)
		a := Line{d[0], d[1], d[2]}
		b := Line{d[3], d[4], d[5]}
		r := Record{id, a, b}

		records[i] = r
	}
	return records
}

func (data *Data) WalkRecords(f WalkFunc) {
	for showID, show := range data.Shows {
		fmt.Printf("Indexing %v\n", show.Title)
		start := time.Now()
		for fileID, file := range show.Files {
			fmt.Printf("Indexing episode %v... ", file.Name)
			startFile := time.Now()
			for _, record := range file.Records {
				f(showID, fileID, record)
			}
			elapsedFile := time.Since(startFile)
			fmt.Printf("(%v)\n", elapsedFile)
		}
		elapsed := time.Since(start)
		fmt.Printf("Finished %v (%v)\n", show.Title, elapsed)
	}
}

func parseRecordID(s string) *recordID {
	parts := strings.Split(s, ".")
	fileID, _ := strconv.Atoi(parts[1])
	subID, _ := strconv.Atoi(parts[2])

	return &recordID{
		showID: parts[0],
		fileID: fileID,
		subID:  subID,
	}
}
