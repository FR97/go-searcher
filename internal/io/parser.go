package io

import (
	"bytes"
	"encoding/xml"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"jaytaylor.com/html2text"
)

func ParseFiles(
	path string,
	threads int,
	fileFilter func(string, int64) bool,
	withContent func(string, int64, string),
	withError func(error),
) error {
	ch := makeWorkerChannel(threads, withContent, withError)
	defer close(ch)

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
			for loop := true; loop; {
				select {
				case ch <- req:
					loop = false
				default:
				}
			}

			return nil
		})

	return err
}

func makeWorkerChannel(
	threads int,
	withContent func(string, int64, string),
	withError func(error)) chan ParseReq {
	ch := make(chan ParseReq, threads*2)

	for i := 0; i < threads; i++ {
		go ParseWorker(i, ch, func(pr ParseReq) {
			processFile(pr, withContent, withError)
		})
	}
	return ch
}

func processFile(
	req ParseReq,
	withContent func(string, int64, string),
	withError func(error)) {
	content, err := ParseFile(req.Path)
	if err != nil {
		withError(err)
	} else {
		withContent(req.Path, req.ModTime, content)
	}
}

func ParseFile(filePath string) (string, error) {

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
