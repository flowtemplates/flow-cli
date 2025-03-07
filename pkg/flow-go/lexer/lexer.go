package lexer

import (
	"github.com/flowtemplates/cli/pkg/flow-go/token"
)

const eof = 0

type Lexer struct {
	input  string
	start  token.Position
	pos    token.Position
	tokens chan token.Token
}

func LexString(input string) *Lexer {
	l := &Lexer{
		input: input,
		start: token.Position{
			Offset: 0,
			Line:   1,
			Column: 1,
		},
		tokens: make(chan token.Token, 2),
	}
	l.pos = l.start

	go l.run()
	return l
}

func LexStringTokens(input string) []token.Token {
	l := LexString(input)
	var tokens []token.Token
	for {
		tok := l.NextToken()
		if tok.Typ == token.EOF {
			break
		}
		tokens = append(tokens, tok)
	}

	return tokens
}

func (l *Lexer) emit(t token.Type) {
	if l.start.Offset < l.pos.Offset {
		l.tokens <- token.Token{
			Typ: t,
			Val: l.input[l.start.Offset:l.pos.Offset],
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
	if l.pos.Offset >= len(l.input) {
		return eof
	}

	r := rune(l.input[l.pos.Offset])
	l.pos.Offset++

	// Handle newline properly
	if r == '\n' {
		l.pos.Line++
		l.pos.Column = 1 // Reset column on new line
	} else {
		l.pos.Column++
	}

	return r
}

func (l *Lexer) back() {
	if l.pos.Offset <= 0 {
		return
	}

	l.pos.Offset--
	if l.input[l.pos.Offset] == '\n' {
		l.pos.Line-- // Move back a line
		// Find previous line start to restore column correctly
		l.pos.Column = 1
		for i := l.pos.Offset - 1; i >= 0 && l.input[i] != '\n'; i-- {
			l.pos.Column++
		}
	} else {
		l.pos.Column--
	}
}

func (l *Lexer) peek() rune {
	if l.pos.Offset < len(l.input) {
		r := rune(l.input[l.pos.Offset])
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
