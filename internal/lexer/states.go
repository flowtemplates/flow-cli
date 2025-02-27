package lexer

import (
	"strings"
	"unicode"
)

type stateFn func(*Lexer) stateFn

func (l *Lexer) StartsWith(s string) bool {
	return strings.HasPrefix(l.input[l.pos:], s)
}

func lexText(l *Lexer) stateFn {
	for {
		if l.StartsWith(LeftExpr) {
			if l.start < l.pos {
				l.emit(TokenText)
			}
			return lexLeftExpr
		}
		if l.next() == 0 {
			if l.start < l.pos {
				l.emit(TokenText)
			}
			return nil
		}
	}
}

func lexLeftExpr(l *Lexer) stateFn {
	l.pos += len(LeftExpr)
	l.emit(TokenLeftExpr)
	return lexWhitespace(lexSymbol)
}

func lexSymbol(l *Lexer) stateFn {
	for {
		switch r := l.peek(); {
		case unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_':
			l.next()
		default:
			l.emit(TokenSymbol)
			return lexWhitespace(lexRightExpr)
		}
	}
}

func lexWhitespace(nextState stateFn) stateFn {
	return func(l *Lexer) stateFn {
		for unicode.IsSpace(l.peek()) {
			l.next()
		}
		l.emit(TokenWhitespace)
		return nextState
	}
}

func lexRightExpr(l *Lexer) stateFn {
	if l.StartsWith(RightExpr) {
		l.pos += len(RightExpr)
		l.emit(TokenRightExpr)
		return lexWhitespace(lexText)
	}
	return l.errorf("unexpected EOF")
}
