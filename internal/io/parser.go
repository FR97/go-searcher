package io

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"jaytaylor.com/html2text"
)

func ParseFiles(path string,
	fileFilter func(string, int64) bool,
	withContent func(string, int64, string),
	withError func(error)) error {
	err := filepath.Walk(path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			} else if info.IsDir() {
				return nil
			} else if !fileFilter(path, info.ModTime().UnixMilli()) {
				return nil
			}

			processFile(path,
				func(p, c string) {
					withContent(p, info.ModTime().UnixMilli(), c)
				},
				withError)

			return nil
		})

	return err
}

func ParseFilesMT(
	wg *sync.WaitGroup,
	path string,
	threads int,
	fileFilter func(string, int64) bool,
	withContent func(string, int64, string),
	withError func(error),
) error {

	ch := make(chan ParseReq, threads*2)
	defer close(ch)
	fmt.Println("starting", threads, "workers")
	for i := 0; i < threads; i++ {
		go ParseWorker(i, ch, func(pr ParseReq) {
			processFileMT(pr, withContent, withError)
		})
	}

	err := filepath.Walk(path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			} else if info.IsDir() {
				return nil
			} else if !fileFilter(path, info.ModTime().UnixMilli()) {
				return nil
			}

			req := ParseReq{path, info.ModTime().UnixMilli()}
			wg.Add(1)
			for loop := true; loop; {
				select {
				case ch <- req:
					loop = false
				default:
				}
			}

			return nil
		})

	wg.Wait()
	return err
}

func processFileMT(
	req ParseReq,
	withContent func(string, int64, string),
	withError func(error)) {
	content, err := parseFile(req.Path)
	if err != nil {
		withError(err)
	} else {
		withContent(req.Path, req.ModTime, content)
	}
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

	extension := filepath.Ext(filePath)

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

func readRawFileToString(filePath string) (string, error) {
	fb, err := os.ReadFile(filePath)

	if err != nil {
		return "", err
	}

	return string(fb), nil
}

func readXmlFileToString(filePath string) (string, error) {
	fb, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	decoder := xml.NewDecoder(bytes.NewBuffer(fb))

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
	fb, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	opts := html2text.Options{TextOnly: true}
	str, err := html2text.FromString(string(fb), opts)
	if err != nil {
		return "", err
	}

	return str, nil
}
