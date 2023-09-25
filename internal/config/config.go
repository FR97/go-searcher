package config

import (
	"flag"
	"fmt"
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
	Command       Command
	CacheFilePath string
	IndexConfig   IndexConfig
	SearchConfig  SearchQuery
	ServerConfig  ServerConfig
}

type IndexConfig struct {
	IndexingPath string
	Threads      int
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
	flagOffset := 2

	config := &Config{
		Command:      command,
		IndexConfig:  IndexConfig{},
		SearchConfig: SearchQuery{},
		ServerConfig: ServerConfig{}}
	switch command {
	case Index:
		if len(os.Args) < 3 {
			panic("<index-path> is required")
		}
		config.IndexConfig.IndexingPath = os.Args[2]
		flagOffset = 3
		flagSet.IntVar(&config.IndexConfig.Threads, "threads", 1, "number of threads (default 1 for non parallel execution)")
	case Search:
		if len(os.Args) < 3 {
			panic("<search-input> is required")
		}

		input := os.Args[2]
		if strings.TrimSpace(input) == "" {
			panic("<search-input> cannot be empty")
		}
		config.SearchConfig.Input = input
		flagOffset = 3
		flagSet.IntVar(&config.SearchConfig.Limit, "limit", 20, "search limit (default 20)")
		flagSet.IntVar(&config.SearchConfig.Offset, "offset", 0, "search offset (default 0)")
	case Serve:
		flagSet.IntVar(&config.ServerConfig.Port, "port", 8000, "server port (default 8000)")
	}

	flagSet.StringVar(&config.CacheFilePath, "cache-file", "./cache.json", "path to index file for reading/saving indeces data")

	if len(os.Args) > flagOffset {
		fmt.Println("parsing args")
		flagSet.Parse(os.Args[flagOffset:])
	}

	fmt.Println("config:", config)

	return *config
}
