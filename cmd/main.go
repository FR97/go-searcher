package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/fr97/go-searcher/internal/config"
	"github.com/fr97/go-searcher/internal/indexer"
	"github.com/fr97/go-searcher/internal/io"
	"github.com/fr97/go-searcher/internal/searcher"
	"github.com/fr97/go-searcher/internal/server"
)

type FileIndex map[string]map[string]uint

func main() {
	cfg := config.LoadConfig()

	fmt.Println("config:", cfg)

	bytes, err := os.ReadFile(cfg.IndexFilePath)
	indexed := FileIndex{}
	err = json.Unmarshal(bytes, &indexed)
	if err != nil {
		fmt.Println("indexed file not found, reindexing...")
	}

	switch cfg.Mode {
	case config.Index:
		timed(
			func() { indexer.Index(cfg, indexed) },
			func(d time.Duration) {
				fmt.Println("indexing took:", d.Milliseconds(), "ms")
			})
	case config.Search:
		query := searcher.SearchQuery{
			Input:  cfg.SearchQuery.Input,
			Limit:  cfg.SearchQuery.Limit,
			Offset: cfg.SearchQuery.Offset,
		}
		timed(
			func() {
				sr := searcher.Search(query, searcher.Index(indexed))
				fmt.Println("results:", sr)
			},
			func(d time.Duration) {
				fmt.Println("search took:", d.Milliseconds(), "ms")
			})
	case config.Serve:
		server.Serve(cfg)
	default:
		fmt.Println("unsupported mode:", cfg.Mode)
	}

	indexedJson, err := json.Marshal(indexed)
	if err != nil {
		panic("failed to marshal indexed data to json")
	}

	io.SaveToFile(cfg.IndexFilePath, indexedJson)
}

func timed(f func(), cb func(time.Duration)) {
	start := time.Now()
	f()
	end := time.Now()
	cb(end.Sub(start))
}
