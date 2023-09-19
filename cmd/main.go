package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/fr97/go-searcher/internal/indexer"
	"github.com/fr97/go-searcher/internal/io"
	"os"
	"path/filepath"
)

type SearcherFlags struct {
	SearchPath    string
	IndexFilePath string
}

type FileIndex map[string]indexer.TermFrequency

func main() {
	flags := parseSearcherFlags()

	fmt.Println("search path:", flags.SearchPath)
	fmt.Println("index file path:", flags.IndexFilePath)

	bytes, err := os.ReadFile(flags.IndexFilePath)
	indexed := FileIndex{}
	err = json.Unmarshal(bytes, &indexed)
	if err != nil {
		fmt.Println("indexed file not found, reindexing...")
	}

	err = parseFiles(flags.SearchPath, func(file, content string) {
		_, exists := indexed[file]
		if !exists {
			tf := indexer.IndexTermFreq(content)

			fmt.Println("term frequency:", tf)

			indexed[file] = tf
		}
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("indexed:", indexed)

	indexedJson, err := json.Marshal(indexed)
	if err != nil {
		panic("failed to marshal indexed data to json")
	}

	io.SaveToFile(flags.IndexFilePath, indexedJson)
}

func parseFiles(rootPath string, withContent func(string, string)) error {
	err := filepath.Walk(rootPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			} else if info.IsDir() {
				return nil
			}

			processFile(path, withContent)

			return nil
		})

	return err
}

func processFile(path string, withContent func(string, string)) {
	content, err := io.ParseFile(path)
	if err != nil {
		fmt.Println(fmt.Errorf("error: %w", err))
	} else {
		withContent(path, content)
	}
}

func parseSearcherFlags() SearcherFlags {
	pSearchPath := flag.String("path", "./", "path tofile or directory to be indexed")
	pIndexPath := flag.String("index-file-path", "./index.json", "path to index file for reading/saving index data")
	flag.Parse()

	path := ""
	if pSearchPath != nil {
		path = *pSearchPath
	}

	indexPath := ""
	if pIndexPath != nil {
		indexPath = *pIndexPath
	}

	return SearcherFlags{SearchPath: path, IndexFilePath: indexPath}
}
