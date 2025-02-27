package lexer

import "fmt"

type TokenType int

const (
	TokenEOF TokenType = iota
	TokenError
	TokenWhitespace
	TokenText
	TokenSymbol
	TokenLeftExpr     // {{
	TokenRightExpr    // }}
	TokenLeftComment  // {#
	TokenRightComment // #}
	TokenLeftStmt     // {%
	TokenRightStmt    // %}
)

type Token struct {
	typ TokenType
	val string
}

func (t Token) String() string {
	return fmt.Sprintf("{Type: %v, Value: %q}", t.typ, t.val)
}
