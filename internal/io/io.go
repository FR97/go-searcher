package io

import (
	"encoding/json"
	"github.com/fr97/go-searcher/internal/cache"
	"os"
	"path/filepath"
)

func ReadCache(path string) (cache.Cache, error) {
	cache := cache.NewCache()
	bytes, err := os.ReadFile(path)
	if err != nil {
		return cache, err
	}
	err = json.Unmarshal(bytes, &cache)
	if err != nil {
		return cache, err
	}
	return cache, nil
}

func WriteCache(path string, cache cache.Cache) error {
	data, err := json.Marshal(cache)
	if err != nil {
		panic("failed to marshal indexed data to json")
	}

	return os.WriteFile(path, data, 0644)
}

func GetFileNameForFilePath(path string) string {
	_, file := filepath.Split(path)
	return file
}
