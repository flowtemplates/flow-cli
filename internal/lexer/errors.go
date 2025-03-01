package lexer

import "errors"

var (
	ErrUnexpectedEOF = errors.New("unexpected EOF")
	ErrUnexpected    = errors.New("unexpected")
	ErrUnclosedExpr  = errors.New("unclosed expression")
	ErrUnknown       = errors.New("unknown")
)
