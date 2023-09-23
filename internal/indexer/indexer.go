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
	cache, err := io.ReadCache(cfg.CacheFilePath)
	if err != nil {
		fmt.Println("Index file not found, creating new index")
	}

	io.ParseFiles(cfg,
		func(path string, fi os.FileInfo) bool {
			cached, exists := cache.FileToTermFreq[path]

			if !exists {
				fmt.Println("Indexing new file:", path)
				return true
			}

			if fi.ModTime().UnixMilli() > cached.IndexTime {
				fmt.Println("Reindexing modified file:", path)
				return true
			}

			fmt.Println("Skipping already indexed file:", path)
			return false
		},
		func(file, content string) {
			st := time.Now()
			tf := IndexFileTermFreq(content)
			et := time.Now()

			cache.FileToTermFreq[file] = tf

			for k := range tf.TF {
				cache.TermToFileFreq[k] = cache.TermToFileFreq[k] + 1
			}

			fmt.Println("Indexing", file, "took", et.Sub(st).Milliseconds(), "ms")
		},
		withError)

	io.WriteCache(cfg.CacheFilePath, cache)
}

func withError(err error) {
	fmt.Println(fmt.Errorf("error:%w", err))
}

func IndexFileTermFreq(content string) cache.FileTermFrequency {
	lexer := lexer.NewLexer(content)
	ftf := cache.FileTermFrequency{TF: map[string]uint{}, IndexTime: time.Now().UnixMilli()}

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
