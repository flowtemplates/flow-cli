package analyzer_test

import (
	"testing"

	"github.com/flowtemplates/cli/pkg/flow-go/analyzer"
	"github.com/flowtemplates/cli/pkg/flow-go/parser"
	"github.com/flowtemplates/cli/pkg/flow-go/token"
	"github.com/flowtemplates/cli/pkg/flow-go/types"
)

func TestGetTypeMap(t *testing.T) {
	testCases := []testCase{
		{
			name: "Plain text",
			str:  "Hello world",
			input: []parser.Node{
				parser.Text{
					Pos: 0,
					Val: "Hello world",
				},
			},
			expected:    analyzer.TypeMap{},
			errExpected: false,
		},
		{
			name: "Single var",
			str:  "{{name}}",
			input: []parser.Node{
				parser.ExprBlock{
					LBrace: 0,
					Body:   parser.Ident{Pos: 0, Name: "name"},
					RBrace: 0,
				},
			},
			expected: analyzer.TypeMap{
				"name": types.String,
			},
			errExpected: false,
		},
		{
			name: "Var + integer literal",
			str:  "{{age+123}}",
			input: []parser.Node{
				parser.ExprBlock{
					LBrace: 0,
					Body: parser.BinaryExpr{
						X: parser.Ident{
							Pos:  0,
							Name: "age",
						},
						OpPos: 0,
						Op:    token.ADD,
						Y: parser.Lit{
							Pos: 0,
							Val: "123",
							Typ: token.INT,
						},
					},
					RBrace: 0,
				},
			},
			expected: analyzer.TypeMap{
				"age": types.Number,
			},
			errExpected: false,
		},
		{
			name: "Integer literal + var",
			str:  "{{123+age}}",
			input: []parser.Node{
				parser.ExprBlock{
					LBrace: 0,
					Body: parser.BinaryExpr{
						X: parser.Lit{
							Pos: 0,
							Val: "123",
							Typ: token.INT,
						},
						OpPos: 0,
						Op:    token.ADD,
						Y: parser.Ident{
							Pos:  0,
							Name: "age",
						},
					},
					RBrace: 0,
				},
			},
			expected: analyzer.TypeMap{
				"age": types.Number,
			},
			errExpected: false,
		},
		{
			name: "Var + var",
			str:  "{{age+time}}",
			input: []parser.Node{
				parser.ExprBlock{
					LBrace: 0,
					Body: parser.BinaryExpr{
						X: parser.Ident{
							Pos:  0,
							Name: "age",
						},
						OpPos: 0,
						Op:    token.ADD,
						Y: parser.Ident{
							Pos:  0,
							Name: "time",
						},
					},
					RBrace: 0,
				},
			},
			expected: analyzer.TypeMap{
				"age":  types.Any,
				"time": types.Any,
			},
			errExpected: false,
		},
		{
			name: "Var + string literal",
			str:  "{{name+'ish'}}",
			input: []parser.Node{
				parser.ExprBlock{
					LBrace: 0,
					Body: parser.BinaryExpr{
						X: parser.Ident{
							Pos:  0,
							Name: "name",
						},
						OpPos: 0,
						Op:    token.ADD,
						Y: parser.Lit{
							Pos: 0,
							Val: "'ish'",
							Typ: token.STR,
						},
					},
					RBrace: 0,
				},
			},
			expected: analyzer.TypeMap{
				"name": types.String,
			},
			errExpected: false,
		},
	}
	runTestCases(t, testCases)
}
