package indexer

import (
	"fmt"
	"strings"
	"sync"

	"github.com/fr97/go-searcher/internal/cache"
	"github.com/fr97/go-searcher/internal/config"
	"github.com/fr97/go-searcher/internal/io"
	"github.com/fr97/go-searcher/internal/lexer"
)

type indexRes struct {
	file string
	tf   cache.FileTermFrequency
}

func Index(cfg config.Config) {
	cached, err := io.ReadCache(cfg.CacheFilePath)
	if err != nil {
		fmt.Println("Index file not found, creating new index")
	}

	ch := make(chan indexRes, cfg.IndexConfig.Threads*4)
	doneParse := make(chan struct{})
	defer close(ch)
	defer close(doneParse)

	wg := &sync.WaitGroup{}
	mutex := &sync.Mutex{}
	go func() {
		io.ParseFiles(
			wg,
			cfg.IndexConfig.IndexingPath,
			cfg.IndexConfig.Threads,
			func(path string, modTime int64) bool {
				mutex.Lock()
				cached, exists := cached.FileToTermFreq[path]
				mutex.Unlock()

				if !exists {
					fmt.Println("Indexing new file:", path)
					return true
				}

				if modTime > cached.IndexTime {
					fmt.Println("Reindexing modified file:", path)
					return true
				}

				fmt.Println("Skipping already indexed file:", path)
				return false
			},
			func(file string, modTime int64, content string) {
				tf := IndexFileTermFreq(modTime, content)
				ch <- indexRes{file, tf}
			},
			func(err error) {
				withError(err)
				wg.Done()
			})
		fmt.Println("Parsing finished closing channel")
		doneParse <- struct{}{}
	}()

	go func() {
		for res := range ch {
			mutex.Lock()
			cached.FileToTermFreq[res.file] = res.tf
			for k := range res.tf.TF {
				cached.TermToFileFreq[k] = cached.TermToFileFreq[k] + 1
			}
			mutex.Unlock()
			fmt.Println("Indexed file:", res.file)
			wg.Done()
		}
	}()

	fmt.Println("Waiting for parse to finish...")
	<-doneParse

	io.WriteCache(cfg.CacheFilePath, cached)
}

func withError(err error) {
	fmt.Println(fmt.Errorf("error:%w", err))
}

func IndexFileTermFreq(modTime int64, content string) cache.FileTermFrequency {
	lexer := lexer.NewStemmingLexer(content)
	ftf := cache.FileTermFrequency{
		TF:        map[string]uint{},
		IndexTime: modTime,
	}

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
