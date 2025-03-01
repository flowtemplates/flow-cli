package lexer_test

import (
	"testing"

	"github.com/templatesflow/cli/internal/lexer"
)

func TestStatement(t *testing.T) {
	testCases := []testCase{
		{
			name:  "Simple if statement",
			input: "{%if name%}",
			expectedTokens: []lexer.Token{
				{Typ: lexer.TokenLeftStmt, Val: "{%"},
				{Typ: lexer.TokenIfStmt, Val: "if"},
				{Typ: lexer.TokenWhitespace, Val: " "},
				{Typ: lexer.TokenIdentifier, Val: "name"},
				{Typ: lexer.TokenRightStmt, Val: "%}"},
				{},
			},
		},
	}

	runTestCases(t, testCases)
}
