package lexer_test

import (
	"reflect"
	"testing"

	"github.com/templatesflow/cli/internal/lexer"
)

func TestLexer(t *testing.T) {
	testCases := []struct {
		name           string
		input          string
		expectedTokens []lexer.Token
	}{
		{
			name:  "Empty input",
			input: "",
			expectedTokens: []lexer.Token{
				{Typ: lexer.TokenEOF, Val: ""},
			},
		},
		{
			name:  "Text with no delimiters",
			input: "Hello, world!",
			expectedTokens: []lexer.Token{
				{Typ: lexer.TokenText, Val: "Hello, world!"},
				{Typ: lexer.TokenEOF, Val: ""},
			},
		},
		{
			name:  "Simple expression",
			input: "Hello, {{name}}!",
			expectedTokens: []lexer.Token{
				{Typ: lexer.TokenText, Val: "Hello, "},
				{Typ: lexer.TokenLeftExpr, Val: "{{"},
				{Typ: lexer.TokenSymbol, Val: "name"},
				{Typ: lexer.TokenRightExpr, Val: "}}"},
				{Typ: lexer.TokenText, Val: "!"},
				{Typ: lexer.TokenEOF, Val: ""},
			},
		},
		{
			name:  "Whitespaces inside expr",
			input: "Hello, {{ name		}}!",
			expectedTokens: []lexer.Token{
				{Typ: lexer.TokenText, Val: "Hello, "},
				{Typ: lexer.TokenLeftExpr, Val: "{{"},
				{Typ: lexer.TokenWhitespace, Val: " "},
				{Typ: lexer.TokenSymbol, Val: "name"},
				{Typ: lexer.TokenWhitespace, Val: "\t\t"},
				{Typ: lexer.TokenRightExpr, Val: "}}"},
				{Typ: lexer.TokenText, Val: "!"},
				{Typ: lexer.TokenEOF, Val: ""},
			},
		},
		{
			name:  "Multiple expressions",
			input: "{{greeting}}, {{name}}!",
			expectedTokens: []lexer.Token{
				{Typ: lexer.TokenLeftExpr, Val: "{{"},
				{Typ: lexer.TokenSymbol, Val: "greeting"},
				{Typ: lexer.TokenRightExpr, Val: "}}"},
				{Typ: lexer.TokenText, Val: ", "},
				{Typ: lexer.TokenLeftExpr, Val: "{{"},
				{Typ: lexer.TokenSymbol, Val: "name"},
				{Typ: lexer.TokenRightExpr, Val: "}}"},
				{Typ: lexer.TokenText, Val: "!"},
				{Typ: lexer.TokenEOF, Val: ""},
			},
		},
		{
			name:  "Var name underscores and digits",
			input: "{{GREETING}} {{user_name}} {{user123}}",
			expectedTokens: []lexer.Token{
				{Typ: lexer.TokenLeftExpr, Val: "{{"},
				{Typ: lexer.TokenSymbol, Val: "GREETING"},
				{Typ: lexer.TokenRightExpr, Val: "}}"},
				{Typ: lexer.TokenWhitespace, Val: " "},
				{Typ: lexer.TokenLeftExpr, Val: "{{"},
				{Typ: lexer.TokenSymbol, Val: "user_name"},
				{Typ: lexer.TokenRightExpr, Val: "}}"},
				{Typ: lexer.TokenWhitespace, Val: " "},
				{Typ: lexer.TokenLeftExpr, Val: "{{"},
				{Typ: lexer.TokenSymbol, Val: "user123"},
				{Typ: lexer.TokenRightExpr, Val: "}}"},
				{Typ: lexer.TokenEOF, Val: ""},
			},
		},
		{
			name:  "Text before, after, and between expressions",
			input: "Hello, {{greeting}}, {{name}}! Welcome!",
			expectedTokens: []lexer.Token{
				{Typ: lexer.TokenText, Val: "Hello, "},
				{Typ: lexer.TokenLeftExpr, Val: "{{"},
				{Typ: lexer.TokenSymbol, Val: "greeting"},
				{Typ: lexer.TokenRightExpr, Val: "}}"},
				{Typ: lexer.TokenText, Val: ", "},
				{Typ: lexer.TokenLeftExpr, Val: "{{"},
				{Typ: lexer.TokenSymbol, Val: "name"},
				{Typ: lexer.TokenRightExpr, Val: "}}"},
				{Typ: lexer.TokenText, Val: "! Welcome!"},
				{Typ: lexer.TokenEOF, Val: ""},
			},
		},
		{
			name:  "Unclosed expression",
			input: "Hello, {{name",
			expectedTokens: []lexer.Token{
				{Typ: lexer.TokenText, Val: "Hello, "},
				{Typ: lexer.TokenLeftExpr, Val: "{{"},
				{Typ: lexer.TokenSymbol, Val: "name"},
				{Typ: lexer.TokenError, Val: "unexpected EOF"},
				{Typ: lexer.TokenEOF, Val: ""},
			},
		},
		{
			name:  "Only left expr",
			input: "{{",
			expectedTokens: []lexer.Token{
				{Typ: lexer.TokenLeftExpr, Val: "{{"},
				{Typ: lexer.TokenError, Val: "unexpected EOF"},
				{Typ: lexer.TokenEOF, Val: ""},
			},
		},
		{
			name:  "Only right expr",
			input: "}}",
			expectedTokens: []lexer.Token{
				{Typ: lexer.TokenText, Val: "}}"},
				{Typ: lexer.TokenEOF, Val: ""},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			l := lexer.Lex(tc.input)
			var tokens []lexer.Token
			for {
				tok := l.NextToken()
				tokens = append(tokens, tok)
				if tok.Typ == lexer.TokenEOF {
					break
				}
			}
			if !reflect.DeepEqual(tokens, tc.expectedTokens) {
				t.Errorf("Test Case: %s\nExpected:\n%v\nGot:\n%v",
					tc.name, tc.expectedTokens, tokens)
			}
		})
	}
}
