package server

import (
	"fmt"
	"net/http"

	"github.com/fr97/go-searcher/internal/config"
)

func Serve(cfg config.Config) {

	http.HandleFunc("/", handler)

	err := http.ListenAndServe(fmt.Sprint(":", cfg.ServerConfig.Port), nil)
	if err != nil {
		panic(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
  <div>
    Hello...
  <div>
  `
	htmlOK(w, html)
}

func htmlOK(w http.ResponseWriter, html string) {
	w.Header().Add("content-type", "html")
	fmt.Fprint(w, html)
}
