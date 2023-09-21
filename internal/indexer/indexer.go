package indexer

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fr97/go-searcher/internal/config"
	"github.com/fr97/go-searcher/internal/io"
	"github.com/fr97/go-searcher/internal/lexer"
)

type TermFrequency map[string]uint

func Index(cfg config.Config) {
	index, err := io.ReadIndexFile(cfg.IndicesFilePath)
	if err != nil {
		fmt.Println("Index file not found, creating new index")
	}
	io.ParseFiles(cfg,
		func(path string, fi os.FileInfo) bool {
			_, exists := index[path]

			if exists {
				fmt.Println("Skipping index indexed file:", path)
			}

			return !exists
		},
		func(file, content string) {
			st := time.Now()
			tf := IndexTermFreq(content)
			et := time.Now()

			fmt.Println("Indexing", file, "took", et.Sub(st).Milliseconds(), "ms")
			index[file] = tf
		},
		withError)

	io.WriteIndexFile(cfg.IndicesFilePath, index)
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
