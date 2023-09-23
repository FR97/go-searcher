package indexer

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fr97/go-searcher/internal/cache"
	"github.com/fr97/go-searcher/internal/config"
	"github.com/fr97/go-searcher/internal/io"
	"github.com/fr97/go-searcher/internal/lexer"
)

func Index(cfg config.Config) {
	cache, err := io.ReadCache(cfg.IndicesFilePath)
	if err != nil {
		fmt.Println("Index file not found, creating new index")
	}

	io.ParseFiles(cfg,
		func(path string, fi os.FileInfo) bool {
			_, exists := cache.FileToTermFreq[path]

			if exists {
				fmt.Println("Skipping index indexed file:", path)
			}

			return !exists
		},
		func(file, content string) {
			st := time.Now()
			tf := IndexFileTermFreq(content)
			et := time.Now()

			fmt.Println("Indexing", file, "took", et.Sub(st).Milliseconds(), "ms")
			cache.FileToTermFreq[file] = tf

			for k := range tf.TF {
				cache.TermToFileFreq[k] = cache.TermToFileFreq[k] + 1
			}
		},
		withError)

	io.WriteCache(cfg.IndicesFilePath, cache)
}

func withError(err error) {
	fmt.Println(fmt.Errorf("error:%w", err))
}

func IndexFileTermFreq(content string) cache.FileTermFrequency {
	lexer := lexer.NewLexer(content)
	ftf := cache.FileTermFrequency{TF: map[string]uint{}}

	for {
		token, ok := lexer.NextToken()
		if !ok {
			break
		}

		term := strings.ToLower(string(token))
		count := ftf.TF[term]
		ftf.TF[term] = count + 1
		ftf.TotalTermCount++
	}

	return ftf
}
