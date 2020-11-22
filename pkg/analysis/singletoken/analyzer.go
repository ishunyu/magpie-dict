package singletoken

import (
	"fmt"

	"github.com/blevesearch/bleve/analysis"
	"github.com/blevesearch/bleve/analysis/tokenizer/single"
	"github.com/blevesearch/bleve/registry"
)

const Analyzer = "singletoken_analyzer"

func AnalyzerConstructor(config map[string]interface{}, cache *registry.Cache) (*analysis.Analyzer, error) {
	defer fmt.Println("singletoken analyzer initialized.")
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
