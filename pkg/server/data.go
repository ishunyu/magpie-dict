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

	log "github.com/sirupsen/logrus"
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
	showID   string
	filename string
	subID    int
}

type Showfile struct {
	Name    string
	Records []Record
}

type Show struct {
	ID    string
	Title string
	Files map[string]Showfile
}

type Data struct {
	Shows map[string]Show
}

type manifest struct {
	Title string `json:"title"`
}

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
	log.Info("Finding show data. path: " + showPath)

	manifestPath := filepath.Join(showPath, "manifest.json")
	manifestFile, err := os.Open(manifestPath)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	defer manifestFile.Close()
	bytes, _ := ioutil.ReadAll(manifestFile)

	var manifestData manifest
	json.Unmarshal(bytes, &manifestData)
	id := filepath.Base(showPath)
	title := manifestData.Title

	log.Infof("Loading show data. id: %s, title: %s", id, title)

	return Show{id, title, getRecordFiles(filepath.Join(showPath, "data"), id)}
}

func getRecordFiles(filesPath string, showID string) map[string]Showfile {
	files := make(map[string]Showfile)
	filepath.Walk(filesPath, func(filePath string, info os.FileInfo, err error) error {
		if filePath == filesPath {
			return nil
		}

		fullfilename := filepath.Base(filePath)
		extIndex := strings.LastIndex(fullfilename, ".")
		filename := fullfilename[:extIndex]

		files[filename] = getRecordFile(filePath, showID, filename)
		return nil
	})
	return files
}

func getRecordFile(filePath string, showID string, filename string) Showfile {
	return Showfile{filename, getRecords(filePath, showID, filename)}
}

func getRecords(fileCSV string, showID string, filename string) []Record {
	csvfile, _ := os.Open(fileCSV)
	defer csvfile.Close()
	data := csv.NewReader(csvfile)

	recordsData, _ := data.ReadAll()
	records := make([]Record, len(recordsData))
	for i, d := range recordsData {
		id := fmt.Sprintf("%s.%s.%d", showID, filename, i)
		a := Line{d[0], d[1], d[2]}
		b := Line{d[3], d[4], d[5]}
		r := Record{id, a, b}

		records[i] = r
	}
	return records
}

func parseRecordID(s string) *recordID {
	filenameStartIndex := strings.Index(s, ".")
	filenameEndIndex := strings.LastIndex(s, ".")

	subID, _ := strconv.Atoi(s[filenameEndIndex+1:])

	return &recordID{
		showID:   s[:filenameStartIndex],
		filename: s[filenameStartIndex+1 : filenameEndIndex],
		subID:    subID,
	}
}

type DataVisitor interface {
	start()
	end(time.Time)
	startShow(*Show) bool
	endShow(*Show, time.Time)
	startFile(*Show, *Showfile) bool
	endFile(*Show, *Showfile, time.Time)
	visitRecord(*Show, *Showfile, *Record)
}

func (data *Data) Visit(visitor DataVisitor) {
	visitor.start()
	start := time.Now()

	for _, show := range data.Shows {
		show.visit(visitor)
	}

	visitor.end(start)
}

func (show *Show) visit(visitor DataVisitor) {
	start := time.Now()
	if visitor.startShow(show) {
		for _, file := range show.Files {
			file.visit(visitor, show)
		}
		visitor.endShow(show, start)
	}
}

func (file *Showfile) visit(visitor DataVisitor, show *Show) {
	start := time.Now()
	if visitor.startFile(show, file) {
		for _, record := range file.Records {
			visitor.visitRecord(show, file, &record)
		}
		visitor.endFile(show, file, start)
	}
}
