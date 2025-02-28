package lexer_test

import (
	"reflect"
	"testing"

	"github.com/templatesflow/cli/internal/lexer"
)

type testCase struct {
	name           string
	input          string
	expectedTokens []lexer.Token
}

func runTestCases(t *testing.T, testCases []testCase) {
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

func TestExpression(t *testing.T) {
	testCases := []testCase{
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
			name:  "Empty expression",
			input: "{{}}",
			expectedTokens: []lexer.Token{
				{Typ: lexer.TokenLeftExpr, Val: "{{"},
				{Typ: lexer.TokenRightExpr, Val: "}}"},
				{Typ: lexer.TokenEOF, Val: ""},
			},
		}, {
			name:  "Simple expression",
			input: "Hello, {{name}}!",
			expectedTokens: []lexer.Token{
				{Typ: lexer.TokenText, Val: "Hello, "},
				{Typ: lexer.TokenLeftExpr, Val: "{{"},
				{Typ: lexer.TokenIdentifier, Val: "name"},
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
				{Typ: lexer.TokenIdentifier, Val: "name"},
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
				{Typ: lexer.TokenIdentifier, Val: "greeting"},
				{Typ: lexer.TokenRightExpr, Val: "}}"},
				{Typ: lexer.TokenText, Val: ", "},
				{Typ: lexer.TokenLeftExpr, Val: "{{"},
				{Typ: lexer.TokenIdentifier, Val: "name"},
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
				{Typ: lexer.TokenIdentifier, Val: "GREETING"},
				{Typ: lexer.TokenRightExpr, Val: "}}"},
				{Typ: lexer.TokenWhitespace, Val: " "},
				{Typ: lexer.TokenLeftExpr, Val: "{{"},
				{Typ: lexer.TokenIdentifier, Val: "user_name"},
				{Typ: lexer.TokenRightExpr, Val: "}}"},
				{Typ: lexer.TokenWhitespace, Val: " "},
				{Typ: lexer.TokenLeftExpr, Val: "{{"},
				{Typ: lexer.TokenIdentifier, Val: "user123"},
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
				{Typ: lexer.TokenIdentifier, Val: "greeting"},
				{Typ: lexer.TokenRightExpr, Val: "}}"},
				{Typ: lexer.TokenText, Val: ", "},
				{Typ: lexer.TokenLeftExpr, Val: "{{"},
				{Typ: lexer.TokenIdentifier, Val: "name"},
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
				{Typ: lexer.TokenIdentifier, Val: "name"},
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

	runTestCases(t, testCases)
}

func TestNumLiterals(t *testing.T) {
	// TODO: add tests for edge cases
	testCases := []testCase{
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
	}

	runTestCases(t, testCases)
}

func TestStringLiterals(t *testing.T) {
	// TODO: add tests for edge cases
	testCases := []testCase{
		{
			name:  "Simple string literal in double quotes",
			input: `{{"double"}}`,
			expectedTokens: []lexer.Token{
				{Typ: lexer.TokenLeftExpr, Val: "{{"},
				{Typ: lexer.TokenString, Val: `"double"`},
				{Typ: lexer.TokenRightExpr, Val: "}}"},
				{Typ: lexer.TokenEOF, Val: ""},
			},
		},
		{
			name:  "Simple string literal in single quotes",
			input: `{{'single'}}`,
			expectedTokens: []lexer.Token{
				{Typ: lexer.TokenLeftExpr, Val: "{{"},
				{Typ: lexer.TokenString, Val: `'single'`},
				{Typ: lexer.TokenRightExpr, Val: "}}"},
				{Typ: lexer.TokenEOF, Val: ""},
			},
		},
		{
			name:  "Empty string literal",
			input: `{{ "" }}`,
			expectedTokens: []lexer.Token{
				{Typ: lexer.TokenLeftExpr, Val: "{{"},
				{Typ: lexer.TokenWhitespace, Val: " "},
				{Typ: lexer.TokenString, Val: `""`},
				{Typ: lexer.TokenWhitespace, Val: " "},
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
			name:  "Multiple strings",
			input: `{{"123 falseasd" "bsdbq12 )_ asd" }}`,
			expectedTokens: []lexer.Token{
				{Typ: lexer.TokenLeftExpr, Val: "{{"},
				{Typ: lexer.TokenString, Val: `"123 falseasd"`},
				{Typ: lexer.TokenWhitespace, Val: " "},
				{Typ: lexer.TokenError, Val: "unclosed expression"},
				{Typ: lexer.TokenEOF, Val: ""},
			},
		},
	}

	runTestCases(t, testCases)
}

func TestBooleanLiterals(t *testing.T) {
	// TODO: add tests for edge cases
	testCases := []testCase{
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
			name:  "Multiple values",
			input: "{{false  true }}",
			expectedTokens: []lexer.Token{
				{Typ: lexer.TokenLeftExpr, Val: "{{"},
				{Typ: lexer.TokenBoolean, Val: "false"},
				{Typ: lexer.TokenWhitespace, Val: "  "},
				{Typ: lexer.TokenError, Val: "unclosed expression"},
				{Typ: lexer.TokenEOF, Val: ""},
			},
		},
	}

	runTestCases(t, testCases)
}

func TestFilters(t *testing.T) {
	// TODO: add tests for edge cases
	testCases := []testCase{
		{
			name:  "Simple filter",
			input: "{{name->upper}}",
			expectedTokens: []lexer.Token{
				{Typ: lexer.TokenLeftExpr, Val: "{{"},
				{Typ: lexer.TokenIdentifier, Val: "name"},
				{Typ: lexer.TokenPipe, Val: "->"},
				{Typ: lexer.TokenFilter, Val: "upper"},
				{Typ: lexer.TokenRightExpr, Val: "}}"},
				{Typ: lexer.TokenEOF, Val: ""},
			},
		},
		{
			name:  "Nested filters",
			input: "{{ name -> upper -> camel }}",
			expectedTokens: []lexer.Token{
				{Typ: lexer.TokenLeftExpr, Val: "{{"},
				{Typ: lexer.TokenWhitespace, Val: " "},
				{Typ: lexer.TokenIdentifier, Val: "name"},
				{Typ: lexer.TokenWhitespace, Val: " "},
				{Typ: lexer.TokenPipe, Val: "->"},
				{Typ: lexer.TokenWhitespace, Val: " "},
				{Typ: lexer.TokenFilter, Val: "upper"},
				{Typ: lexer.TokenWhitespace, Val: " "},
				{Typ: lexer.TokenPipe, Val: "->"},
				{Typ: lexer.TokenWhitespace, Val: " "},
				{Typ: lexer.TokenFilter, Val: "camel"},
				{Typ: lexer.TokenWhitespace, Val: " "},
				{Typ: lexer.TokenRightExpr, Val: "}}"},
				{Typ: lexer.TokenEOF, Val: ""},
			},
		},
	}

	runTestCases(t, testCases)
}

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
	}

	runTestCases(t, testCases)
}
