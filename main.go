package main

import (
	_ "embed"
	"fmt"
	"time"

	"github.com/fr97/go-searcher/internal/cache"
	"github.com/fr97/go-searcher/internal/config"
	"github.com/fr97/go-searcher/internal/indexer"
	"github.com/fr97/go-searcher/internal/io"
	"github.com/fr97/go-searcher/internal/searcher"
	"github.com/fr97/go-searcher/internal/server"
)

//go:embed public/view/index.gohtml
var html string

func main() {
	cfg := config.LoadConfig()

	switch cfg.Command {
	case config.Index:
		timed(
			func() { indexer.Index(cfg) },
			func(d time.Duration) {
				fmt.Println("indexing took:", d.Milliseconds(), "ms")
			})
	case config.Search:
		indices := loadCacheFile(cfg)
		query := searcher.SearchQuery{
			Input:  cfg.SearchConfig.Input,
			Limit:  cfg.SearchConfig.Limit,
			Offset: cfg.SearchConfig.Offset,
		}
		timed(
			func() {
				sr := searcher.Search(query, indices)
				fmt.Println("results:", sr)
			},
			func(d time.Duration) {
				fmt.Println("search took:", d.Milliseconds(), "ms")
			})
	case config.Serve:
		indices := loadCacheFile(cfg)
		server.Serve(cfg, indices, html)
	default:
		config.ShowHelp()
	}
}

func loadCacheFile(cfg config.Config) cache.Cache {
	cache, err := io.ReadCache(cfg.CacheFilePath)
	if err != nil {
		panic("cache file not found, please run index command first")
	}
	return cache
}

func timed(f func(), cb func(time.Duration)) {
	start := time.Now()
	f()
	end := time.Now()
	cb(end.Sub(start))
}
