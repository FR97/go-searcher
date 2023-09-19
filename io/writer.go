package io

import "os"

func SaveToFile(path string, data []byte) error {
	return os.WriteFile(path, data, 0644)
}
