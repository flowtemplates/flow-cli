package lexer

import "fmt"

const (
	RightExpr = "}}"
	LeftExpr  = "{{"
	LeftComm  = "{#"
	RightComm = "#}"
	LeftStmt  = "{%"
	RightStmt = "%}"
)

type TokenType int

const (
	TokenEOF TokenType = iota
	TokenError
	TokenWhitespace
	TokenText
	TokenSymbol
	TokenInteger
	TokenString
	TokenBoolean
	TokenFloat
	TokenLeftExpr  // {{
	TokenRightExpr // }}
	TokenLeftComm  // {#
	TokenRightComm // #}
	TokenLeftStmt  // {%
	TokenRightStmt // %}
)

type Token struct {
	Typ TokenType
	Val string
}

func (t Token) String() string {
	return fmt.Sprintf("{Type: %v, Value: %q}", t.Typ, t.Val)
}
