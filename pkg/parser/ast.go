package parser

import "github.com/flowtemplates/cli/pkg/token"

type Node interface{} // nolint: iface

type Expr interface{} // nolint: iface

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

	IfStmt struct {
		StmtBeg    int
		PostStmtWs string
		IfPos      int
		PostIfWs   string
		Condition  Node
		Body       []Node
		Else       *[]Node
		StmtEnd    int
	}
)
