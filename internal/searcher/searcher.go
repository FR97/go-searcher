package searcher

import (
	"fmt"
	"github.com/fr97/go-searcher/internal/lexer"
	"sort"
)

type Index map[string]map[string]uint

type SearchQuery struct {
	Input  string
	Limit  int
	Offset int
}

type SearchResult struct {
	FilePath string
	Score    float64
}

func Search(query SearchQuery, index Index) []SearchResult {

	lexer := lexer.NewLexer(query.Input)
	terms := []string{}

	for {
		token, ok := lexer.NextToken()
		if !ok {
			break
		}

		fmt.Println("token:", token)
		terms = append(terms, string(token))
	}

	results := []SearchResult{}

	for file, tfMap := range index {
		totalFreq := 0.0
		for _, term := range terms {
			tf := findTermFreqIndex(tfMap, term)
			fmt.Println("tf:", tf)
			totalFreq += tf
		}

		if totalFreq > 0 {
			res := SearchResult{FilePath: file, Score: totalFreq}
			results = append(results, res)
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	return results
}

func findTermFreqIndex(tfMap map[string]uint, term string) float64 {
	if len(tfMap) <= 0 {
		return 0
	}

	totalCount := uint(0)
	for _, count := range tfMap {
		totalCount += count
	}

	tf := float64(tfMap[term])
	return tf / float64(totalCount)
}
