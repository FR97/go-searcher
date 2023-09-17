package parser

import (
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

	if extension == ".txt" {
		return readRawFileToString(filePath)
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
