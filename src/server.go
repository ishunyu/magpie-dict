package main

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/mapping"
)

var fileCSV = "/Users/shun/code/magpie-dict/resource/data/EP21Outfile.csv"
var fileIndex = "/Users/shun/code/magpie-dict/tmp/index.bleve"

var indexData *bleve.Index = nil
var recordData *[][]string = nil

type bleveRecord struct {
	Index string
	AText string
	BText string
}

func searchHandler(w http.ResponseWriter, req *http.Request) {
	start := time.Now()

	searchText := req.FormValue("searchText")
	fmt.Print(searchText)

	searchResults := search(indexData, &searchText)
	if searchResults == nil {
		fmt.Println()
		return
	}

	results := findResults(searchResults)
	if results == nil {
		fmt.Println()
		return
	}

	for _, rec := range *results {
		t := strings.Join(rec[1:], ";")
		fmt.Fprint(w, t+"#")
	}
	elapsed := time.Since(start)
	fmt.Printf(" %s\n", elapsed)
}

func findResults(results *[]string) *[][]string {
	if results == nil {
		return nil
	}

	textResults := make([][]string, len(*results))
	for i, result := range *results {
		x, _ := strconv.Atoi(result)
		textResults[i] = (*recordData)[x]
	}
	return &textResults
}

func search(index *bleve.Index, s *string) *[]string {
	query := bleve.NewQueryStringQuery(*s)
	searchRequest := bleve.NewSearchRequest(query)
	searchResult, _ := (*index).Search(searchRequest)

	hits := len(searchResult.Hits)
	if hits == 0 {
		return nil
	}

	result := make([]string, hits)
	for i, match := range searchResult.Hits {
		result[i] = (*match).ID
	}
	return &result
}

func getRecordData() *[][]string {
	csvfile, _ := os.Open(fileCSV)
	data := csv.NewReader(csvfile)

	records, _ := data.ReadAll()
	indexedRecords := make([][]string, len(records))
	for i, r := range records {
		newRecord := make([]string, 1, len(r))
		newRecord[0] = strconv.Itoa(i)
		newRecord = append(newRecord, r...)

		indexedRecords[i] = newRecord
	}
	return &indexedRecords
}

func indexRecordData(file string, data *[][]string) *bleve.Index {
	fmt.Print("Checking for indexes...")
	index, err := bleve.Open(fileIndex)
	if err == nil {
		fmt.Println("found!")
		return &index
	}
	fmt.Println("not found :(")

	fmt.Println("Indexing...")
	mapping := getMapping()
	index, err = bleve.New(fileIndex, mapping)
	if err != nil {
		panic(err)
	}

	for _, r := range *data {
		message := bleveRecord{r[0], r[3], r[6]}
		fmt.Println(message)
		index.Index(message.Index, message)
	}
	return &index
}

func getMapping() *mapping.IndexMappingImpl {
	mapping := bleve.NewIndexMapping()
	return mapping
}

func main() {
	recordData = getRecordData()
	indexData = indexRecordData(fileIndex, recordData)

	http.Handle("/", http.FileServer(http.Dir("/Users/shun/code/magpie-dict/static")))
	http.HandleFunc("/search", searchHandler)

	port := 8090
	fmt.Printf("Starting server on port %v...\n", port)
	http.ListenAndServe(fmt.Sprint("localhost:", port), nil)
}
