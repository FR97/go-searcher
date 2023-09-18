package indexer

import (
	"strings"
	"unicode"
)

func IndexTermFreq(content string) map[string]uint {
	lexer := NewLexer(content)
	tf := map[string]uint{}

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
