package searcher

import (
	"strings"

	"github.com/fr97/go-searcher/internal/io"
)

func previewContent(res SearchResult, terms []string, termIndexes []int, previewOffset int) string {
	fileContent, err := io.ParseFile(res.FilePath)
	if err != nil {
		println("Failed to read file for preview:", res.FilePath)
	}

	preview := ""
	for i := range terms {
		index := termIndexes[i]
		start := max(0, index-previewOffset)
		end := min(len(fileContent), index+previewOffset)
		preview += strings.ReplaceAll(fileContent[start:end], "\n", " ")
	}

	return preview
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
