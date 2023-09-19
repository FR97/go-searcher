package indexer

import (
	"fmt"
	"os"
	"strings"
	"time"
	"unicode"

	"github.com/fr97/go-searcher/internal/config"
	"github.com/fr97/go-searcher/internal/io"
)

type TermFrequency map[string]uint

func Index(config config.Config, index map[string]TermFrequency) {
	io.ParseFiles(config,
		func(path string, fi os.FileInfo) bool {
			_, exists := index[path]
			return !exists
		},
		func(file, content string) {
			st := time.Now()
			tf := IndexTermFreq(content)
			et := time.Now()

			fmt.Println("Indexing", file, "took", et.Sub(st).Nanoseconds(), "ns")
			index[file] = tf
		},
		withError)
}

func withError(err error) {
	fmt.Println(fmt.Errorf("error:%w", err))
}

func IndexTermFreq(content string) TermFrequency {
	lexer := NewLexer(content)
	tf := TermFrequency{}

	for {
		token, ok := lexer.nextToken()

		if !ok {
			break
		}
		term := strings.ToLower(string(token))
		count := tf[term]
		tf[term] = count + 1
	}

	return tf
}

type Lexer struct {
	content  []rune
	position int
}

func NewLexer(content string) *Lexer {
	return &Lexer{content: []rune(content), position: 0}
}

func (l *Lexer) nextToken() ([]rune, bool) {
	l.incrementWhile(unicode.IsSpace)

	if l.position >= len(l.content) {
		return []rune{}, false
	}

	if unicode.IsLetter(l.content[l.position]) { // word token
		start := l.position
		l.incrementWhile(isAlpaNumeric)
		return l.content[start:l.position], true
	} else if unicode.IsNumber(l.content[l.position]) { // number token
		start := l.position
		l.incrementWhile(unicode.IsNumber)
		return l.content[start:l.position], true
	} else { // other tokens are treated as single chars
		l.position += 1
		return l.content[l.position-1 : l.position], true
	}
}

func (l *Lexer) incrementWhile(filter func(rune) bool) {
	for l.position < len(l.content) && filter(l.content[l.position]) {
		l.position += 1
	}
}

func isAlpaNumeric(r rune) bool {
	if !unicode.IsLetter(r) && !unicode.IsNumber(r) {
		return false
	}

	return true
}
