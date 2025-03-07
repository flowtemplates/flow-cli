package parser_test

import (
	"testing"

	"github.com/flowtemplates/cli/pkg/flow-go/parser"
	"github.com/flowtemplates/cli/pkg/flow-go/token"
)

func TestExpressions(t *testing.T) {
	testCases := []testCase{
		{
			name: "Plain text",
			str:  "hello",
			input: []token.Token{
				{Typ: token.TEXT, Val: "hello"},
			},
			expected: []parser.Node{
				parser.Text{
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
					Body: &parser.Ident{Name: "x"},
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
					Body: &parser.Ident{
						PostWS: " ",
						Name:   "x",
					},
					PostLWS: " ",
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
					Body: &parser.Ident{Name: "x"},
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
					Body: &parser.BinaryExpr{
						X: &parser.Lit{
							Val: "123",
							Typ: token.INT,
						},
						Op: token.ADD,
						Y: &parser.Ident{
							Name: "age",
						},
					},
				},
			},
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
					Body: &parser.BinaryExpr{
						X: &parser.Lit{
							Val:    "123",
							PostWS: " ",
							Typ:    token.INT,
						},
						PostOpWS: " ",
						Op:       token.ADD,
						Y: &parser.Ident{
							Name: "age",
						},
					},
				},
			},
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
					Body: &parser.BinaryExpr{
						X: &parser.Lit{
							Val: "123",
							Typ: token.INT,
						},
						Op: token.SUB,
						Y: &parser.Ident{
							Name: "age",
						},
					},
				},
			},
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
					Body: &parser.BinaryExpr{
						X: &parser.Lit{
							Val: "123",
							Typ: token.INT,
						},
						Op: token.MUL,
						Y: &parser.Ident{
							Name: "age",
						},
					},
				},
			},
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
					Body: &parser.BinaryExpr{
						X: &parser.Lit{
							Val: "123",
							Typ: token.INT,
						},
						Op: token.MUL,
						Y: &parser.Ident{
							Name: "age",
						},
					},
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
					Body: &parser.BinaryExpr{
						X: &parser.Ident{
							Name: "age",
						},
						Op: token.DIV,
						Y: &parser.Lit{
							Val: "2",
							Typ: token.INT,
						},
					},
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
					Body: &parser.Ident{
						Name: "age",
					},
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
					Body: &parser.BinaryExpr{
						X: &parser.BinaryExpr{
							X: &parser.Lit{
								Val: "1",
								Typ: token.INT,
							},
							Op: token.ADD,
							Y: &parser.Lit{
								Val: "2",
								Typ: token.INT,
							},
						},
						Op: token.MUL,
						Y: &parser.Lit{
							Val: "3",
							Typ: token.INT,
						},
					},
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
					Body: &parser.BinaryExpr{
						X: &parser.BinaryExpr{
							X: &parser.Lit{
								PostWS: " ",
								Val:    "1",
								Typ:    token.INT,
							},
							PostOpWS: " ",
							Op:       token.ADD,
							Y: &parser.Lit{
								Val: "2",
								Typ: token.INT,
							},
						},
						PostOpWS: " ",
						Op:       token.MUL,
						Y: &parser.Lit{
							Val: "3",
							Typ: token.INT,
						},
					},
				},
			},
		},
	}
	runTestCases(t, testCases)
}

func TestStatements(t *testing.T) {
	testCases := []testCase{
		{
			name: "If statement",
			str:  "{%if var%}\ntext\n{%end%}",
			input: []token.Token{
				{Typ: token.LSTMT},
				{Typ: token.IF},
				{Typ: token.WS, Val: " "},
				{Typ: token.IDENT, Val: "var"},
				{Typ: token.RSTMT},
				{Typ: token.TEXT, Val: "\ntext\n"},
				{Typ: token.LSTMT},
				{Typ: token.END},
				{Typ: token.RSTMT},
			},
			expected: []parser.Node{
				parser.IfStmt{
					PostIfWs: " ",
					Condition: parser.Ident{
						Name:   "var",
						PostWS: "",
					},
					Body: []parser.Node{
						parser.Text{
							Val: "\ntext\n",
						},
					},
					Else: nil,
				},
			},
		},
		{
			name: "If statement (with whitespaces)",
			str:  "{% if var  %}text{% end %}",
			input: []token.Token{
				{Typ: token.LSTMT},
				{Typ: token.WS, Val: " "},
				{Typ: token.IF},
				{Typ: token.WS, Val: " "},
				{Typ: token.IDENT, Val: "var"},
				{Typ: token.WS, Val: "  "},
				{Typ: token.RSTMT},
				{Typ: token.TEXT, Val: "text"},
				{Typ: token.LSTMT},
				{Typ: token.WS, Val: " "},
				{Typ: token.END},
				{Typ: token.WS, Val: " "},
				{Typ: token.RSTMT},
			},
			expected: []parser.Node{
				parser.IfStmt{
					PostStmtWs: " ",
					PostIfWs:   " ",
					Condition: parser.Ident{
						Name:   "var",
						PostWS: "  ",
					},
					Body: []parser.Node{
						parser.Text{
							Val: "text",
						},
					},
					Else: nil,
				},
			},
		},
		{
			name: "Nested if statements",
			str:  "{%if var%}1{%if name%}text{%end%}2{%end%}",
			input: []token.Token{
				{Typ: token.LSTMT},
				{Typ: token.IF},
				{Typ: token.WS, Val: " "},
				{Typ: token.IDENT, Val: "var"},
				{Typ: token.RSTMT},
				{Typ: token.TEXT, Val: "1"},
				{Typ: token.LSTMT},
				{Typ: token.IF},
				{Typ: token.WS, Val: " "},
				{Typ: token.IDENT, Val: "name"},
				{Typ: token.RSTMT},
				{Typ: token.TEXT, Val: "text"},
				{Typ: token.LSTMT},
				{Typ: token.END},
				{Typ: token.RSTMT},
				{Typ: token.TEXT, Val: "2"},
				{Typ: token.LSTMT},
				{Typ: token.END},
				{Typ: token.RSTMT},
			},
			expected: []parser.Node{
				parser.IfStmt{
					PostIfWs: " ",
					Condition: parser.Ident{
						Name: "var",
					},
					Body: []parser.Node{
						parser.Text{
							Val: "1",
						},
						parser.IfStmt{
							PostIfWs: " ",
							Condition: parser.Ident{
								Name: "name",
							},
							Body: []parser.Node{
								parser.Text{
									Val: "text",
								},
							},
						},
						parser.Text{
							Val: "2",
						},
					},
					Else: nil,
				},
			},
		},
	}
	runTestCases(t, testCases)
}
