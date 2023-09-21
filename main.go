package main

import (
	_ "embed"
	"fmt"
	"github.com/fr97/go-searcher/internal/config"
	"github.com/fr97/go-searcher/internal/indexer"
	"github.com/fr97/go-searcher/internal/io"
	"github.com/fr97/go-searcher/internal/searcher"
	"github.com/fr97/go-searcher/internal/server"
	"runtime"
	"strings"
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

	osExt := ""
	if runtime.GOOS == "windows" {
		osExt = ".exe"
	}

	builder := strings.Builder{}

	builder.WriteString("Usage: go-searcher" + osExt + " <command> <required-arg-value> [--opt-arg opt-arg-value]\n")
	builder.WriteString("Commands:\n")
	builder.WriteString("  index <index-path> [--indices-file (default: ./indices.json)]\n")
	builder.WriteString("  search <search-input> [--indices-file (default: ./indices.json)]\n")
	builder.WriteString("  serve [--indices-file (default: ./indices.json)] [--port (default: 8080)]\n\n")
	builder.WriteString("Example for indexing all dirs/files under current dir and save in default indices file:\n")
	builder.WriteString("  go-searcher" + osExt + " index ./\n")
	builder.WriteString("\n")
	builder.WriteString("Example for indexing everything under ./custom dir and save in custom-indices.json file:\n")
	builder.WriteString("  go-searcher" + osExt + " index ./custom --indices-file ./custom-indices.json\n")
	builder.WriteString("\n")
	builder.WriteString("Example for searching 'hello world' in default indices file:\n")
	builder.WriteString("  go-searcher" + osExt + " search 'hello world'\n")
	builder.WriteString("\n")
	builder.WriteString("Example for starting server with default indices.json file on custom port 9999:\n")
	builder.WriteString("  go-searcher" + osExt + " serve --port 9999\n")

	fmt.Println(builder.String())
}

func getIndexFile(cfg config.Config) searcher.Index {
	indexed, err := io.ReadIndexFile(cfg.IndicesFilePath)
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
