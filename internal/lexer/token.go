package lexer

import "fmt"

const (
	FalseLiteral = "false"
	TrueLiteral  = "true"
	RightExpr    = "}}"
	LeftExpr     = "{{"
	LeftComm     = "{#"
	RightComm    = "#}"
	LeftStmt     = "{%"
	RightStmt    = "%}"
	Pipe         = "->"
)

type TokenType int

const (
	TokenEOF TokenType = iota
	TokenError
	TokenWhitespace
	TokenText
	TokenIdentifier
	TokenFilterIdentifier
	TokenInteger
	TokenString
	TokenBoolean
	TokenFloat
	TokenPipe
	TokenCommText
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
