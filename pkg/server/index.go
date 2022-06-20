package main

import (
	"os"
	"path/filepath"
	"time"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/mapping"
	"github.com/blevesearch/bleve/search/query"

	"github.com/ishunyu/magpie-dict/pkg/analysis/singletoken"
	"github.com/ishunyu/magpie-dict/pkg/analysis/wholesentence"

	log "github.com/sirupsen/logrus"
)

func GetIndex(config *Config) *Index {
	data := GetData(config.DataPath)
	index := indexData(config.IndexPath, &data)
	return &Index{data, index}
}

type Index struct {
	Data   Data
	BIndex *bleve.Index
}

func (index *Index) Search(searchText string, showID string) []*recordID {
	queryString := "*" + searchText + "*"
	wildcardQuery := bleve.NewWildcardQuery(queryString)

	var newQuery query.Query
	if showID != "" {
		fieldQuery := bleve.NewQueryStringQuery("ShowID:" + showID)
		booleanQuery := bleve.NewBooleanQuery()
		booleanQuery.AddMust(fieldQuery, wildcardQuery)
		newQuery = booleanQuery
	} else {
		newQuery = wildcardQuery
	}

	bSearchRequest := bleve.NewSearchRequest(newQuery)
	bSearchResult, err := (*index.BIndex).Search(bSearchRequest)
	if err != nil {
		log.Error(err)
		return nil
	}

	numHits := len(bSearchResult.Hits)
	if numHits == 0 {
		return nil
	}

	searchResults := make([]*recordID, numHits)
	for i, match := range bSearchResult.Hits {
		searchResults[i] = parseRecordID(match.ID)
	}

	return searchResults
}

type message struct {
	ID     string
	ShowID string
	AText  string
	BText  string
}

func (msg message) Type() string {
	return "message"
}

type indexingDataVisitor struct {
	bIndex *bleve.Index
	added  *IndexManifest
}

func (visitor *indexingDataVisitor) start() {
	log.Info("===== Indexing BEGIN =====")
}

func (visitor *indexingDataVisitor) end(elapsed time.Duration) {
	log.Infof("===== Indexing COMPLETED (%v) =====", elapsed.Truncate(time.Millisecond))
}

func (visitor *indexingDataVisitor) startShow(show *Show) bool {
	if !visitor.added.HasShow(show.ID) {
		log.Infof("Skipped indexing \"%v\"", show.Title)
		return false
	}
	log.Infof("Begin indexing \"%v\"", show.Title)
	return true
}

func (visitor *indexingDataVisitor) endShow(show *Show, elapsed time.Duration) {
	log.Infof("Completed indexing \"%s\" (%v)", show.Title, elapsed.Truncate(time.Millisecond))
}

func (visitor *indexingDataVisitor) startFile(show *Show, file *Showfile) bool {
	if !visitor.added.Has(show.ID, file.Name) {
		log.Infof("Skipped indexing \"%v\" - %v...", show.Title, file.Name)
		return false
	}
	log.Infof("Begin indexing \"%v\" - %v...", show.Title, file.Name)
	return true
}

func (visitor *indexingDataVisitor) endFile(show *Show, file *Showfile, elapsed time.Duration) {
	log.Infof("Completed indexing \"%v\" - %s...DONE (%v)", show.Title, file.Name, elapsed.Truncate(time.Millisecond))
}

func (visitor *indexingDataVisitor) visitRecord(show *Show, file *Showfile, record *Record) {
	bMessage := message{record.ID, show.ID, record.A.Text, record.B.Text}
	(*visitor.bIndex).Index(bMessage.ID, bMessage)
}

func indexData(indexPath string, data *Data) *bleve.Index {
	indexManifestPath := filepath.Join(indexPath, "manifest.json")
	bleveIndexPath := filepath.Join(indexPath, "bleve")

	existingManifest, err := LoadIndexManifestFromFile(indexManifestPath)
	if err == nil {
		log.Info("Existing index manifest found.")
		log.Info("existing: ", existingManifest)
	} else {
		log.Info("No index manifest found.")
		os.RemoveAll(bleveIndexPath)
		existingManifest = NewIndexManifest()
	}

	bIndex, err := bleve.Open(bleveIndexPath)
	if err == nil {
		log.Info("Existing index found.")
	} else {
		log.Info("No index found.")

		mapping := getNewMapping()
		bIndex, err = bleve.New(bleveIndexPath, mapping)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	}

	newManifest := GetIndexManifest(data)
	log.Info("data: ", newManifest)

	added, removed := existingManifest.Compare(newManifest)
	log.Info("added: ", added)
	log.Info("removed: ", removed)

	data.Visit(&indexingDataVisitor{&bIndex, added})
	err = newManifest.SaveToFile(indexManifestPath)
	if err != nil {
		log.Warn("Failed to save/update index manifest file.")
	} else {
		log.Info("Index manifest file updated.")
	}

	return &bIndex
}

func getNewMapping() *mapping.IndexMappingImpl {
	indexMapping := bleve.NewIndexMapping()
	documentMapping := bleve.NewDocumentMapping()
	indexMapping.AddDocumentMapping("message", documentMapping)

	idFieldMapping := bleve.NewTextFieldMapping()
	idFieldMapping.Index = false
	documentMapping.AddFieldMappingsAt("ID", idFieldMapping)

	showIDFieldMapping := bleve.NewTextFieldMapping()
	showIDFieldMapping.Store = false
	showIDFieldMapping.Analyzer = singletoken.Analyzer
	documentMapping.AddFieldMappingsAt("ShowID", showIDFieldMapping)

	aTextFieldMapping := bleve.NewTextFieldMapping()
	aTextFieldMapping.Store = false
	aTextFieldMapping.Analyzer = wholesentence.Analyzer
	documentMapping.AddFieldMappingsAt("AText", aTextFieldMapping)

	bTextFieldMapping := bleve.NewTextFieldMapping()
	bTextFieldMapping.Store = false
	documentMapping.AddFieldMappingsAt("BText", bTextFieldMapping)

	return indexMapping
}
