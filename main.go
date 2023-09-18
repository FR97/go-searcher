package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fr97/go-searcher/parser"
)

func main() {

	pPath := flag.String("path", "./", "file or directory to be indexed")

	flag.Parse()

	path := ""
	if pPath != nil {
		path = *pPath
	}

	fmt.Println("path: ", path)

	err := parseFiles(path)
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

func isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return fileInfo.IsDir(), err
}
