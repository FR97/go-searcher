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

type FileIndex map[string]indexer.TermFrequency

func main() {
	flags := parseSearcherFlags()

	fmt.Println("search path:", flags.SearchPath)

	indexed := make(map[string]map[string]uint)

	err := parseFiles(flags.SearchPath, func(file, content string) {
		_, exists := indexed[file]
		if !exists {
			tf := indexer.IndexTermFreq(content)

			fmt.Println("term frequency:", tf)

			indexed[file] = tf
		}
	})
	if err != nil {
		fmt.Println("error:", err)
	}

	fmt.Println("indexed:", indexed)
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
