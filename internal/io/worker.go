package io

type ParseReq struct {
	Path    string
	ModTime int64
}

func ParseWorker(id int, ch chan ParseReq, f func(ParseReq)) {
	for req := range ch {
		f(req)
	}
}
