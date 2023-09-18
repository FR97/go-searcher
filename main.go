package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fr97/go-searcher/parser"
)

type SearcherFlags struct {
	SearchPath string
}

func main() {

	flags := parseSearcherFlags()

	fmt.Println("search path: ", flags.SearchPath)

	err := parseFiles(flags.SearchPath)
	if err != nil {
		fmt.Println("error: ", err)
	}

}

func parseFiles(rootPath string) error {
	err := filepath.Walk(rootPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() {
				_, err := parser.ParseFile(path)
				if err != nil {
					fmt.Println(fmt.Errorf("error: %w", err))
				} else {
					fmt.Println("parsed: ", path)
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
