package main

import (
	"fmt"
	"os"

	"github.com/fr97/go-searcher/parser"
)

func main() {

	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("File path not provided")
		return
	}

	filePath := args[0]

	str, err := parser.ParseFile(filePath)
	if err != nil {
		fmt.Println("error:", err)
	}

	fmt.Println(str)
}
