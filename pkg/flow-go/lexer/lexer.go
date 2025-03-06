package lexer

import (
	"github.com/flowtemplates/cli/pkg/flow-go/token"
)

const eof = 0

type Lexer struct {
	input  string
	start  int
	pos    int
	tokens chan token.Token
}

func LexString(input string) *Lexer {
	l := &Lexer{
		input:  input,
		start:  0,
		pos:    0,
		tokens: make(chan token.Token, 2),
	}

	go l.run()
	return l
}

func (l *Lexer) emit(t token.Type) {
	if l.start < l.pos {
		l.tokens <- token.Token{
			Typ: t,
			Val: l.input[l.start:l.pos],
			Pos: l.start,
		}
		l.start = l.pos
	}
}

func (l *Lexer) NextToken() token.Token {
	return <-l.tokens
}

func (l *Lexer) run() {
	for state := lexText; state != nil; {
		state = state(l)
	}
	close(l.tokens)
}

func (l *Lexer) next() rune {
	if l.pos < len(l.input) {
		r := rune(l.input[l.pos])
		l.pos++
		return r
	}
	return eof
}

func (l *Lexer) back() {
	l.pos--
}

func (l *Lexer) peek() rune {
	if l.pos < len(l.input) {
		r := rune(l.input[l.pos])
		return r
	}
	return eof
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
