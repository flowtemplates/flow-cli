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

func lexExpr(l *Lexer) stateFn {
	r := l.peek()
	if unicode.IsLetter(r) || r == '$' || r == '_' {
		return lexSymbolOrBoolean
	}

	if unicode.IsDigit(r) {
		return lexPositiveNum
	}

	if r == '-' {
		return lexNegativeNum
	}

	if r == '"' || r == '\'' {
		return lexString
	}

	return l.errorf("unexpected EOF")
}

func lexPositiveNum(l *Lexer) stateFn {
	for {
		switch r := l.peek(); {
		case unicode.IsDigit(r):
			l.next()
		case r == '.':
			return lexFloatNumber
		default:
			l.emit(TokenInteger)
			return lexWhitespace(lexRightExpr)
		}
	}
}

func lexNegativeNum(l *Lexer) stateFn {
	l.next()
	for {
		switch r := l.peek(); {
		case unicode.IsDigit(r):
			l.next()
		case r == '.':
			return lexFloatNumber
		default:
			l.emit(TokenInteger)
			return lexWhitespace(lexRightExpr)
		}
	}
}

func lexFloatNumber(l *Lexer) stateFn {
	l.next()
	for {
		switch r := l.peek(); {
		case unicode.IsDigit(r):
			l.next()
		default:
			l.emit(TokenFloat)
			return lexWhitespace(lexRightExpr)
		}
	}
}

func lexString(l *Lexer) stateFn {
	l.next()
	for {
		r := l.next()
		if r == '"' || r == '\'' {
			break
		}
	}
	l.emit(TokenString)
	return lexWhitespace(lexRightExpr)
}

func lexSymbolOrBoolean(l *Lexer) stateFn {
	for {
		switch r := l.peek(); {
		case unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || r == '$':
			l.next()
		default:
			word := l.input[l.start:l.pos]

			if word == "true" || word == "false" {
				l.emit(TokenBoolean)
				// check for other reserverd names
			} else {
				l.emit(TokenIdentifier)
			}

			return lexWhitespace(lexPipelineOrRightExpr)
		}
	}
}

func lexPipelineOrRightExpr(l *Lexer) stateFn {
	if l.StartsWith(Pipe) {
		l.pos += len(Pipe)
		l.emit(TokenPipe)
		return lexWhitespace(lexFilter)
	}

	return lexRightExpr(l)
}

func lexFilter(l *Lexer) stateFn {
	for {
		switch r := l.peek(); {
		case unicode.IsLetter(r):
			l.next()
		default:
			l.emit(TokenFilter)
			return lexWhitespace(lexPipelineOrRightExpr)
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

func lexLeftExpr(l *Lexer) stateFn {
	l.pos += len(LeftExpr)
	l.emit(TokenLeftExpr)
	return lexWhitespace(lexExpr)
}

func lexRightExpr(l *Lexer) stateFn {
	if l.StartsWith(RightExpr) {
		l.pos += len(RightExpr)
		l.emit(TokenRightExpr)
		return lexWhitespace(lexText)
	}

	if l.peek() == 0 {
		return l.errorf("unexpected EOF")
	}

	return l.errorf("unclosed expression")
}
