package searcher

import (
	"github.com/fr97/go-searcher/internal/cache"
	"github.com/fr97/go-searcher/internal/lexer"
	"math"
	"sort"
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

	lexer := lexer.NewLexer(query.Input)
	terms := []string{}

	for {
		token, ok := lexer.NextToken()
		if !ok {
			break
		}

		terms = append(terms, string(token))
	}

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

	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	return results
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
