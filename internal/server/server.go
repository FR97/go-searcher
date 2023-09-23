package server

import (
	"fmt"
	"github.com/fr97/go-searcher/internal/cache"
	"github.com/fr97/go-searcher/internal/config"
	"github.com/fr97/go-searcher/internal/io"
	"github.com/fr97/go-searcher/internal/searcher"
	"html/template"
	"net/http"
)

func Serve(cfg config.Config, cache cache.Cache, html string) {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, cache, html)
	})

	fmt.Println("Starting server on port:", cfg.ServerConfig.Port)
	err := http.ListenAndServe(fmt.Sprint(":", cfg.ServerConfig.Port), nil)
	if err != nil {
		panic(err)
	}
}

type response struct {
	HasQuery bool
	Results  []searchResult
}

type searchResult struct {
	FileName string
	FilePath string
	Score    float64
}

func handler(w http.ResponseWriter, r *http.Request, cache cache.Cache, html string) {
	query := r.URL.Query()
	if query == nil || len(query) <= 0 {
		htmlOK(w, response{HasQuery: false}, html)
	} else {
		input := query.Get("search-input")
		results := searcher.Search(searcher.SearchQuery{
			Input:  input,
			Limit:  10,
			Offset: 0,
		}, cache)
		res := response{HasQuery: true, Results: []searchResult{}}
		for _, result := range results {
			res.Results = append(res.Results,
				searchResult{
					FileName: io.GetFileNameForFilePath(result.FilePath),
					FilePath: result.FilePath,
					Score:    result.Score,
				})
		}
		htmlOK(w, res, html)
	}
}

func htmlOK(w http.ResponseWriter, res response, tmplFile string) {
	w.Header().Add("content-type", "html")

	tmpl, err := template.New("index.gohtml").Parse(tmplFile)
	if err != nil {
		fmt.Println("failed to parse template:", err)
	}

	tmpl.Execute(w, res)
}
