package lexer

import (
	"unicode"

	"github.com/surgebase/porter2"
)

type Lexer interface {
	NextToken() (string, bool)
}

type SimpleTermLexer struct {
	content  []rune
	position int
}

func NewLexer(content string, stemming bool) Lexer {
	if stemming {
		return newStemmingLexer(content)
	}
	return &SimpleTermLexer{content: []rune(content), position: 0}
}

func (l *SimpleTermLexer) NextToken() (string, bool) {
	l.incrementWhile(unicode.IsSpace)

	if l.position >= len(l.content) {
		return "", false
	}

	if unicode.IsLetter(l.content[l.position]) { // word token
		start := l.position
		l.incrementWhile(isAlpaNumeric)
		return string(l.content[start:l.position]), true
	} else if unicode.IsNumber(l.content[l.position]) { // number token
		start := l.position
		l.incrementWhile(unicode.IsNumber)
		return string(l.content[start:l.position]), true
	} else { // other tokens are treated as single chars
		l.position += 1
		return string(l.content[l.position-1 : l.position]), true
	}
}

func (l *SimpleTermLexer) incrementWhile(filter func(rune) bool) {
	for l.position < len(l.content) && filter(l.content[l.position]) {
		l.position += 1
	}
}

type StemmingLexer struct {
	simpleLexer SimpleTermLexer
}

func newStemmingLexer(content string) Lexer {
	return &StemmingLexer{simpleLexer: SimpleTermLexer{content: []rune(content), position: 0}}
}

func (l *StemmingLexer) NextToken() (string, bool) {
	token, ok := l.simpleLexer.NextToken()
	if !ok {
		return token, ok
	}

	stemmed := porter2.Stem(string(token))

	return stemmed, true
}

func isAlpaNumeric(r rune) bool {
	if !unicode.IsLetter(r) && !unicode.IsNumber(r) {
		return false
	}

	return true
}
