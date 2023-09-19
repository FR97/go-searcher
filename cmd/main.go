package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/fr97/go-searcher/internal/config"
	"github.com/fr97/go-searcher/internal/indexer"
	"github.com/fr97/go-searcher/internal/io"
	"github.com/fr97/go-searcher/internal/server"
)

type FileIndex map[string]indexer.TermFrequency

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
		indexer.Index(cfg, indexed)
	case config.Serve:
		server.Serve(cfg)
	default:
		fmt.Println("unsupported mode:", cfg.Mode)
	}

	fmt.Println("indexed:", indexed)

	indexedJson, err := json.Marshal(indexed)
	if err != nil {
		panic("failed to marshal indexed data to json")
	}

	io.SaveToFile(cfg.IndexFilePath, indexedJson)
}
