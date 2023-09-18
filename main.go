package main

import (
	"flag"
	"fmt"
	"github.com/fr97/go-searcher/indexer"
	"github.com/fr97/go-searcher/parser"
	"os"
	"path/filepath"
)

type SearcherFlags struct {
	SearchPath string
}

func main() {
	flags := parseSearcherFlags()

	fmt.Println("search path: ", flags.SearchPath)

	indexedFiles := make(map[string]map[string]int)

	err := parseFiles(flags.SearchPath, func(file, content string) {
		_, exists := indexedFiles[file]
		if !exists {
			indexed := map[string]int{}
			indexer.Index(content, indexed)
			indexedFiles[file] = indexed
		}
	})
	if err != nil {
		fmt.Println("error: ", err)
	}

	fmt.Println("indexed: ", indexedFiles)
}

func parseFiles(rootPath string, withContent func(string, string)) error {
	err := filepath.Walk(rootPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() {
				content, err := parser.ParseFile(path)
				if err != nil {
					fmt.Println(fmt.Errorf("error: %w", err))
				} else {
					withContent(path, content)
				}
			}

			return nil
		})

	return err
}

func parseSearcherFlags() SearcherFlags {
	pPath := flag.String("path", "./", "file or directory to be indexed")

	flag.Parse()

	path := ""
	if pPath != nil {
		path = *pPath
	}

	return SearcherFlags{SearchPath: path}
}
