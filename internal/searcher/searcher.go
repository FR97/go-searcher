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
	Preview  string
}

func Search(query SearchQuery, cache cache.Cache) []SearchResult {
	terms := parseTerms(query.Input)
	results := []SearchResult{}
	fileTermIndexes := map[string][]int{}

	for file, ftf := range cache.FileToTermFreq {
		score := 0.0
		for _, term := range terms {
			_tf, ok := ftf.TF[term]
			if !ok {
				score = 0
			} else {
				tf := float64(_tf.Count) / float64(ftf.TotalTermCount)
				fileTermIndexes[file] = append(fileTermIndexes[file], int(_tf.FirstIndex))
				idf := calculateIDF(cache, term)
				score += tf * idf
			}
		}

		if score > 0 {
			res := SearchResult{FilePath: file, Score: score}
			results = append(results, res)
		}

	}

	resultPage := sortAndPaginate(results, query)

	for i, res := range resultPage {
		resultPage[i].Preview = previewContent(res, terms, fileTermIndexes[res.FilePath], 20)
	}

	return resultPage
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
