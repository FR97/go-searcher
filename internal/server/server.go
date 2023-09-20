package server

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/fr97/go-searcher/internal/config"
	"github.com/fr97/go-searcher/internal/io"
	"github.com/fr97/go-searcher/internal/searcher"
)

const HTML_TEMPLATES_PATH = "./public/view/"

func Serve(cfg config.Config, indexed searcher.Index) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { handler(w, r, indexed) })

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

func handler(w http.ResponseWriter, r *http.Request, indexed searcher.Index) {
	query := r.URL.Query()
	fmt.Println("query:", query)
	if query == nil || len(query) <= 0 {
		htmlOK(w, response{HasQuery: false}, "index.gohtml")
	} else {
		input := query.Get("search-input")
		fmt.Println("Starting search for:", input)
		results := searcher.Search(searcher.SearchQuery{
			Input:  input,
			Limit:  10,
			Offset: 0,
		}, indexed)
		fmt.Println("Search results:")
		res := response{HasQuery: true, Results: []searchResult{}}
		for _, result := range results {
			res.Results = append(res.Results,
				searchResult{
					FileName: io.GetFileNameForFilePath(result.FilePath),
					FilePath: result.FilePath,
					Score:    result.Score,
				})
		}
		htmlOK(w, res, "index.gohtml")
	}
}

func htmlOK(w http.ResponseWriter, data interface{}, tmplFile string) {
	w.Header().Add("content-type", "html")

	tmpl, err := template.ParseFiles(HTML_TEMPLATES_PATH + tmplFile)
	if err != nil {
		fmt.Println("failed to parse template:", err)
	}

	tmpl.Execute(w, data)
}
