package indexer

import (
	"fmt"
	"github.com/fr97/go-searcher/internal/config"
	"github.com/fr97/go-searcher/internal/io"
	"github.com/fr97/go-searcher/internal/lexer"
	"os"
	"strings"
	"time"
)

type TermFrequency map[string]uint

func Index(cfg config.Config, index map[string]map[string]uint) {
	io.ParseFiles(cfg,
		func(path string, fi os.FileInfo) bool {
			_, exists := index[path]
			return !exists
		},
		func(file, content string) {
			st := time.Now()
			tf := IndexTermFreq(content)
			et := time.Now()

			fmt.Println("Indexing", file, "took", et.Sub(st).Nanoseconds(), "ns")
			index[file] = tf
		},
		withError)
}

func withError(err error) {
	fmt.Println(fmt.Errorf("error:%w", err))
}

func IndexTermFreq(content string) TermFrequency {
	lexer := lexer.NewLexer(content)
	tf := TermFrequency{}

	for {
		token, ok := lexer.NextToken()

		if !ok {
			break
		}
		term := strings.ToLower(string(token))
		count := tf[term]
		tf[term] = count + 1
	}

	return tf
}
