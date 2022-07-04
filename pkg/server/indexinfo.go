package main

import (
	"encoding/json"
	"time"
)

var exists = struct{}{}

type IndexShowManifest struct {
	Files map[string]struct{} `json:"files"`
}

// Records information about which files have already been indexed to reduce unnecessary re-index
type IndexManifest struct {
	Shows map[string]*IndexShowManifest `json:"shows"`
}

func newIndexShowManifest() *IndexShowManifest {
	return &IndexShowManifest{make(map[string]struct{})}
}

func NewIndexManifest() *IndexManifest {
	return &IndexManifest{make(map[string]*IndexShowManifest)}
}

func (manifest *IndexManifest) HasShow(showName string) bool {
	_, c := manifest.Shows[showName]
	return c
}

func (manifest *IndexManifest) Add(showName string, fileName string) {
	if !manifest.HasShow(showName) {
		manifest.Shows[showName] = newIndexShowManifest()
	}
	manifest.Shows[showName].Add(fileName)
}

func (manifest *IndexManifest) Has(showName string, fileName string) bool {
	show, c := manifest.Shows[showName]
	if !c {
		return false
	}
	return show.Has(fileName)
}

func (show *IndexShowManifest) Add(fileName string) {
	if !show.Has(fileName) {
		show.Files[fileName] = exists
	}
}

func (show *IndexShowManifest) Has(fileName string) bool {
	_, c := show.Files[fileName]
	return c
}

func LoadIndexManifestFromFile(filePath string) (*IndexManifest, error) {
	var manifest IndexManifest
	err := JsonLoadFromFile(filePath, &manifest)
	if err != nil {
		return nil, err
	}

	return &manifest, nil
}

func (manifest *IndexManifest) SaveToFile(filePath string) error {
	return JsonWriteToFile(filePath, manifest)
}

func (manifest *IndexManifest) Compare(to *IndexManifest) (added, removed *IndexManifest) {
	added = NewIndexManifest()
	removed = NewIndexManifest()

	for showName, show := range manifest.Shows {
		for fileName := range show.Files {
			if !to.Has(showName, fileName) {
				removed.Add(showName, fileName)
			}
		}
	}

	for showName, show := range to.Shows {
		for fileName := range show.Files {
			if !manifest.Has(showName, fileName) {
				added.Add(showName, fileName)
			}
		}
	}

	return
}

func (manifest *IndexManifest) String() string {
	byte, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return err.Error()
	}
	return string(byte)
}

func GetIndexManifest(data *Data) *IndexManifest {
	visitor := &indexManifestDataVisitor{NewIndexManifest()}
	data.Visit(visitor)
	return visitor.manifest
}

type indexManifestDataVisitor struct {
	manifest *IndexManifest
}

func (visitor *indexManifestDataVisitor) start()                                                 {}
func (visitor *indexManifestDataVisitor) end(start time.Time)                                    {}
func (visitor *indexManifestDataVisitor) endShow(show *Show, start time.Time)                    {}
func (visitor *indexManifestDataVisitor) visitRecord(show *Show, file *Showfile, record *Record) {}
func (visitor *indexManifestDataVisitor) endFile(show *Show, file *Showfile, start time.Time)    {}

func (visitor *indexManifestDataVisitor) startShow(show *Show) bool {
	return true
}

func (visitor *indexManifestDataVisitor) startFile(show *Show, file *Showfile) bool {
	visitor.manifest.Add(show.ID, file.Name)
	return false
}
