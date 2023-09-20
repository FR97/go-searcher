package io

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func ReadIndexFile(path string) (map[string]map[string]uint, error) {
	indexed := map[string]map[string]uint{}
	bytes, err := os.ReadFile(path)
	if err != nil {
		return indexed, err
	}
	err = json.Unmarshal(bytes, &indexed)
	if err != nil {
		return indexed, err
	}
	return indexed, nil
}

func WriteIndexFile(path string, indexed map[string]map[string]uint) error {
	data, err := json.Marshal(indexed)
	if err != nil {
		panic("failed to marshal indexed data to json")
	}

	return os.WriteFile(path, data, 0644)
}

func GetFileNameForFilePath(path string) string {
	_, file := filepath.Split(path)
	return file
}
