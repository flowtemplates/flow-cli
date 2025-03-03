package parser_test

import (
	"testing"

	"github.com/templatesflow/cli/internal/parser"
	"github.com/templatesflow/cli/internal/token"
)

func TestParser(t *testing.T) {
	testCases := []testCase{
		{
			name: "Plain text",
			str:  "hello",
			input: []token.Token{
				{Typ: token.TEXT, Val: "hello"},
			},
			expected: []parser.Node{
				parser.Text{
					Pos: 0,
					Val: "hello",
				},
			},
		},
		{
			name: "Single expression with var",
			str:  "{{x}}",
			input: []token.Token{
				{Typ: token.LEXPR},
				{Typ: token.IDENT, Val: "x"},
				{Typ: token.REXPR},
			},
			expected: []parser.Node{
				parser.ExprBlock{
					LBrace: 0,
					Body:   &parser.Ident{Pos: 0, Name: "x"},
					RBrace: 0,
				},
			},
		},
		{
			name: "Whitespaces with var",
			str:  "{{ x }}",
			input: []token.Token{
				{Typ: token.LEXPR},
				{Typ: token.WS, Val: " "},
				{Typ: token.IDENT, Val: "x"},
				{Typ: token.WS, Val: " "},
				{Typ: token.REXPR},
			},
			expected: []parser.Node{
				parser.ExprBlock{
					LBrace: 0,
					Body: &parser.Ident{
						PostWS: " ",
						Pos:    0,
						Name:   "x",
					},
					PostLWS: " ",
					RBrace:  0,
				},
			},
		},
		{
			name: "Expressions between text",
			str:  "Hello, {{name}}\n{{var }}Text",
			input: []token.Token{
				{Typ: token.LEXPR},
				{Typ: token.IDENT, Val: "x"},
				{Typ: token.REXPR},
			},
			expected: []parser.Node{
				parser.ExprBlock{
					LBrace: 0,
					Body:   &parser.Ident{Pos: 0, Name: "x"},
					RBrace: 0,
				},
			},
		},
		{
			name: "Addition",
			str:  "{{123+age}}",
			input: []token.Token{
				{Typ: token.LEXPR},
				{Typ: token.INT, Val: "123"},
				{Typ: token.ADD},
				{Typ: token.IDENT, Val: "age"},
				{Typ: token.REXPR},
			},
			expected: []parser.Node{
				parser.ExprBlock{
					LBrace: 0,
					Body: &parser.BinaryExpr{
						X: &parser.Lit{
							Pos: 0,
							Val: "123",
							Typ: token.INT,
						},
						OpPos: 0,
						Op:    token.ADD,
						Y: &parser.Ident{
							Pos:  0,
							Name: "age",
						},
					},
					RBrace: 0,
				}},
		},
		{
			name: "Addition with whitespaces",
			str:  "{{123 + age}}",
			input: []token.Token{
				{Typ: token.LEXPR},
				{Typ: token.INT, Val: "123"},
				{Typ: token.WS, Val: " "},
				{Typ: token.ADD},
				{Typ: token.WS, Val: " "},
				{Typ: token.IDENT, Val: "age"},
				{Typ: token.REXPR},
			},
			expected: []parser.Node{
				parser.ExprBlock{
					LBrace: 0,
					Body: &parser.BinaryExpr{
						X: &parser.Lit{
							Pos:    0,
							Val:    "123",
							PostWS: " ",
							Typ:    token.INT,
						},
						OpPos:    0,
						PostOpWS: " ",
						Op:       token.ADD,
						Y: &parser.Ident{
							Pos:  0,
							Name: "age",
						},
					},
					RBrace: 0,
				}},
		},
		{
			name: "Subtraction",
			str:  "{{123-age}}",
			input: []token.Token{
				{Typ: token.LEXPR},
				{Typ: token.INT, Val: "123"},
				{Typ: token.SUB},
				{Typ: token.IDENT, Val: "age"},
				{Typ: token.REXPR},
			},
			expected: []parser.Node{
				parser.ExprBlock{
					LBrace: 0,
					Body: &parser.BinaryExpr{
						X: &parser.Lit{
							Pos: 0,
							Val: "123",
							Typ: token.INT,
						},
						OpPos: 0,
						Op:    token.SUB,
						Y: &parser.Ident{
							Pos:  0,
							Name: "age",
						},
					},
					RBrace: 0,
				}},
		},
		{
			name: "Multiply",
			str:  "{{123*age}}",
			input: []token.Token{
				{Typ: token.LEXPR},
				{Typ: token.INT, Val: "123"},
				{Typ: token.MUL},
				{Typ: token.IDENT, Val: "age"},
				{Typ: token.REXPR},
			},
			expected: []parser.Node{
				parser.ExprBlock{
					LBrace: 0,
					Body: &parser.BinaryExpr{
						X: &parser.Lit{
							Pos: 0,
							Val: "123",
							Typ: token.INT,
						},
						OpPos: 0,
						Op:    token.MUL,
						Y: &parser.Ident{
							Pos:  0,
							Name: "age",
						},
					},
					RBrace: 0,
				}},
		},

		{
			name: "Multiply",
			str:  "{{123*age}}",
			input: []token.Token{
				{Typ: token.LEXPR},
				{Typ: token.INT, Val: "123"},
				{Typ: token.MUL},
				{Typ: token.IDENT, Val: "age"},
				{Typ: token.REXPR},
			},
			expected: []parser.Node{
				parser.ExprBlock{
					LBrace: 0,
					Body: &parser.BinaryExpr{
						X: &parser.Lit{
							Pos: 0,
							Val: "123",
							Typ: token.INT,
						},
						OpPos: 0,
						Op:    token.MUL,
						Y: &parser.Ident{
							Pos:  0,
							Name: "age",
						},
					},
					RBrace: 0,
				},
			},
		},
		{
			name: "Division",
			str:  "{{age/2}}",
			input: []token.Token{
				{Typ: token.LEXPR},
				{Typ: token.IDENT, Val: "age"},
				{Typ: token.DIV},
				{Typ: token.INT, Val: "2"},
				{Typ: token.REXPR},
			},
			expected: []parser.Node{
				parser.ExprBlock{
					LBrace: 0,
					Body: &parser.BinaryExpr{
						X: &parser.Ident{
							Pos:  0,
							Name: "age",
						},
						OpPos: 0,
						Op:    token.DIV,
						Y: &parser.Lit{
							Pos: 0,
							Val: "2",
							Typ: token.INT,
						},
					},
					RBrace: 0,
				},
			},
		},
		{
			name: "Redundant parens",
			str:  "{{(age)}}",
			input: []token.Token{
				{Typ: token.LEXPR},
				{Typ: token.LPAREN},
				{Typ: token.IDENT, Val: "age"},
				{Typ: token.RPAREN},
				{Typ: token.REXPR},
			},
			expected: []parser.Node{
				parser.ExprBlock{
					LBrace: 0,
					Body: &parser.Ident{
						Pos:  0,
						Name: "age",
					},
					RBrace: 0,
				},
			},
		},
		{
			name: "Parens * int",
			str:  "{{(1+2)*3}}",
			input: []token.Token{
				{Typ: token.LEXPR},
				{Typ: token.LPAREN},
				{Typ: token.INT, Val: "1"},
				{Typ: token.ADD},
				{Typ: token.INT, Val: "2"},
				{Typ: token.RPAREN},
				{Typ: token.MUL},
				{Typ: token.INT, Val: "3"},
				{Typ: token.REXPR},
			},
			expected: []parser.Node{
				parser.ExprBlock{
					LBrace: 0,
					Body: &parser.BinaryExpr{
						X: &parser.BinaryExpr{
							X: &parser.Lit{
								Pos: 0,
								Val: "1",
								Typ: token.INT,
							},
							OpPos: 0,
							Op:    token.ADD,
							Y: &parser.Lit{
								Pos: 0,
								Val: "2",
								Typ: token.INT,
							},
						},
						OpPos: 0,
						Op:    token.MUL,
						Y: &parser.Lit{
							Pos: 0,
							Val: "3",
							Typ: token.INT,
						},
					},
					RBrace: 0,
				},
			},
		},
		{
			name: "Parens * int (with whitespaces)",
			str:  "{{(1 + 2) * 3}}",
			input: []token.Token{
				{Typ: token.LEXPR},
				{Typ: token.LPAREN},
				{Typ: token.INT, Val: "1"},
				{Typ: token.WS, Val: " "},
				{Typ: token.ADD},
				{Typ: token.WS, Val: " "},
				{Typ: token.INT, Val: "2"},
				{Typ: token.RPAREN},
				{Typ: token.WS, Val: " "},
				{Typ: token.MUL},
				{Typ: token.WS, Val: " "},
				{Typ: token.INT, Val: "3"},
				{Typ: token.REXPR},
			},
			expected: []parser.Node{
				parser.ExprBlock{
					LBrace: 0,
					Body: &parser.BinaryExpr{
						X: &parser.BinaryExpr{
							X: &parser.Lit{
								PostWS: " ",
								Pos:    0,
								Val:    "1",
								Typ:    token.INT,
							},
							OpPos:    0,
							PostOpWS: " ",
							Op:       token.ADD,
							Y: &parser.Lit{
								Pos: 0,
								Val: "2",
								Typ: token.INT,
							},
						},
						OpPos:    0,
						PostOpWS: " ",
						Op:       token.MUL,
						Y: &parser.Lit{
							Pos: 0,
							Val: "3",
							Typ: token.INT,
						},
					},
					RBrace: 0,
				},
			},
		},
	}
	runTestCases(t, testCases)
}
