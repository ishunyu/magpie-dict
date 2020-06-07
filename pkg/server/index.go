package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/mapping"
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

func GetIndex(config *Config) *Index {
	data := GetData(config.GetDataPath())
	index := indexData(config.IndexPath, &data)
	return &Index{data, index}
}

func (index *Index) Search(searchText string) [][]int {
	query := bleve.NewQueryStringQuery(searchText)
	searchRequest := bleve.NewSearchRequest(query)
	searchResult, err := (*index.BIndex).Search(searchRequest)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	hits := len(searchResult.Hits)
	if hits == 0 {
		return nil
	}

	hitsIds := make([][]int, hits)
	for i, match := range searchResult.Hits {
		parts := strings.Split(match.ID, ".")
		nums := make([]int, len(parts))
		for j, s := range parts {
			nums[j], _ = strconv.Atoi(s)
		}
		hitsIds[i] = nums
	}

	return hitsIds
}

func indexData(indexPath string, data *Data) *bleve.Index {
	bIndex, err := bleve.Open(indexPath)
	if err == nil {
		fmt.Println("Index found!")
		return &bIndex
	}
	fmt.Println("Index found not found!")

	fmt.Println("Indexing started")
	start := time.Now()
	mapping := getNewMapping()
	bIndex, err = bleve.New(indexPath, mapping)
	if err != nil {
		panic(err)
	}

	data.WalkRecords(func(showId int, fileId int, record Record) {
		id := fmt.Sprintf("%d.%d.%d", showId, fileId, record.ID)
		bMessage := message{id, record.A.Text, record.B.Text}
		bIndex.Index(bMessage.ID, bMessage)
	})

	elapsed := time.Since(start)
	fmt.Printf("Indexing completed! (%v)\n", elapsed)
	return &bIndex
}

func getNewMapping() *mapping.IndexMappingImpl {
	mapping := bleve.NewIndexMapping()
	return mapping
}
