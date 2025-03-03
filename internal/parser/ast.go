package parser

import "github.com/templatesflow/cli/internal/token"

type Node interface {
}

type Expr interface {
}

type (
	Text struct {
		Pos int
		Val string
	}

	Lit struct {
		Pos    int
		Typ    token.Type
		Val    string
		PostWS string
	}

	Ident struct {
		Pos    int
		Name   string
		PostWS string
	}

	BinaryExpr struct {
		X        Expr
		OpPos    int
		PostOpWS string
		Op       token.Type
		Y        Expr
	}

	ExprBlock struct {
		LBrace  int
		PostLWS string
		Body    Expr
		RBrace  int
	}
)
