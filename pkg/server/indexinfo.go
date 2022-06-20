package main

import (
	"encoding/json"
	"time"
)

var exists = struct{}{}

type indexShowInfo struct {
	Files map[string]struct{} `json:"files"`
}

// Records information about which files have already been indexed to reduce unnecessary re-index
type IndexInfo struct {
	Shows map[string]*indexShowInfo `json:"shows"`
}

func newIndexShowInfo() *indexShowInfo {
	return &indexShowInfo{make(map[string]struct{})}
}

func NewIndexInfo() *IndexInfo {
	return &IndexInfo{make(map[string]*indexShowInfo)}
}

func (info *IndexInfo) HasShow(showName string) bool {
	_, c := info.Shows[showName]
	return c
}

func (info *IndexInfo) Add(showName string, fileName string) {
	if !info.HasShow(showName) {
		info.Shows[showName] = newIndexShowInfo()
	}
	info.Shows[showName].Add(fileName)
}

func (info *IndexInfo) Has(showName string, fileName string) bool {
	show, c := info.Shows[showName]
	if !c {
		return false
	}
	return show.Has(fileName)
}

func (show *indexShowInfo) Add(fileName string) {
	if !show.Has(fileName) {
		show.Files[fileName] = exists
	}
}

func (show *indexShowInfo) Has(fileName string) bool {
	_, c := show.Files[fileName]
	return c
}

func LoadFromFile(filePath string) (*IndexInfo, error) {
	var info IndexInfo
	err := JsonLoadFromFile(filePath, &info)
	if err != nil {
		return nil, err
	}

	return &info, nil
}

func (info *IndexInfo) SaveToFile(filePath string) error {
	return JsonWriteToFile(filePath, info)
}

func (from *IndexInfo) Compare(to *IndexInfo) (added, removed *IndexInfo) {
	added = NewIndexInfo()
	removed = NewIndexInfo()

	for showName, show := range from.Shows {
		for fileName := range show.Files {
			if !to.Has(showName, fileName) {
				removed.Add(showName, fileName)
			}
		}
	}

	for showName, show := range to.Shows {
		for fileName := range show.Files {
			if !from.Has(showName, fileName) {
				added.Add(showName, fileName)
			}
		}
	}

	return
}

func (info *IndexInfo) String() string {
	byte, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		return err.Error()
	}
	return string(byte)
}

func TransformToIndexInfo(data *Data) *IndexInfo {
	visitor := &indexInfoDataVisitor{NewIndexInfo()}
	data.Visit(visitor)
	return visitor.indexInfo
}

type indexInfoDataVisitor struct {
	indexInfo *IndexInfo
}

func (visitor *indexInfoDataVisitor) start()                                                    {}
func (visitor *indexInfoDataVisitor) end(elapsed time.Duration)                                 {}
func (visitor *indexInfoDataVisitor) endShow(show *Show, elapsed time.Duration)                 {}
func (visitor *indexInfoDataVisitor) visitRecord(show *Show, file *Showfile, record *Record)    {}
func (visitor *indexInfoDataVisitor) endFile(show *Show, file *Showfile, elapsed time.Duration) {}

func (visitor *indexInfoDataVisitor) startShow(show *Show) bool {
	return true
}

func (visitor *indexInfoDataVisitor) startFile(show *Show, file *Showfile) bool {
	visitor.indexInfo.Add(show.ID, file.Name)
	return false
}
