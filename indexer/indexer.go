package indexer

import (
	"fmt"
	"strings"
	"unicode"
)

func Index(content string, index map[string]uint) {

	lexer := NewLexer(content)

	for {
		token, ok := lexer.nextToken()

		if !ok {
			break
		}
		fmt.Println("token:", string(token), "position:", lexer.position)
		strToken := strings.ToLower(string(token))
		count := index[strToken]
		index[strToken] = count + 1
	}
}

type Lexer struct {
	content  []rune
	position int
}

func NewLexer(content string) *Lexer {
	return &Lexer{content: []rune(content), position: 0}
}

func (l *Lexer) nextToken() ([]rune, bool) {
	l.lTrim()

	if l.position >= len(l.content) {
		return []rune{}, false
	}

	if unicode.IsLetter(l.content[l.position]) {
		i := l.position
		for i < len(l.content) && isAlpaNumeric(l.content[i]) {
			i += 1
		}
		start := l.position
		l.position = i
		return l.content[start:i], true
	} else if unicode.IsNumber(l.content[l.position]) {
		i := l.position
		for i < len(l.content) && unicode.IsNumber(l.content[i]) {
			i += 1
		}
		start := l.position
		l.position = i
		return l.content[start:i], true
	}

	return []rune{}, true
}

func (l *Lexer) lTrim() {
	for l.position < len(l.content) && unicode.IsSpace(l.content[l.position]) {
		l.position += 1
	}
}

func isAlpaNumeric(r rune) bool {
	if !unicode.IsLetter(r) && !unicode.IsNumber(r) {
		return false
	}

	return true
}
