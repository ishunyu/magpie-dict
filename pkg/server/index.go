package main

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/mapping"
	"github.com/blevesearch/bleve/search/query"

	"github.com/ishunyu/magpie-dict/pkg/analysis/singletoken"
	"github.com/ishunyu/magpie-dict/pkg/analysis/wholesentence"

	log "github.com/sirupsen/logrus"
)

func NewIndex(dataPath string, indexPath string) (*Index, error) {
	data, err := getData(dataPath)
	if err != nil {
		return nil, err
	}

	index, err := indexData(indexPath, data)
	if err != nil {
		return nil, err
	}

	return &Index{data, index}, nil
}

type Index struct {
	Data   *Data
	BIndex *bleve.Index
}

func (index *Index) Search(searchText string, showID string) []*recordID {
	queryString := "*" + strings.ToLower(searchText) + "*"
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
	bBatch *bleve.Batch
}

func (visitor *indexingDataVisitor) start() {
	log.Info("===== Indexing BEGIN =====")
}

func (visitor *indexingDataVisitor) end(start time.Time) {
	log.Infof("===== Indexing COMPLETED (%v) =====", time.Since(start).Truncate(time.Millisecond))
}

func (visitor *indexingDataVisitor) startShow(show *Show) bool {
	if !visitor.added.HasShow(show.ID) {
		log.Infof("Skipped indexing \"%v\"", show.Title)
		return false
	}
	log.Infof("Begin indexing \"%v\"", show.Title)
	return true
}

func (visitor *indexingDataVisitor) endShow(show *Show, start time.Time) {
	log.Infof("Completed indexing \"%s\" (%v)", show.Title, time.Since(start).Truncate(time.Millisecond))
}

func (visitor *indexingDataVisitor) startFile(show *Show, file *Showfile) bool {
	if !visitor.added.Has(show.ID, file.Name) {
		log.Infof("Skipped indexing \"%v\" - %v...", show.Title, file.Name)
		return false
	}
	log.Infof("Begin indexing \"%v\" - %v...", show.Title, file.Name)
	if visitor.bBatch != nil {
		log.Error("Batch should be null.")
		os.Exit(1)
	}

	visitor.bBatch = (*visitor.bIndex).NewBatch()
	return true
}

func (visitor *indexingDataVisitor) endFile(show *Show, file *Showfile, start time.Time) {
	(*visitor.bIndex).Batch(visitor.bBatch)
	visitor.bBatch = nil
	log.Infof("Completed indexing \"%v\" - %s...DONE (%v)", show.Title, file.Name, time.Since(start).Truncate(time.Millisecond))
}

func (visitor *indexingDataVisitor) visitRecord(show *Show, file *Showfile, record *Record) {
	bMessage := message{record.ID, show.ID, record.A.Text, record.B.Text}
	(*visitor.bBatch).Index(bMessage.ID, bMessage)
}

func indexData(indexPath string, data *Data) (*bleve.Index, error) {
	manifestPath := filepath.Join(indexPath, "manifest.json")
	blevePath := filepath.Join(indexPath, "bleve")

	manifestFromFile, manifestErr := LoadIndexManifestFromFile(manifestPath)
	bIndex, bleveErr := bleve.Open(blevePath)

	if manifestErr != nil || bleveErr != nil {
		log.Info("Index files not complete, recreating index.")
		if bleveErr == nil {
			bleveCloseErr := bIndex.Close()
			if bleveCloseErr != nil {
				return nil, bleveCloseErr
			}
		}

		os.RemoveAll(blevePath)
		os.Remove(blevePath)
		os.Remove(manifestPath)

		mapping := getNewMapping()
		bIndex, bleveErr = bleve.New(blevePath, mapping)
		if bleveErr != nil {
			log.Error(bleveErr)
			return nil, bleveErr
		}
		manifestFromFile = NewIndexManifest()
	}

	manifestFromData := GetIndexManifest(data)
	log.Info("manifest(data): ", manifestFromData)

	added, removed := manifestFromFile.Compare(manifestFromData)
	log.Info("manifest(added): ", added)
	log.Info("manifest(removed): ", removed)

	data.Visit(&indexingDataVisitor{&bIndex, added, nil})

	if len(added.Shows) != 0 || len(removed.Shows) != 0 {
		manifestSaveErr := manifestFromData.SaveToFile(manifestPath)
		if manifestSaveErr != nil {
			return nil, manifestSaveErr
		}
		log.Info("Index manifest file saved.")
	}

	return &bIndex, nil
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
