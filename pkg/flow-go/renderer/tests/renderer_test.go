package renderer_test

import (
	"testing"

	"github.com/flowtemplates/cli/pkg/flow-go/parser"
	"github.com/flowtemplates/cli/pkg/flow-go/renderer"
)

func TestRenderer(t *testing.T) {
	testCases := []testCase{
		{
			name:     "Plain text",
			str:      "Hello world",
			expected: "Hello world",
			input: []parser.Node{
				parser.Text{
					Val: "Hello world",
				},
			},
			context:     renderer.Scope{},
			errExpected: false,
		},
		{
			name:     "Expression with string var",
			str:      "{{name}}",
			expected: "useuse",
			input: []parser.Node{
				parser.ExprBlock{
					Body: parser.Ident{Name: "name"},
				},
			},
			context: renderer.Scope{
				"name": "useuse",
			},
			errExpected: false,
		},
		{
			name:     "Multiple expressions",
			str:      "Hello {{name}}!\nFrom {{ flow }} templates",
			expected: "Hello world!\nFrom flow templates",
			input: []parser.Node{
				parser.Text{
					Val: "Hello ",
				},
				parser.ExprBlock{
					Body: parser.Ident{Name: "name"},
				},
				parser.Text{
					Val: "!\nFrom ",
				},
				parser.ExprBlock{
					Body: parser.Ident{Name: "flow"},
				},
				parser.Text{
					Val: " templates",
				},
			},
			context: renderer.Scope{
				"name": "world",
				"flow": "flow",
			},
			errExpected: false,
		},
		{
			name:     "Truthy if statement",
			str:      "{%if var%}\ntext\n{%end%}",
			expected: "\ntext\n",
			input: []parser.Node{
				parser.IfStmt{
					PostKwWs: " ",
					Condition: parser.Ident{
						Name:   "var",
						PostWs: "",
					},
					Body: []parser.Node{
						parser.Text{
							Val: "\ntext\n",
						},
					},
					Else: nil,
				},
			},
			context: renderer.Scope{
				"var": "true",
			},
			errExpected: false,
		},
		{
			name:     "Falsy if statement",
			str:      "{%if var%}\ntext\n{%end%}",
			expected: "",
			input: []parser.Node{
				parser.IfStmt{
					PostKwWs: " ",
					Condition: parser.Ident{
						Name:   "var",
						PostWs: "",
					},
					Body: []parser.Node{
						parser.Text{
							Val: "\ntext\n",
						},
					},
					Else: nil,
				},
			},
			context: renderer.Scope{
				"var": "false",
			},
			errExpected: false,
		},
	}
	runTestCases(t, testCases)
}
