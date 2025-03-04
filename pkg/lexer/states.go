package lexer

import (
	"strings"
	"unicode"

	"github.com/flowtemplates/cli/pkg/token"
)

type stateFn func(*Lexer) stateFn

func (l *Lexer) lexToken(t token.Type, next stateFn) stateFn {
	l.pos += len(token.TokenString(t))
	l.emit(t)
	return next
}

func (l *Lexer) StartsWith(t token.Type) bool {
	return strings.HasPrefix(l.input[l.pos:], token.TokenString(t))
}

func (l *Lexer) tryTokens(nextState stateFn, tokens ...token.Type) stateFn {
	for _, token := range tokens {
		if l.StartsWith(token) {
			return l.lexToken(token, nextState)
		}
	}

	return nil
}

func lexText(l *Lexer) stateFn {
	for {
		if l.StartsWith(token.LEXPR) {
			l.emit(token.TEXT)
			return l.lexToken(token.LEXPR, lexExpr)
		}

		if l.StartsWith(token.RARR) {
			l.emit(token.TEXT)
			return l.lexToken(token.RARR, lexComm)
		}

		if l.StartsWith(token.LSTMT) {
			l.emit(token.TEXT)
			return l.lexToken(token.LSTMT, lexStmt)
		}

		if l.StartsWith(token.LCOMM) {
			l.emit(token.TEXT)
			return l.lexToken(token.LCOMM, lexComm)
		}

		// if l.StartsWith(RightExpr) {
		// 	return lexRightExpr
		// }

		// if l.StartsWith(RightStmt) {
		// 	return lexRightStmt
		// }

		// if l.StartsWith(RightComm) {
		// 	return lexRightComm
		// }

		if l.next() == eof {
			break
		}
	}

	l.emit(token.TEXT)
	return nil
}

// TODO: rename
func lexRealExpr(nextState stateFn) stateFn {
	return func(l *Lexer) stateFn {
		switch r := l.next(); {
		case r == eof:
			return nil
		case r == '\n' || r == '\r':
			return lexText
		case unicode.IsSpace(r):
			return lexLineWhitespace(nextState)
		case r == '"' || r == '\'':
			return lexString
		case r == token.TokenRune(token.LPAREN):
			l.back()
			return l.lexToken(token.LPAREN, nextState)
		case r == token.TokenRune(token.RPAREN):
			l.back()
			return l.lexToken(token.RPAREN, nextState)
		case r == token.TokenRune(token.SUB):
			l.back()
			return l.lexToken(token.SUB, nextState)
		case r == token.TokenRune(token.ADD):
			l.back()
			return l.lexToken(token.ADD, nextState)
		case r == token.TokenRune(token.MUL):
			l.back()
			return l.lexToken(token.MUL, nextState)
		case r == token.TokenRune(token.DIV):
			l.back()
			return l.lexToken(token.DIV, nextState)
		case unicode.IsDigit(r):
			return lexNum(nextState)
		case token.IsNotOp(r) && r != '.':
			return lexIdent(nextState)
		default:
			l.emit(token.EXPECTED_EXPR)
			return nextState
		}
	}
}

func lexExpr(l *Lexer) stateFn {
	if l.StartsWith(token.REXPR) {
		return l.lexToken(token.REXPR, lexText)
	}

	if state := l.tryTokens(lexExpr,
		token.RARR,
		token.EQL,
		token.NEQL,
		token.IS,
		token.NOT,
		token.AND,
		token.EXCL,
		token.OR,
		token.LAND,
		token.LOR,
		token.LESS,
		token.GTR,
		token.QUESTION,
		token.COLON,
	); state != nil {
		return state
	}

	return lexRealExpr(lexExpr)
}

func lexComm(l *Lexer) stateFn {
	// ? try to lex something to do not cause commenting whole thing if there is no closing tag
	for {
		if l.StartsWith(token.RCOMM) {
			l.emit(token.COMM_TEXT)
			return l.lexToken(token.RCOMM, lexText)
		}

		r := l.next()
		if r == eof {
			l.emit(token.COMM_TEXT)
			return nil
		}
	}
}

func lexNum(nextState stateFn) stateFn {
	return func(l *Lexer) stateFn {
		digits := "0123456789"

		l.acceptRun(digits)
		if l.accept(".") {
			l.acceptRun(digits)
			l.emit(token.FLOAT)
		} else {
			l.emit(token.INT)
		}

		return nextState
	}
}
func lexString(l *Lexer) stateFn {
	for {
		r := l.next()
		switch r {
		case eof:
			l.emit(token.NOT_TERMINATED_STR)
			return lexText
		case '"', '\'':
			l.emit(token.STR)
			return lexExpr
		}
	}
}

func lexIdent(nextState stateFn) stateFn {
	return func(l *Lexer) stateFn {
		for {
			switch r := l.next(); {
			case r == eof:
				l.emit(token.IDENT)
				return nil
			case !token.IsNotOp(r) || unicode.IsSpace(r):
				l.back()
				l.emit(token.IDENT)
				return nextState
			}
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
				return lexText
			default:
				l.emit(token.WS)
				return nextState
			}
		}
	}
}

func lexStmt(l *Lexer) stateFn {
	if l.StartsWith(token.RSTMT) {
		return l.lexToken(token.RSTMT, lexText)
	}

	if state := l.tryTokens(lexStmt,
		token.RARR,
		token.IF,
		token.SWITCH,
		token.CASE,
		token.DEFAULT,
		token.END,
		token.EQL,
		token.NEQL,
		token.IS,
		token.NOT,
		token.EXCL,
		token.AND,
		token.OR,
		token.LAND,
		token.LOR,
		token.LESS,
		token.GTR,
		token.QUESTION,
		token.COLON,
	); state != nil {
		return state
	}

	return lexRealExpr(lexStmt)
}
