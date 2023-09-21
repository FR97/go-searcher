package config

import (
	"flag"
	"os"
	"strings"
)

type Command string

const (
	Index  Command = "index"
	Search Command = "search"
	Serve  Command = "serve"
	Help   Command = "help"
)

var Commands = map[string]Command{
	"index":  Index,
	"search": Search,
	"serve":  Serve,
	"help":   Help,
}

type Config struct {
	Command         Command
	IndexingPath    string
	IndicesFilePath string
	SearchQuery     SearchQuery
	ServerConfig    ServerConfig
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
	command := Help
	if len(os.Args) > 1 {
		c := os.Args[1]
		cmd, ok := Commands[c]
		if ok {
			command = cmd
		}
	}

	flagSet := flag.NewFlagSet("", flag.ExitOnError)

	config := &Config{Command: command}
	switch command {
	case Index:
		if len(os.Args) < 3 {
			panic("<index-path> is required")
		}
		config.IndexingPath = os.Args[2]
	case Search:
		if len(os.Args) < 3 {
			panic("<search-input> is required")
		}

		input := os.Args[2]
		if strings.TrimSpace(input) == "" {
			panic("<search-input> cannot be empty")
		}
		config.SearchQuery.Input = input

		flagSet.IntVar(&config.SearchQuery.Limit, "search-limit", 10, "search limit (default 10)")
		flagSet.IntVar(&config.SearchQuery.Offset, "search-offset", 0, "search offset (default 0)")
	case Serve:
		flagSet.IntVar(&config.ServerConfig.Port, "port", 8000, "server port (default 8000)")
	}

	flagSet.StringVar(&config.IndicesFilePath, "indices-file", "./indices.json", "path to indeces file for reading/saving indeces data")

	if len(os.Args) >= 3 {
		flagSet.Parse(os.Args[2:])
	}

	return *config
}
