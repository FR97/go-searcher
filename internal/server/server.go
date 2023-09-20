package server

import (
	"fmt"
	"github.com/fr97/go-searcher/internal/config"
	"html/template"
	"net/http"
)

const HTML_TEMPLATES_PATH = "./public/view/"

func Serve(cfg config.Config) {

	http.HandleFunc("/", handler)

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
}

func handler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	fmt.Println("query:", query)
	if query == nil || len(query) <= 0 {
		fmt.Println("I m here1")
		htmlOK(w, response{HasQuery: false}, "index.gohtml")
	} else {
		fmt.Println("I m here2")
		htmlOK(w, response{HasQuery: true}, "index.gohtml")
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
