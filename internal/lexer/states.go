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
	defer l.emit(TokenText)
	for {
		if l.StartsWith(LeftExpr) {
			return lexLeftExpr
		}

		if l.StartsWith(LeftComm) {
			return lexLeftComm
		}

		if l.next() == 0 {
			return nil
		}
	}
}

func lexExpr(l *Lexer) stateFn {
	r := l.peek()
	if unicode.IsLetter(r) || r == '$' || r == '_' {
		return lexIdentifierOrBoolean
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

	if l.StartsWith(RightExpr) {
		return lexRightExpr
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
			return lexLineWhitespace(lexPipelineOrRightExpr)
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
			return lexLineWhitespace(lexPipelineOrRightExpr)
		}
	}
}

func lexFloatNumber(l *Lexer) stateFn {
	l.next() // skips dot
	for {
		r := l.peek()
		if !unicode.IsDigit(r) {
			l.emit(TokenFloat)
			return lexLineWhitespace(lexPipelineOrRightExpr)
		}

		l.next()
	}
}

func lexString(l *Lexer) stateFn {
	l.next() // skips leading "
	for {
		r := l.next()
		if r == '"' || r == '\'' {
			break
		}
	}

	l.emit(TokenString)
	return lexLineWhitespace(lexPipelineOrRightExpr)
}

func lexIdentifierOrBoolean(l *Lexer) stateFn {
	for {
		switch r := l.peek(); {
		case unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || r == '$':
			l.next()
		default:
			word := l.input[l.start:l.pos]

			if word == FalseLiteral || word == TrueLiteral {
				l.emit(TokenBoolean)
				// TODO: check for other reserverd names
			} else {
				l.emit(TokenIdentifier)
			}

			return lexLineWhitespace(lexPipelineOrRightExpr)
		}
	}
}

func lexPipelineOrRightExpr(l *Lexer) stateFn {
	if l.StartsWith(Pipe) {
		l.pos += len(Pipe)
		l.emit(TokenPipe)
		return lexLineWhitespace(lexFilterIdentifier)
	}

	if l.StartsWith(RightExpr) {
		return lexRightExpr
	}

	if l.peek() == 0 {
		return l.errorf("unexpected EOF")
	}

	return l.errorf("unclosed expression")
}

func lexFilterIdentifier(l *Lexer) stateFn {
	for {
		switch r := l.peek(); {
		case unicode.IsLetter(r):
			l.next()
		default:
			l.emit(TokenFilterIdentifier)
			return lexLineWhitespace(lexPipelineOrRightExpr)
		}
	}
}

func lexLineWhitespace(nextState stateFn) stateFn {
	return func(l *Lexer) stateFn {
		for {
			switch r := l.peek(); {
			case r == ' ' || r == '\t':
				l.next()
			case unicode.IsSpace(r):
				return l.errorf("line break")
			default:
				l.emit(TokenWhitespace)
				return nextState
			}
		}
	}
}

// func lexWhitespace(nextState stateFn) stateFn {
// 	return func(l *Lexer) stateFn {
// 		for unicode.IsSpace(l.peek()) {
// 			l.next()
// 		}
// 		l.emit(TokenWhitespace)
// 		return nextState
// 	}
// }

func lexLeftExpr(l *Lexer) stateFn {
	l.pos += len(LeftExpr)
	l.emit(TokenLeftExpr)
	return lexLineWhitespace(lexExpr)
}

func lexRightExpr(l *Lexer) stateFn {
	l.pos += len(RightExpr)
	l.emit(TokenRightExpr)
	return lexText
}

func lexLeftComm(l *Lexer) stateFn {
	l.pos += len(LeftComm)
	l.emit(TokenLeftComm)
	return lexComm
}

func lexRightComm(l *Lexer) stateFn {
	l.pos += len(RightComm)
	l.emit(TokenRightComm)
	return lexText
}

func lexComm(l *Lexer) stateFn {
	for {
		if l.peek() == 0 {
			return l.errorf("unexpected EOF")
		}

		if l.StartsWith(RightComm) {
			l.emit(TokenCommText)
			return lexRightComm
		}

		l.next()
	}
}
