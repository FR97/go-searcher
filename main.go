package main

import (
	_ "embed"
	"fmt"
	"github.com/fr97/go-searcher/internal/config"
	"github.com/fr97/go-searcher/internal/indexer"
	"github.com/fr97/go-searcher/internal/io"
	"github.com/fr97/go-searcher/internal/searcher"
	"github.com/fr97/go-searcher/internal/server"
	"time"
)

type FileIndex map[string]map[string]uint

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
		indexed := getIndexFile(cfg)
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
		indexed := getIndexFile(cfg)
		server.Serve(cfg, indexed, html)
	default:
		help()
	}
}

func help() {
	fmt.Println("Usage: go-searcher [command]")
	fmt.Println("Commands:")
	fmt.Println("  index <index-path>")
	fmt.Println("  search <search-input>")
	fmt.Println("  serve")
}

func getIndexFile(cfg config.Config) searcher.Index {
	indexed, err := io.ReadIndexFile(cfg.IndexFilePath)
	if err != nil {
		panic("invalid index file, please run indexer first")
	}
	return searcher.Index(indexed)
}

func timed(f func(), cb func(time.Duration)) {
	start := time.Now()
	f()
	end := time.Now()
	cb(end.Sub(start))
}
