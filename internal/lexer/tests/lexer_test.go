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
			name:  "Text with no custom syntax",
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
		{
			name:  "Integer value",
			input: "{{10}}",
			expectedTokens: []lexer.Token{
				{Typ: lexer.TokenLeftExpr, Val: "{{"},
				{Typ: lexer.TokenInteger, Val: "10"},
				{Typ: lexer.TokenRightExpr, Val: "}}"},
				{Typ: lexer.TokenEOF, Val: ""},
			},
		},
		{
			name:  "Integer value with unclosed expression",
			input: "{{10} text",
			expectedTokens: []lexer.Token{
				{Typ: lexer.TokenLeftExpr, Val: "{{"},
				{Typ: lexer.TokenInteger, Val: "10"},
				{Typ: lexer.TokenError, Val: "unclosed expression"},
				{Typ: lexer.TokenEOF, Val: ""},
			},
		},
		{
			name:  "Negative integer value",
			input: "{{-123}}",
			expectedTokens: []lexer.Token{
				{Typ: lexer.TokenLeftExpr, Val: "{{"},
				{Typ: lexer.TokenInteger, Val: "-123"},
				{Typ: lexer.TokenRightExpr, Val: "}}"},
				{Typ: lexer.TokenEOF, Val: ""},
			},
		},
		{
			name:  "True value",
			input: "{{true}}",
			expectedTokens: []lexer.Token{
				{Typ: lexer.TokenLeftExpr, Val: "{{"},
				{Typ: lexer.TokenBoolean, Val: "true"},
				{Typ: lexer.TokenRightExpr, Val: "}}"},
				{Typ: lexer.TokenEOF, Val: ""},
			},
		},
		{
			name:  "False value",
			input: "{{false}}",
			expectedTokens: []lexer.Token{
				{Typ: lexer.TokenLeftExpr, Val: "{{"},
				{Typ: lexer.TokenBoolean, Val: "false"},
				{Typ: lexer.TokenRightExpr, Val: "}}"},
				{Typ: lexer.TokenEOF, Val: ""},
			},
		},
		{
			name:  "Simple string literal",
			input: `{{"some_string"}}`,
			expectedTokens: []lexer.Token{
				{Typ: lexer.TokenLeftExpr, Val: "{{"},
				{Typ: lexer.TokenString, Val: `"some_string"`},
				{Typ: lexer.TokenRightExpr, Val: "}}"},
				{Typ: lexer.TokenEOF, Val: ""},
			},
		},
		{
			name:  "String literal with whitespaces",
			input: `{{"word1 word2  	word3"}}`,
			expectedTokens: []lexer.Token{
				{Typ: lexer.TokenLeftExpr, Val: "{{"},
				{Typ: lexer.TokenString, Val: `"word1 word2  	word3"`},
				{Typ: lexer.TokenRightExpr, Val: "}}"},
				{Typ: lexer.TokenEOF, Val: ""},
			},
		},
		{
			name:  "String literal with numbers and booleans",
			input: `{{"123 false -22.0"}}`,
			expectedTokens: []lexer.Token{
				{Typ: lexer.TokenLeftExpr, Val: "{{"},
				{Typ: lexer.TokenString, Val: `"123 false -22.0"`},
				{Typ: lexer.TokenRightExpr, Val: "}}"},
				{Typ: lexer.TokenEOF, Val: ""},
			},
		},
		{
			name:  "float value",
			input: "{{12.3}}",
			expectedTokens: []lexer.Token{
				{Typ: lexer.TokenLeftExpr, Val: "{{"},
				{Typ: lexer.TokenFloat, Val: "12.3"},
				{Typ: lexer.TokenRightExpr, Val: "}}"},
				{Typ: lexer.TokenEOF, Val: ""},
			},
		},
		{
			name:  "Negative float value",
			input: "{{-12.3}}",
			expectedTokens: []lexer.Token{
				{Typ: lexer.TokenLeftExpr, Val: "{{"},
				{Typ: lexer.TokenFloat, Val: "-12.3"},
				{Typ: lexer.TokenRightExpr, Val: "}}"},
				{Typ: lexer.TokenEOF, Val: ""},
			},
		},
		// {
		// 	name:  "Function",
		// 	input: "{{ name -> upper }}",
		// 	expectedTokens: []lexer.Token{
		// 		{Typ: lexer.TokenLeftExpr, Val: "{{"},
		// 		{Typ: lexer.TokenWhitespace, Val: " "},
		// 		{Typ: lexer.TokenSymbol, Val: "name"},
		// 		{Typ: lexer.TokenRightExpr, Val: "}}"},
		// 		{Typ: lexer.TokenEOF, Val: ""},
		// 	},
		// },
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
