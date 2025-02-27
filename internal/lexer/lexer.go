package lexer

import (
	"fmt"
	"strings"
)

type stateFn func(*Lexer) stateFn

type Lexer struct {
	input  string
	start  int
	pos    int
	tokens []Token
}

func (l *Lexer) emit(t TokenType) {
	token := Token{
		typ: t,
		val: l.input[l.start:l.pos],
	}
	l.tokens = append(l.tokens, token)
	l.start = l.pos
}

func (l *Lexer) errorf(format string, args ...any) stateFn {
	l.tokens = append(l.tokens, Token{
		typ: TokenError,
		val: fmt.Sprintf(format, args...),
	})
	return nil
}

func (l *Lexer) next() rune {
	if l.pos >= len(l.input) {
		return 0
	}
	r := rune(l.input[l.pos])
	l.pos++
	return r
}

func (l *Lexer) backup() {
	if l.pos > 0 {
		l.pos--
	}
}

func (l *Lexer) isLeftExpr() bool {
	return strings.HasPrefix(l.input[l.pos:], "{{")
}

func (l *Lexer) isRightExpr() bool {
	return strings.HasPrefix(l.input[l.pos:], "}}")
}

func isWhitespace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n' || r == '\r'
}

func New(input string) *Lexer {
	return &Lexer{
		input:  input,
		start:  0,
		pos:    0,
		tokens: []Token{},
	}
}

func Lex(input string) []Token {
	l := New(input)
	for state := lexText; state != nil; {
		state = state(l)
	}
	return l.tokens
}

func lexText(l *Lexer) stateFn {
	for {
		if l.isLeftExpr() {
			if l.pos > l.start {
				l.emit(TokenText)
			}
			return lexLeftExpr
		}
		if l.next() == 0 {
			break
		}
	}
	if l.pos > l.start {
		l.emit(TokenText)
	}
	l.emit(TokenEOF)
	return nil
}

func lexLeftExpr(l *Lexer) stateFn {
	l.pos += 2
	l.emit(TokenLeftExpr)
	return lexInsideExpr
}

func lexRightExpr(l *Lexer) stateFn {
	l.pos += 2
	l.emit(TokenRightExpr)
	return lexText
}

func lexInsideExpr(l *Lexer) stateFn {
	for {
		if l.isRightExpr() {
			return lexRightExpr
		}

		r := l.next()
		switch {
		case r == 0:
			return l.errorf("unexpected end of file inside action")
		case isWhitespace(r):
			return lexWhitespace
		case isAlphaNumeric(r) || r == '_':
			l.backup()
			return lexSymbol
		default:
			return l.errorf("invalid character %q inside action", r)
		}
	}
}

func isAlphaNumeric(r rune) bool {
	return (r >= 'a' && r <= 'z') ||
		(r >= 'A' && r <= 'Z') ||
		(r >= '0' && r <= '9')
}

func lexSymbol(l *Lexer) stateFn {
	for {
		switch r := l.next(); {
		case r == 0:
			l.emit(TokenSymbol)
			return l.errorf("unexpected end of file inside action")
		case !isAlphaNumeric(r) && r != '_':
			l.backup()
			l.emit(TokenSymbol)
			return lexInsideExpr
		}
	}
}

func lexWhitespace(l *Lexer) stateFn {
	for {
		switch r := l.next(); {
		case !isWhitespace(r):
			l.backup()
			l.emit(TokenWhitespace)
			return lexInsideExpr
		case r == 0:
			l.emit(TokenWhitespace)
			return l.errorf("unexpected end of file inside action")
		}
	}
}
