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
		{
			name:  "Simple switch statement",
			input: "{%switch name%}",
			expectedTokens: []lexer.Token{
				{Typ: lexer.TokenLeftStmt, Val: "{%"},
				{Typ: lexer.TokenSwitchStmt, Val: "switch"},
				{Typ: lexer.TokenWhitespace, Val: " "},
				{Typ: lexer.TokenIdentifier, Val: "name"},
				{Typ: lexer.TokenRightStmt, Val: "%}"},
				{},
			},
		},
		{
			name:  "Simple case statement",
			input: "{%case value%}",
			expectedTokens: []lexer.Token{
				{Typ: lexer.TokenLeftStmt, Val: "{%"},
				{Typ: lexer.TokenCaseStmt, Val: "case"},
				{Typ: lexer.TokenWhitespace, Val: " "},
				{Typ: lexer.TokenIdentifier, Val: "value"},
				{Typ: lexer.TokenRightStmt, Val: "%}"},
				{},
			},
		},
		{
			name:  "Simple default statement",
			input: "{%default%}",
			expectedTokens: []lexer.Token{
				{Typ: lexer.TokenLeftStmt, Val: "{%"},
				{Typ: lexer.TokenDefaultStmt, Val: "default"},
				{Typ: lexer.TokenRightStmt, Val: "%}"},
				{},
			},
		},
	}

	runTestCases(t, testCases)
}
