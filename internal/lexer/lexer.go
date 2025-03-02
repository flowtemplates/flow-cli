package lexer

import (
	"github.com/templatesflow/cli/internal/token"
)

const eof = 0

type Lexer struct {
	input  string
	start  int
	pos    int
	tokens chan token.Token
}

func Lex(input string) *Lexer {
	l := &Lexer{
		input:  input,
		start:  0,
		pos:    0,
		tokens: make(chan token.Token, 2),
	}

	go l.run()
	return l
}

func (l *Lexer) emit(t token.Type) string {
	var s string
	if l.start < l.pos {
		s = l.input[l.start:l.pos]
		l.tokens <- token.Token{
			Typ: t,
			Val: s,
			Pos: l.start,
		}
		l.start = l.pos
	}

	return s
}

func (l *Lexer) NextToken() token.Token {
	return <-l.tokens
}

func (l *Lexer) run() {
	defer close(l.tokens)
	for state := lexText; state != nil; {
		state = state(l)
	}
}

func (l *Lexer) next() rune {
	if l.pos >= len(l.input) {
		return eof
	}
	r := rune(l.input[l.pos])
	l.pos++
	return r
}

func (l *Lexer) back() {
	l.pos--
}

func (l *Lexer) peek() rune {
	r := l.next()
	l.back()
	return r
}

func (l *Lexer) accept(valid string) bool {
	r := l.next()
	for _, c := range valid {
		if c == r {
			return true
		}
	}

	l.back()
	return false
}

func (l *Lexer) acceptRun(valid string) {
	for l.accept(valid) {
	}
}
