package searcher

import (
	"fmt"
	"github.com/fr97/go-searcher/internal/lexer"
	"math"
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
		score := 0.0
		for _, term := range terms {
			tf := findTermFreq(tfMap, term)
			idf := findInverseDocFreq(index, term)
			score += tf * idf
		}

		if score > 0 {
			res := SearchResult{FilePath: file, Score: score}
			results = append(results, res)
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	return results
}

func findTermFreq(tfMap map[string]uint, term string) float64 {
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

func findInverseDocFreq(index Index, term string) float64 {
	docCount := float64(len(index))
	termOccurrence := float64(0)
	for _, tfMap := range index {
		if _, ok := tfMap[term]; ok {
			termOccurrence++
		}
	}

	if termOccurrence == 0 {
		termOccurrence = 1
	}

	idf := math.Log10(docCount / termOccurrence)
	return idf
}
