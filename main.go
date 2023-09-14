package main

import "os"

func main() {

	args := os.Args[1:]

	if len(args) == 0 {
		println("File path not provided")
		return
	}

	filePath := args[0]

	bytes, err := os.ReadFile(filePath)

	if err != nil {
		println("error reading a file: ", err.Error())
		return
	}

	str := string(bytes)
	println(str)
}
