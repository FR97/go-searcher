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
	SearchPath    string
	IndexFilePath string
	Mode          Mode
	Debug         bool
	ServerConfig  ServerConfig
}

type ServerConfig struct {
	Port int
}

func LoadConfig() Config {

	flags := &Config{}

	mode := flag.String("mode", "index", "mode to run (index|search|serve)")

	flag.StringVar(&flags.SearchPath, "path", "./", "path tofile or directory to be indexed")
	flag.StringVar(&flags.IndexFilePath, "index-file-path", "./index.json", "path to index file for reading/saving index data")
	flag.BoolVar(&flags.Debug, "debug", false, "enable debug mode")

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
