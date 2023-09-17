package parser

import (
	"encoding/xml"
	"errors"
	"os"
	"strings"
)

func ParseFile(filePath string) (string, error) {

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
	}

	return "", errors.New("unsupported extension: " + extension)
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
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	var v string
	if err := xml.Unmarshal(bytes, &v); err != nil {
		return "", err
	}

	return v, nil
}
