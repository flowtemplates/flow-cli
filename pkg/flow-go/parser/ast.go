package parser

import (
	"github.com/flowtemplates/cli/pkg/flow-go/token"
)

type Node interface{} // nolint: iface

type Expr interface{} // nolint: iface

type (
	Text struct {
		Pos token.Position
		Val string
	}

	Lit struct {
		Pos    token.Position
		Typ    token.Type
		Val    string
		PostWS string
	}

	Ident struct {
		Pos    token.Position
		Name   string
		PostWS string
	}

	BinaryExpr struct {
		X        Expr
		OpPos    token.Position
		PostOpWS string
		Op       token.Type
		Y        Expr
	}

	ExprBlock struct {
		LBrace  token.Position
		PostLWS string
		Body    Expr
		RBrace  token.Position
	}

	IfStmt struct {
		StmtBeg    token.Position
		PostStmtWs string
		IfPos      token.Position
		PostIfWs   string
		Condition  Expr
		Body       []Node
		Else       []Node
		StmtEnd    token.Position
	}
)
