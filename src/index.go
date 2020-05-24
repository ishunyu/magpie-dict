package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/mapping"
)

type Index struct {
	Records []Record
	BIndex  *bleve.Index
}

type Line struct {
	Start string `json:"start"`
	End   string `json:"end"`
	Text  string `json:"text"`
}

type Record struct {
	ID int `json:"id"`
	A  Line `json:"a"`
	B  Line `json:"b"`
}

type bRecord struct {
	ID    string
	AText string
	BText string
}

func GetIndex(config *Config) *Index {
	records := getRecords(config.GetDataPath())
	index := indexRecords(config.IndexPath, records)
	return &Index{records, index}
}

func (index *Index) Search(searchText string) []*Record {
	query := bleve.NewQueryStringQuery(searchText)
	searchRequest := bleve.NewSearchRequest(query)
	searchResult, err := (*index.BIndex).Search(searchRequest)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	hits := len(searchResult.Hits)
	if hits == 0 {
		fmt.Println("Zero hits")
		return nil
	}

	hitIds := make([]int, hits)
	for i, match := range searchResult.Hits {
		hitIds[i], err = strconv.Atoi((*match).ID)
		if err != nil {
			fmt.Println(err)
			return nil
		}
	}

	// fmt.Printf("Hits: %v\n", hitIds)
	return index.getResults(hitIds)
}

func (index *Index) GetRecord(id int) *Record {
	if id < 0 || id >= len(index.Records) {
		return nil
	}
	return &index.Records[id]
}

func (index *Index) getResults(hitIds []int) []*Record {
	if hitIds == nil {
		return nil
	}

	results := make([]*Record, len(hitIds))
	for i, id := range hitIds {
		results[i] = &(index.Records)[id]
	}
	return results
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

func indexRecords(indexPath string, records []Record) *bleve.Index {
	fmt.Print("Checking for indexes...")
	bIndex, err := bleve.Open(indexPath)
	if err == nil {
		fmt.Println("found!")
		return &bIndex
	}
	fmt.Println("not found :(")

	fmt.Println("Indexing...")
	mapping := getNewMapping()
	bIndex, err = bleve.New(indexPath, mapping)
	if err != nil {
		panic(err)
	}

	for _, r := range records {
		bMessage := bRecord{strconv.Itoa(r.ID), r.A.Text, r.B.Text}
		bIndex.Index(bMessage.ID, bMessage)
	}
	return &bIndex
}

func getNewMapping() *mapping.IndexMappingImpl {
	mapping := bleve.NewIndexMapping()
	return mapping
}
