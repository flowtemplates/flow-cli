package lexer_test

import (
	"testing"

	"github.com/templatesflow/cli/internal/lexer"
)

func TestComments(t *testing.T) {
	// TODO: add tests for edge cases
	testCases := []testCase{
		{
			name:  "Empty comment",
			input: "{##}",
			expectedTokens: []lexer.Token{
				{Typ: lexer.TokenLeftComm, Val: "{#"},
				{Typ: lexer.TokenRightComm, Val: "#}"},
				{Typ: lexer.TokenEOF, Val: ""},
			},
		},
		{
			name:  "Single comment",
			input: "{# no comments.. #}",
			expectedTokens: []lexer.Token{
				{Typ: lexer.TokenLeftComm, Val: "{#"},
				{Typ: lexer.TokenCommText, Val: ` no comments.. `},
				{Typ: lexer.TokenRightComm, Val: "#}"},
				{Typ: lexer.TokenEOF, Val: ""},
			},
		},
		{
			name:  "Multiline comment",
			input: "{# line 1\nline 2\r\n\nline 3 #}",
			expectedTokens: []lexer.Token{
				{Typ: lexer.TokenLeftComm, Val: "{#"},
				{Typ: lexer.TokenCommText, Val: " line 1\nline 2\r\n\nline 3 "},
				{Typ: lexer.TokenRightComm, Val: "#}"},
				{Typ: lexer.TokenEOF, Val: ""},
			},
		},
	}

	runTestCases(t, testCases)
}
