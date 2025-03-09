package parser

import (
	"github.com/flowtemplates/cli/pkg/flow-go/token"
)

type Node interface{} // nolint: iface

type Expr interface{} // nolint: iface

type Stmt interface{} // nolint: iface

type (
	Text struct {
		Pos token.Position
		Val string
	}

	Lit struct {
		Pos    token.Position
		Typ    token.Type
		Val    string
		PostWs string
	}

	Ident struct {
		Pos    token.Position
		Name   string
		PostWs string
	}

	BinaryExpr struct {
		X        Expr
		OpPos    token.Position
		PostOpWs string
		Op       token.Type
		Y        Expr
	}

	ExprBlock struct {
		LBrace  token.Position
		PostLWs string
		Body    Expr
		RBrace  token.Position
	}

	IfStmt struct {
		StmtBeg    token.Position
		PostStmtWs string
		KwPos      token.Position
		PostKwWs   string
		Condition  Expr
		Body       []Node
		Else       []Node
		StmtEnd    token.Position
	}

	GenIfStmt struct {
		StmtBeg    token.Position
		PostStmtWs string
		KwPos      token.Position
		PostKwWs   string
		Condition  Expr
		StmtEnd    token.Position
	}
)
