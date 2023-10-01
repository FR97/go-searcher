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
	file    string
	tf      cache.FileTermFrequency
	indexed bool
}

func Index(cfg config.Config) {
	cached, err := io.ReadCache(cfg.CacheFilePath)
	if err != nil {
		fmt.Println("Index file not found, creating new index")
	}

	newCache := cache.NewCache()

	ch := make(chan indexRes, cfg.IndexConfig.Threads*16)
	defer close(ch)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		io.ParseFiles(
			cfg.IndexConfig.IndexingPath,
			cfg.IndexConfig.Threads,
			func(path string, modTime int64) bool {
				cached, exists := cached.FileToTermFreq[path]
				wg.Add(1)

				if !exists {
					fmt.Println("Indexing new file:", path)
					return true
				}

				if modTime > cached.IndexTime {
					fmt.Println("Reindexing modified file:", path)
					return true
				}

				fmt.Println("Skipping already indexed file:", path)
				ch <- indexRes{path, cached, true}
				return false
			},
			func(file string, modTime int64, content string) {
				tf := IndexFileTermFreq(modTime, content)
				ch <- indexRes{file, tf, false}
			},
			func(err error) {
				withError(err)
				wg.Done()
			})
		wg.Done()
	}()

	go func() { collectIndexResults(ch, newCache, wg) }()

	wg.Wait()

	io.WriteCache(cfg.CacheFilePath, newCache)
}

func collectIndexResults(ch chan indexRes, cache cache.Cache, wg *sync.WaitGroup) {
	for res := range ch {
		cache.FileToTermFreq[res.file] = res.tf
		if !res.indexed {
			for k := range res.tf.TF {
				cache.TermToFileFreq[k] = cache.TermToFileFreq[k] + 1
			}
		}
		wg.Done()
	}
}

func withError(err error) {
	fmt.Println(fmt.Errorf("error:%w", err))
}

func IndexFileTermFreq(modTime int64, content string) cache.FileTermFrequency {
	lexer := lexer.NewLexer(content, true)
	ftf := cache.FileTermFrequency{
		TF:        map[string]uint{},
		IndexTime: modTime,
	}

	for {
		token, ok := lexer.NextToken()
		if !ok {
			break
		}

		term := strings.ToLower(token)
		count := ftf.TF[term]
		ftf.TF[term] = count + 1
		ftf.TotalTermCount++
	}

	return ftf
}
