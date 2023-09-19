package io

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/fr97/go-searcher/internal/config"
	"jaytaylor.com/html2text"
)

func ParseFiles(config config.Config,
	fileFilter func(string, os.FileInfo) bool,
	withContent func(string, string),
	withError func(error)) error {
	err := filepath.Walk(config.SearchPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			} else if info.IsDir() {
				return nil
			} else if !fileFilter(path, info) {
				fmt.Println("skipping indexed file:", info.Name())
				return nil
			}

			processFile(path, withContent, withError)

			return nil
		})

	return err
}

func processFile(path string, withContent func(string, string), withError func(error)) {
	content, err := parseFile(path)
	if err != nil {
		withError(err)
	} else {
		withContent(path, content)
	}
}

func parseFile(filePath string) (string, error) {

	if len(filePath) <= 0 {
		return "", errors.New("empty file path")
	}

	extension := exractExtensionWithDot(filePath)

	if len(extension) <= 0 {
		extension = ".txt"
	}

	switch extension {
	case ".txt", ".md":
		return readRawFileToString(filePath)
	case ".xml", ".xhtml":
		return readXmlFileToString(filePath)
	case ".html":
		return readHtmlFileToString(filePath)
	}

	return "", errors.New("unsupported extension " + extension)
}

func exractExtensionWithDot(filePath string) string {
	dotIndex := strings.LastIndex(filePath, ".")
	if dotIndex == -1 {
		return ""
	}
	return filePath[dotIndex:]
}

func readRawFileToString(filePath string) (string, error) {
	bytes, err := os.ReadFile(filePath)

	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func readXmlFileToString(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}

	defer f.Close()

	decoder := xml.NewDecoder(f)

	var sb strings.Builder

	for {
		tok, err := decoder.Token()
		if tok == nil || err == io.EOF {
			return sb.String(), nil
		} else if err != nil {
			return "", err
		}

		if cd, ok := tok.(xml.CharData); ok {
			str := string(cd)
			if len(strings.TrimSpace(str)) > 0 {
				sb.WriteString(str)
				sb.WriteString(" ")
			}
		}
	}
}

func readHtmlFileToString(filePath string) (string, error) {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	str, err := html2text.FromString(string(bytes), html2text.Options{TextOnly: true})
	if err != nil {
		return "", err
	}

	return str, nil
}
