package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Line struct {
	Start string `json:"start"`
	End   string `json:"end"`
	Text  string `json:"text"`
}

type Record struct {
	ID int  `json:"id"`
	A  Line `json:"a"`
	B  Line `json:"b"`
}

type Showfile struct {
	Name    string
	Records []Record
}

type Show struct {
	Title string
	Files []Showfile
}

type Data struct {
	Shows []Show
}

type manifest struct {
	Title string `json:"title"`
}

type WalkFunc func(showId int, fileId int, record Record)

func GetData(dataPath string) Data {
	shows := make([]Show, 0, 1)
	filepath.Walk(dataPath, func(showPath string, info os.FileInfo, err error) error {
		if showPath == dataPath || !info.IsDir() || filepath.Dir(showPath) != dataPath {
			return nil
		}

		shows = append(shows, getShow(showPath))
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
	title := manifestData.Title
	fmt.Println(", title: " + title)

	return Show{title, getRecordFiles(filepath.Join(showPath, "data"))}
}

func getRecordFiles(filesPath string) []Showfile {
	files := make([]Showfile, 0, 10)
	filepath.Walk(filesPath, func(filePath string, info os.FileInfo, err error) error {
		if filePath == filesPath || strings.Split(filepath.Base(filePath), ".")[0] == "" {
			return nil
		}

		files = append(files, getRecordFile(filePath))
		return nil
	})
	return files
}

func getRecordFile(filePath string) Showfile {
	filename := filepath.Base(filePath)
	name := strings.Split(filename, ".")[0]
	return Showfile{name, getRecords(filePath)}
}

func getRecords(fileCSV string) []Record {
	csvfile, _ := os.Open(fileCSV)
	data := csv.NewReader(csvfile)

	recordsData, _ := data.ReadAll()
	records := make([]Record, len(recordsData))
	for i, d := range recordsData {
		a := Line{d[0], d[1], d[2]}
		b := Line{d[3], d[4], d[5]}
		r := Record{i, a, b}

		records[i] = r
	}
	return records
}

func (data *Data) WalkRecords(f WalkFunc) {
	for showId, show := range data.Shows {
		fmt.Printf("Indexing %v\n", show.Title)
		start := time.Now()
		for fileId, file := range show.Files {
			fmt.Printf("Indexing episode %v... ", file.Name)
			startFile := time.Now()
			for _, record := range file.Records {
				f(showId, fileId, record)
			}
			elapsedFile := time.Since(startFile)
			fmt.Printf("(%v)\n", elapsedFile)
		}
		elapsed := time.Since(start)
		fmt.Printf("Finished %v (%v)\n", show.Title, elapsed)
	}
}
