package jieba

import (
	"fmt"

	"github.com/blevesearch/bleve/analysis"
	"github.com/blevesearch/bleve/registry"
)

const Analyzer = "jieba_analyzer"

func AnalyzerConstructor(config map[string]interface{}, cache *registry.Cache) (*analysis.Analyzer, error) {
	defer fmt.Println("Jieba Analyzer initialized.")
	tokenizer, err := cache.TokenizerNamed(Tokenizer)
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
