package lexer

import (
	"reflect"
	"testing"
)

func TestLexer(t *testing.T) {
	testCases := []struct {
		name           string
		input          string
		expectedTokens []Token
	}{
		{
			name:  "Empty input",
			input: "",
			expectedTokens: []Token{
				{typ: TokenEOF, val: ""},
			},
		},
		{
			name:  "Text with no delimiters",
			input: "Hello, world!",
			expectedTokens: []Token{
				{typ: TokenText, val: "Hello, world!"},
				{typ: TokenEOF, val: ""},
			},
		},
		{
			name:  "Simple expression",
			input: "Hello, {{name}}!",
			expectedTokens: []Token{
				{typ: TokenText, val: "Hello, "},
				{typ: TokenLeftExpr, val: "{{"},
				{typ: TokenSymbol, val: "name"},
				{typ: TokenRightExpr, val: "}}"},
				{typ: TokenText, val: "!"},
				{typ: TokenEOF, val: ""},
			},
		},
		{
			name:  "Whitespaces inside expr",
			input: "Hello, {{ name		}}!",
			expectedTokens: []Token{
				{typ: TokenText, val: "Hello, "},
				{typ: TokenLeftExpr, val: "{{"},
				{typ: TokenWhitespace, val: " "},
				{typ: TokenSymbol, val: "name"},
				{typ: TokenWhitespace, val: "\t\t"},
				{typ: TokenRightExpr, val: "}}"},
				{typ: TokenText, val: "!"},
				{typ: TokenEOF, val: ""},
			},
		},
		{
			name:  "Multiple expressions",
			input: "{{greeting}}, {{name}}!",
			expectedTokens: []Token{
				{typ: TokenLeftExpr, val: "{{"},
				{typ: TokenSymbol, val: "greeting"},
				{typ: TokenRightExpr, val: "}}"},
				{typ: TokenText, val: ", "},
				{typ: TokenLeftExpr, val: "{{"},
				{typ: TokenSymbol, val: "name"},
				{typ: TokenRightExpr, val: "}}"},
				{typ: TokenText, val: "!"},
				{typ: TokenEOF, val: ""},
			},
		},
		{
			name:  "Text before, after, and between expressions",
			input: "Hello, {{greeting}}, {{name}}! Welcome!",
			expectedTokens: []Token{
				{typ: TokenText, val: "Hello, "},
				{typ: TokenLeftExpr, val: "{{"},
				{typ: TokenSymbol, val: "greeting"},
				{typ: TokenRightExpr, val: "}}"},
				{typ: TokenText, val: ", "},
				{typ: TokenLeftExpr, val: "{{"},
				{typ: TokenSymbol, val: "name"},
				{typ: TokenRightExpr, val: "}}"},
				{typ: TokenText, val: "! Welcome!"},
				{typ: TokenEOF, val: ""},
			},
		},
		{
			name:  "Unclosed expression",
			input: "Hello, {{name",
			expectedTokens: []Token{
				{typ: TokenText, val: "Hello, "},
				{typ: TokenLeftExpr, val: "{{"},
				{typ: TokenSymbol, val: "name"},
				{typ: TokenError, val: "unexpected end of file inside action"},
			},
		},
		{
			name:  "Only left expr",
			input: "{{",
			expectedTokens: []Token{
				{typ: TokenLeftExpr, val: "{{"},
				{typ: TokenError, val: "unexpected end of file inside action"},
			},
		},
		{
			name:  "Only right expr",
			input: "}}",
			expectedTokens: []Token{
				{typ: TokenText, val: "}}"},
				{typ: TokenEOF, val: ""},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualTokens := Lex(tc.input)
			if !reflect.DeepEqual(actualTokens, tc.expectedTokens) {
				t.Errorf("Test Case: %s\nExpected:\n%v\nGot:\n%v",
					tc.name, tc.expectedTokens, actualTokens)
			}
		})
	}
}
