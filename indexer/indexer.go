package indexer

import "fmt"

func Index(content string, index map[string]int) {
	index[fmt.Sprint(len(content))] = len(content)
}
