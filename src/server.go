package main

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/mapping"
)

var fileA = ""
var fileCSV = "/Users/shun/code/magpie-dict/resource/data/EP21Outfile.csv"
var fileIndex = "/Users/shun/code/magpie-dict/tmp/index.bleve"

var indexAll *bleve.Index = nil
var recordsAll *[][]string = nil

func searchHandler(w http.ResponseWriter, req *http.Request) {
	searchText := req.FormValue("searchText")
	fmt.Printf("searching... %s\n", searchText)

	results := search(indexAll, &searchText)
	if results == nil {
		return
	}

	textResults := matchTextResults(results)
	if textResults == nil {
		return
	}

	for _, rec := range *textResults {
		t := strings.Join(rec[1:], ";")
		fmt.Fprint(w, t+"#")
	}
}

func matchTextResults(results *[]string) *[][]string {
	if results == nil {
		return nil
	}

	textResults := make([][]string, len(*results))
	for i, result := range *results {
		x, _ := strconv.Atoi(result)
		textResults[i] = (*recordsAll)[x]
	}
	return &textResults
}

func data() *[][]string {
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

type Record struct {
	Index string
	AText string
	BText string
}

func indexRecords(records *[][]string) *bleve.Index {
	index, err := bleve.Open(fileIndex)
	if err == nil {
		return &index
	}

	mapping := indexMapping()
	index, err = bleve.New(fileIndex, mapping)
	if err != nil {
		panic(err)
	}

	for _, r := range *records {
		message := Record{r[0], r[3], r[6]}
		fmt.Println(message)
		index.Index(message.Index, message)
	}
	return &index
}

func indexMapping() *mapping.IndexMappingImpl {
	mapping := bleve.NewIndexMapping()
	return mapping
}

func search(index *bleve.Index, s *string) *[]string {
	query := bleve.NewQueryStringQuery(*s)
	searchRequest := bleve.NewSearchRequest(query)
	searchResult, _ := (*index).Search(searchRequest)

	// fmt.Printf("searchResult: %s", searchResult)

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

func main() {
	recordsAll = data()
	indexAll = indexRecords(recordsAll)

	http.Handle("/", http.FileServer(http.Dir("/Users/shun/code/magpie-dict/static")))
	http.HandleFunc("/search", searchHandler)

	port := 8090
	fmt.Printf("Starting server on port %v...\n", port)
	http.ListenAndServe(fmt.Sprint("localhost:", port), nil)
}
