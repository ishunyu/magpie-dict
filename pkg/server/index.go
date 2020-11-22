package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/mapping"

	"github.com/ishunyu/magpie-dict/pkg/analysis/wholesentence"
)

type Index struct {
	Data   Data
	BIndex *bleve.Index
}

type message struct {
	ID    string
	AText string
	BText string
}

type searchResult struct {
	showID string
	fileID int
	subID  int
}

func (msg message) Type() string {
	return "message"
}

func GetIndex(config *Config) *Index {
	data := GetData(config.DataPath)
	index := indexData(config.IndexPath, &data)
	return &Index{data, index}
}

func (index *Index) Search(searchText string) []searchResult {
	queryString := "*" + searchText + "*"
	query := bleve.NewWildcardQuery(queryString)

	bSearchRequest := bleve.NewSearchRequest(query)
	bSearchResult, err := (*index.BIndex).Search(bSearchRequest)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	numHits := len(bSearchResult.Hits)
	if numHits == 0 {
		return nil
	}

	searchResults := make([]searchResult, numHits)
	for i, match := range bSearchResult.Hits {
		parts := strings.Split(match.ID, ".")
		fileID, _ := strconv.Atoi(parts[1])
		subID, _ := strconv.Atoi(parts[2])
		searchResults[i] = searchResult{
			showID: parts[0],
			fileID: fileID,
			subID:  subID,
		}
	}

	return searchResults
}

func indexData(indexPath string, data *Data) *bleve.Index {
	bIndex, err := bleve.Open(indexPath)
	if err == nil {
		fmt.Println("Index found.")
		return &bIndex
	}
	fmt.Println("Index not found.")

	mapping := getNewMapping()
	bIndex, err = bleve.New(indexPath, mapping)
	if err != nil {
		panic(err)
	}

	fmt.Println("Indexing started.")
	start := time.Now()
	data.WalkRecords(func(showID string, fileID int, record Record) {
		bMessage := message{record.ID, record.A.Text, record.B.Text}
		bIndex.Index(bMessage.ID, bMessage)
	})

	elapsed := time.Since(start)
	fmt.Printf("Indexing completed. (%v)\n", elapsed)
	return &bIndex
}

func getNewMapping() *mapping.IndexMappingImpl {
	indexMapping := bleve.NewIndexMapping()
	documentMapping := bleve.NewDocumentMapping()
	indexMapping.AddDocumentMapping("message", documentMapping)

	idFieldMapping := bleve.NewTextFieldMapping()
	idFieldMapping.Index = false
	documentMapping.AddFieldMappingsAt("ID", idFieldMapping)

	aTextFieldMapping := bleve.NewTextFieldMapping()
	aTextFieldMapping.Analyzer = wholesentence.Analyzer
	documentMapping.AddFieldMappingsAt("AText", aTextFieldMapping)

	return indexMapping
}
