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
					Body: &parser.Ident{Name: "name"},
				},
			},
			context: renderer.Scope{
				"name": "useuse",
			},
			errExpected: false,
		},
		// {
		// 	name:     "If statement",
		// 	str:      "{%if var%}\ntext\n{%end%}",
		// 	expected: "\ntext\n",
		// 	input: []parser.Node{
		// 		parser.IfStmt{
		// 			StmtBeg:  0,
		// 			IfPos:    0,
		// 			PostIfWs: " ",
		// 			Condition: parser.Ident{
		// 				Pos:    0,
		// 				Name:   "var",
		// 				PostWS: "",
		// 			},
		// 			Body: []parser.Node{
		// 				parser.Text{
		// 					Pos: 0,
		// 					Val: "\ntext\n",
		// 				},
		// 			},
		// 			Else:    nil,
		// 			StmtEnd: 0,
		// 		},
		// 	},
		// 	context: renderer.Scope{
		// 		"var": true,
		// 	},
		// 	errExpected: false,
		// },
	}
	runTestCases(t, testCases)
}
