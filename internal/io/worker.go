package io

import "fmt"

type ParseReq struct {
	Path    string
	ModTime int64
}

func ParseWorker(id int, ch chan ParseReq, f func(ParseReq)) {
	fmt.Println("Worker", id, "started")
	for req := range ch {
		fmt.Println("Worker", id, "received work")
		f(req)
		fmt.Println("Worker", id, "done work")
	}
	fmt.Println("Worker", id, "stopped")
}
