package token

import (
	"fmt"
	"slices"
)

type Type int

func (t Type) String() string {
	return TokenString(t)
}

const (
	EOF Type = iota
	ILLEGAL

	valueable_beg
	COMM_TEXT
	TEXT
	WS

	IDENT // main
	INT   // 12345
	FLOAT // 123.45
	STR   // "abc"

	errors_beg
	// Errors
	NOT_TERMINATED_STR
	EXPECTED_EXPR
	errors_end
	valueable_end

	operator_beg
	// Operators and delimiters
	ADD // +
	SUB // -
	MUL // *
	DIV // /
	MOD // %

	ASSIGN     // =
	ADD_ASSIGN // +=
	SUB_ASSIGN // -=
	MUL_ASSIGN // *=
	QUO_ASSIGN // /=
	REM_ASSIGN // %=

	QUESTION // ?
	COLON    // :

	comparison_op_beg
	EQL  // ==
	LESS // <
	GTR  // >
	EXCL // !
	NEQL // !=
	LEQ  // <=
	GEQ  // >=
	LAND // &&
	LOR  // ||
	comparison_op_end

	LPAREN // (
	LBRACK // [
	LBRACE // {

	RPAREN // )
	RBRACK // ]
	RBRACE // }

	COMMA // ,

	LEXPR // {{
	REXPR // }}
	LCOMM // {#
	RCOMM // #}
	LSTMT // {%
	RSTMT // %}

	RARR // ->
	operator_end

	keyword_beg
	// Keywords
	FOR     // for
	LET     // let
	IF      // if
	GENIF   // genif
	ELSE    // else
	SWITCH  // switch
	END     // end
	CASE    // case
	DEFAULT // default
	EXTEND  // extend
	AND     // and
	OR      // or
	IS      // is
	NOT     // not
	DO      // do
	keyword_end
)

var tokens = []string{
	EOF:     "EOF",
	ILLEGAL: "ILLEGAL",

	COMM_TEXT: "COMMENT",
	TEXT:      "TEXT",
	WS:        "WHITESPACE",

	IDENT: "IDENT",
	INT:   "INT",
	FLOAT: "FLOAT",
	STR:   "STRING",

	NOT_TERMINATED_STR: "NOT_TERMINATED_STR",
	EXPECTED_EXPR:      "EXPECTED_EXPR",

	QUESTION: "?",
	COLON:    ":",

	ADD: "+",
	SUB: "-",
	MUL: "*",
	DIV: "/",
	MOD: "%",

	ASSIGN:     "=",
	ADD_ASSIGN: "+=",
	SUB_ASSIGN: "-=",
	MUL_ASSIGN: "*=",
	QUO_ASSIGN: "/=",
	REM_ASSIGN: "%=",

	EQL:  "==",
	LESS: "<",
	GTR:  ">",
	EXCL: "!",
	NEQL: "!=",
	LEQ:  "<=",
	GEQ:  ">=",
	LAND: "&&",
	LOR:  "||",

	LPAREN: "(",
	LBRACK: "[",
	LBRACE: "{",

	RPAREN: ")",
	RBRACK: "]",
	RBRACE: "}",

	COMMA: ",",

	LEXPR: "{{",
	REXPR: "}}",
	LCOMM: "{#",
	RCOMM: "#}",
	LSTMT: "{%",
	RSTMT: "%}",

	RARR: "->",

	FOR:     "for",
	LET:     "let",
	IF:      "if",
	GENIF:   "genif",
	ELSE:    "else",
	SWITCH:  "switch",
	END:     "end",
	CASE:    "case",
	DEFAULT: "default",
	EXTEND:  "extend",
	AND:     "and",
	OR:      "or",
	IS:      "is",
	NOT:     "not",
	DO:      "do",
}

func TokenString(t Type) string {
	return tokens[t]
}

func TokenRune(t Type) rune {
	return rune(tokens[t][0])
}

type Token struct {
	Typ Type
	Val string
	Pos Position
}

func (t Token) String() string {
	if t.IsValueable() {
		switch t.Typ {
		case EOF:
			return "EOF"
		case TEXT:
			return fmt.Sprintf("{Typ: %s, Val: %.10q, Pos: %s}", TokenString(t.Typ), t.Val, t.Pos)
		default:
			return fmt.Sprintf("{Typ: %s, Val: %q, Pos: %s}", TokenString(t.Typ), t.Val, t.Pos)
		}
	}

	return fmt.Sprintf("{Typ: %[1]s, Val: %[1]q, Pos: %s}", TokenString(t.Typ), t.Pos)
}

func (t Token) IsOneOfMany(types ...Type) bool {
	return slices.Contains(types, t.Typ)
}

func (t Token) IsValueable() bool {
	return valueable_beg < t.Typ && t.Typ < valueable_end
}

func (t Token) IsComparisonOp() bool {
	return comparison_op_beg < t.Typ && t.Typ < comparison_op_end
}

func IsNotOp(r rune) bool {
	for i := operator_beg + 1; i < operator_end; i++ {
		t := tokens[i]
		if t != "" && r == rune(t[0]) {
			return false
		}
	}

	return true
}
