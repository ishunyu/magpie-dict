package jieba

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/blevesearch/bleve/analysis"
	"github.com/blevesearch/bleve/registry"
	"github.com/yanyiwu/gojieba"
)

const Tokenizer = "jieba_tokenizer"

var ideographRegexp = regexp.MustCompile(`\p{Han}+`)

type JiebaTokenizer struct {
	tokenizer *gojieba.Jieba
}

func NewJiebaTokenizer() (analysis.Tokenizer, error) {
	return &JiebaTokenizer{
		tokenizer: gojieba.NewJieba(),
	}, nil
}

func (jt *JiebaTokenizer) Tokenize(input []byte) analysis.TokenStream {
	rv := make(analysis.TokenStream, 0)
	runeStart := 0
	start := 0
	end := 0
	pos := 1
	var width int
	words := jt.tokenizer.CutForSearch(string(input), true)
	for _, word := range words {
		end = start + len(word)
		token := analysis.Token{
			Term:     []byte(word),
			Start:    start,
			End:      end,
			Position: pos,
			Type:     detectTokenType(word),
		}
		// fmt.Println("jieba:", &token)
		rv = append(rv, &token)
		pos++
		runeStart += width
		start = end
	}
	return rv
}

func JiebaTokenizerConstructor(config map[string]interface{}, cache *registry.Cache) (analysis.Tokenizer, error) {
	defer fmt.Println("New Jieba Tokenizer initialized!")
	return NewJiebaTokenizer()
}

func detectTokenType(term string) analysis.TokenType {
	if ideographRegexp.MatchString(term) {
		if len(term) == 6 {
			// fmt.Printf("HI")
			return analysis.Double
		}
		return analysis.Ideographic
	}
	_, err := strconv.ParseFloat(term, 64)
	if err == nil {
		return analysis.Numeric
	}
	return analysis.AlphaNumeric
}

func init() {
	registry.RegisterTokenizer(Tokenizer, JiebaTokenizerConstructor)
}
