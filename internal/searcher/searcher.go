package searcher

import (
	"math"
	"sort"

	"github.com/fr97/go-searcher/internal/cache"
	"github.com/fr97/go-searcher/internal/lexer"
)

type SearchQuery struct {
	Input  string
	Limit  int
	Offset int
}

type SearchResult struct {
	FilePath string
	Score    float64
}

func Search(query SearchQuery, cache cache.Cache) []SearchResult {
	terms := parseTerms(query.Input)
	results := []SearchResult{}

	for file, ftf := range cache.FileToTermFreq {
		score := 0.0
		for _, term := range terms {
			tf := caclulateTF(ftf, term)
			idf := calculateIDF(cache, term)
			score += tf * idf
		}

		if score > 0 {
			res := SearchResult{FilePath: file, Score: score}
			results = append(results, res)
		}
	}

	return sortAndPaginate(results, query)
}

func parseTerms(input string) []string {
	lexer := lexer.NewStemmingLexer(input)
	terms := []string{}
	for {
		token, ok := lexer.NextToken()
		if !ok {
			break
		}

		terms = append(terms, token)
	}
	return terms
}

func caclulateTF(ftf cache.FileTermFrequency, term string) float64 {
	if len(ftf.TF) <= 0 {
		return 0
	}

	tf := float64(ftf.TF[term])
	return tf / float64(ftf.TotalTermCount)
}

func calculateIDF(cache cache.Cache, term string) float64 {
	docCount := float64(len(cache.FileToTermFreq))
	termOccurrence := float64(cache.TermToFileFreq[term])

	if termOccurrence == 0 {
		termOccurrence = 1
	}

	idf := math.Log10(docCount / termOccurrence)

	return idf
}

func sortAndPaginate(results []SearchResult, query SearchQuery) []SearchResult {
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	start := query.Limit * query.Offset
	end := query.Limit * (query.Offset + 1)

	if start > len(results) {
		return []SearchResult{}
	}

	if end > len(results) {
		return results[start:]
	}

	return results[start:end]
}
