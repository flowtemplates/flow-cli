package lexer

import "fmt"

const (
	FalseLiteral = "false"
	TrueLiteral  = "true"
	IfStmt       = "if"
	SwitchStmt   = "switch"
	CaseStmt     = "case"
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
	TokenIfStmt
	TokenSwitchStmt
	TokenCaseStmt
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
	Pos int
}

func (t Token) String() string {
	return fmt.Sprintf("{Type: %v, Value: %q, Pos: %d}", t.Typ, t.Val, t.Pos)
}
