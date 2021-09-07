package singletoken

import (
	"github.com/blevesearch/bleve/analysis"
	"github.com/blevesearch/bleve/analysis/tokenizer/single"
	"github.com/blevesearch/bleve/registry"

	log "github.com/sirupsen/logrus"
)

const Analyzer = "singletoken_analyzer"

func AnalyzerConstructor(config map[string]interface{}, cache *registry.Cache) (*analysis.Analyzer, error) {
	defer log.Info("singletoken analyzer initialized.")
	tokenizer, err := cache.TokenizerNamed(single.Name)
	if err != nil {
		return nil, err
	}

	rv := analysis.Analyzer{
		Tokenizer: tokenizer,
	}
	return &rv, nil
}

func init() {
	registry.RegisterAnalyzer(Analyzer, AnalyzerConstructor)
}
