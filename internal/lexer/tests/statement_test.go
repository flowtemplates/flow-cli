package lexer_test

import (
	"testing"

	"github.com/templatesflow/cli/internal/token"
)

func TestStatement(t *testing.T) {
	testCases := []testCase{
		{
			name:  "Simple if statement",
			input: "{%if name%}",
			expectedTokens: []token.Token{
				{Typ: token.LSTMT},
				{Typ: token.IF},
				{Typ: token.WHITESPACE, Val: " "},
				{Typ: token.IDENT, Val: "name"},
				{Typ: token.RSTMT},
			},
		},
		{
			name:  "Simple switch statement",
			input: "{%switch name%}",
			expectedTokens: []token.Token{
				{Typ: token.LSTMT},
				{Typ: token.SWITCH},
				{Typ: token.WHITESPACE, Val: " "},
				{Typ: token.IDENT, Val: "name"},
				{Typ: token.RSTMT},
			},
		},
		{
			name:  "Simple case statement",
			input: "{%case value%}",
			expectedTokens: []token.Token{
				{Typ: token.LSTMT},
				{Typ: token.CASE},
				{Typ: token.WHITESPACE, Val: " "},
				{Typ: token.IDENT, Val: "value"},
				{Typ: token.RSTMT},
			},
		},
		{
			name:  "Simple default statement",
			input: "{%default%}",
			expectedTokens: []token.Token{
				{Typ: token.LSTMT},
				{Typ: token.DEFAULT, Val: "default"},
				{Typ: token.RSTMT},
			},
		},
		{
			name:  "If with equal expression",
			input: "{%if name==3%}",
			expectedTokens: []token.Token{
				{Typ: token.LSTMT},
				{Typ: token.IF},
				{Typ: token.WHITESPACE, Val: " "},
				{Typ: token.IDENT, Val: "name"},
				{Typ: token.EQL},
				{Typ: token.INT, Val: "3"},
				{Typ: token.RSTMT},
			},
		},
	}
	runTestCases(t, testCases)
}
