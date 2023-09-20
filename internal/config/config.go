package config

import (
	"flag"
)

type Mode string

const (
	Index  Mode = "index"
	Search Mode = "search"
	Serve  Mode = "serve"
)

type Config struct {
	IndexingPath  string
	IndexFilePath string
	SearchQuery   SearchQuery
	Mode          Mode
	Debug         bool
	ServerConfig  ServerConfig
}

type SearchQuery struct {
	Input  string
	Limit  int
	Offset int
}

type ServerConfig struct {
	Port int
}

func LoadConfig() Config {

	flags := &Config{}

	mode := flag.String("mode", "index", "mode to run (index|search|serve)")

	flag.StringVar(&flags.IndexingPath, "path", "./", "path tofile or directory to be indexed")
	flag.StringVar(&flags.IndexFilePath, "index-file-path", "./index.json", "path to index file for reading/saving index data")
	flag.BoolVar(&flags.Debug, "debug", false, "enable debug mode")

	flag.StringVar(&flags.SearchQuery.Input, "search-input", "", "search input")
	flag.IntVar(&flags.SearchQuery.Limit, "search-limit", 10, "search limit (default 10)")
	flag.IntVar(&flags.SearchQuery.Offset, "search-offset", 0, "search offset (default 0)")

	flag.IntVar(&flags.ServerConfig.Port, "server-port", 8000, "server port (default 8000)")

	flag.Parse()

	if mode == nil {
		panic("mode must be defined")
	}

	switch *mode {
	case "index":
		flags.Mode = Index
	case "search":
		flags.Mode = Search
	case "serve":
		flags.Mode = Serve
	}

	return *flags
}
