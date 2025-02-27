package lexer

import "fmt"

type Lexer struct {
	input  string
	start  int
	pos    int
	tokens chan Token
}

func Lex(input string) *Lexer {
	l := &Lexer{
		input:  input,
		start:  0,
		pos:    0,
		tokens: make(chan Token, 2),
	}

	go l.run()
	return l
}

func (l *Lexer) errorf(format string, args ...any) stateFn {
	l.tokens <- Token{
		Typ: TokenError,
		Val: fmt.Sprintf(format, args...),
	}
	l.pos = l.start
	return nil
}

func (l *Lexer) NextToken() Token {
	return <-l.tokens
}

func (l *Lexer) run() {
	for state := lexText; state != nil; {
		state = state(l)
	}
	close(l.tokens)
}

func (l *Lexer) emit(t TokenType) {
	if l.start != l.pos {
		l.tokens <- Token{
			Typ: t,
			Val: l.input[l.start:l.pos],
		}
		l.start = l.pos
	}
}

func (l *Lexer) next() rune {
	if l.pos >= len(l.input) {
		return 0
	}
	r := rune(l.input[l.pos])
	l.pos++
	return r
}

func (l *Lexer) peek() rune {
	if l.pos >= len(l.input) {
		return 0
	}
	return rune(l.input[l.pos])
}
