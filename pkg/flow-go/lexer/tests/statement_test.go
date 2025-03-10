package lexer_test

import (
	"testing"

	"github.com/flowtemplates/cli/pkg/flow-go/token"
)

func TestStatement(t *testing.T) {
	testCases := []testCase{
		{
			name:  "Simple if statement",
			input: "{%if name%}",
			expected: []token.Token{
				{Typ: token.LSTMT},
				{Typ: token.IF},
				{Typ: token.WS, Val: " "},
				{Typ: token.IDENT, Val: "name"},
				{Typ: token.RSTMT},
			},
		},
		{
			name:  "Simple switch statement",
			input: "{%switch name%}",
			expected: []token.Token{
				{Typ: token.LSTMT},
				{Typ: token.SWITCH},
				{Typ: token.WS, Val: " "},
				{Typ: token.IDENT, Val: "name"},
				{Typ: token.RSTMT},
			},
		},
		{
			name:  "Simple case statement",
			input: "{%case value%}",
			expected: []token.Token{
				{Typ: token.LSTMT},
				{Typ: token.CASE},
				{Typ: token.WS, Val: " "},
				{Typ: token.IDENT, Val: "value"},
				{Typ: token.RSTMT},
			},
		},
		{
			name:  "Simple default statement",
			input: "{%default%}",
			expected: []token.Token{
				{Typ: token.LSTMT},
				{Typ: token.DEFAULT},
				{Typ: token.RSTMT},
			},
		},
		{
			name:  "Simple end statement",
			input: "{%end%}",
			expected: []token.Token{
				{Typ: token.LSTMT},
				{Typ: token.END},
				{Typ: token.RSTMT},
			},
		},
		{
			name:  "Simple genif statement",
			input: "{%genif name%}",
			expected: []token.Token{
				{Typ: token.LSTMT},
				{Typ: token.GENIF},
				{Typ: token.WS, Val: " "},
				{Typ: token.IDENT, Val: "name"},
				{Typ: token.RSTMT},
			},
		},
		{
			name:  "If with equal expression",
			input: "{%if name==3%}",
			expected: []token.Token{
				{Typ: token.LSTMT},
				{Typ: token.IF},
				{Typ: token.WS, Val: " "},
				{Typ: token.IDENT, Val: "name"},
				{Typ: token.EQL},
				{Typ: token.INT, Val: "3"},
				{Typ: token.RSTMT},
			},
		},
	}
	runTestCases(t, testCases)
}
