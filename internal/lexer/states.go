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

		if l.StartsWith(LeftStmt) {
			return lexLeftStmt
		}

		if l.StartsWith(LeftComm) {
			return lexLeftComm
		}

		if l.StartsWith(RightExpr) {
			return lexRightExpr
		}

		if l.StartsWith(RightStmt) {
			return lexRightStmt
		}

		if l.StartsWith(RightComm) {
			return lexRightComm
		}

		if l.next() == 0 {
			return nil
		}
	}
}

func lexExpr(l *Lexer) stateFn {
	r := l.peek() // change to next
	if r == 0 {
		return nil
	}

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

	if l.StartsWith(RightStmt) {
		return lexRightStmt
	}

	// l.emitError(ErrUnknown)
	return nil
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
		switch r := l.peek(); {
		case unicode.IsDigit(r):
			l.next()
		case r == '.':
			l.emit(TokenFloat)
			l.next()
			// l.emitError(ErrUnexpected)
			return lexLineWhitespace(lexExpr)
		default:
			l.emit(TokenFloat)
			return lexLineWhitespace(lexPipelineOrRightExpr)
		}
	}
}

func lexString(l *Lexer) stateFn {
	if l.next() == 0 { // skips leading "
		return nil
	}

	for {
		r := l.next()
		if r == 0 {
			return nil
		}

		if r == '"' || r == '\'' {
			l.emit(TokenString)
			return lexLineWhitespace(lexPipelineOrRightExpr)
		}
	}

	// for r := l.next(); r != 0; {
	// 	if r == '"' || r == '\'' {
	// 		l.emit(TokenString)
	// 		return lexLineWhitespace(lexPipelineOrRightExpr)
	// 	}
	// }
	// return nil
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
	if l.peek() == 0 {
		return nil
	}

	if l.StartsWith(Pipe) {
		l.pos += len(Pipe)
		l.emit(TokenPipe)
		return lexLineWhitespace(lexFilterIdentifier)
	}

	if l.StartsWith(RightExpr) {
		return lexRightExpr
	}

	if l.StartsWith(RightStmt) {
		return lexRightStmt
	}

	return lexText
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
				l.next()
				l.emit(TokenError)
				return lexText
			default:
				l.emit(TokenWhitespace)
				return nextState
			}
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

func lexStmt(l *Lexer) stateFn {
	r := l.peek() // change to next
	if r == 0 {
		// l.emitError(ErrUnexpectedEOF)
		return nil
	}

	if l.StartsWith(IfStmt) {
		return lexIf
	}

	if l.StartsWith(SwitchStmt) {
		return lexSwitch
	}

	if l.StartsWith(CaseStmt) {
		return lexCase
	}

	if l.StartsWith(DefaultStmt) {
		return lexDefault
	}

	if l.StartsWith(RightStmt) {
		return lexRightStmt
	}

	// l.emitError(ErrUnknown)
	return nil
}

func lexIf(l *Lexer) stateFn {
	l.pos += len(IfStmt)
	l.emit(TokenIfStmt)
	return lexWhitespace(lexExpr)
}

func lexSwitch(l *Lexer) stateFn {
	l.pos += len(SwitchStmt)
	l.emit(TokenSwitchStmt)
	return lexWhitespace(lexExpr)
}

func lexCase(l *Lexer) stateFn {
	l.pos += len(CaseStmt)
	l.emit(TokenCaseStmt)
	return lexWhitespace(lexExpr)
}

func lexDefault(l *Lexer) stateFn {
	l.pos += len(DefaultStmt)
	l.emit(TokenDefaultStmt)
	return lexWhitespace(lexRightStmt)
}

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

func lexLeftStmt(l *Lexer) stateFn {
	l.pos += len(LeftStmt)
	l.emit(TokenLeftStmt)
	return lexLineWhitespace(lexStmt)
}

func lexRightStmt(l *Lexer) stateFn {
	l.pos += len(RightStmt)
	l.emit(TokenRightStmt)
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
			// l.emitError(ErrUnexpectedEOF)
			return nil
		}

		if l.StartsWith(RightComm) {
			l.emit(TokenCommText)
			return lexRightComm
		}

		l.next()
	}
}
