package lexer_test

import (
	"testing"

	"github.com/templatesflow/cli/internal/lexer"
)

func TestExpressions(t *testing.T) {
	testCases := []testCase{
		{
			name:  "Empty input",
			input: "",
			expectedTokens: []lexer.Token{
				{},
			},
		},
		{
			name:  "Text with no custom syntax",
			input: "Hello, world!",
			expectedTokens: []lexer.Token{
				{Typ: lexer.TokenText, Val: "Hello, world!", Pos: 0},
				{},
			},
		},
		{
			name:  "Empty expression",
			input: "{{}}",
			expectedTokens: []lexer.Token{
				{Typ: lexer.TokenLeftExpr, Val: "{{", Pos: 0},
				{Typ: lexer.TokenRightExpr, Val: "}}", Pos: 2},
				{},
			},
		}, {
			name:  "Simple expression",
			input: "Hello, {{name}}!",
			expectedTokens: []lexer.Token{
				{Typ: lexer.TokenText, Val: "Hello, ", Pos: 0},
				{Typ: lexer.TokenLeftExpr, Val: "{{", Pos: 0},
				{Typ: lexer.TokenIdentifier, Val: "name", Pos: 0},
				{Typ: lexer.TokenRightExpr, Val: "}}", Pos: 0},
				{Typ: lexer.TokenText, Val: "!", Pos: 0},
				{},
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
				{},
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
				{Typ: lexer.TokenText, Val: " "},
				{Typ: lexer.TokenLeftExpr, Val: "{{"},
				{Typ: lexer.TokenIdentifier, Val: "user_name"},
				{Typ: lexer.TokenRightExpr, Val: "}}"},
				{Typ: lexer.TokenText, Val: " "},
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
			name:  "Multiline input with several blocks",
			input: "Hello,\n {{greeting}}\r\n{{name}}{{surname}}",
			expectedTokens: []lexer.Token{
				{Typ: lexer.TokenText, Val: "Hello,\n "},
				{Typ: lexer.TokenLeftExpr, Val: "{{"},
				{Typ: lexer.TokenIdentifier, Val: "greeting"},
				{Typ: lexer.TokenRightExpr, Val: "}}"},
				{Typ: lexer.TokenText, Val: "\r\n"},
				{Typ: lexer.TokenLeftExpr, Val: "{{"},
				{Typ: lexer.TokenIdentifier, Val: "name"},
				{Typ: lexer.TokenRightExpr, Val: "}}"},
				{Typ: lexer.TokenLeftExpr, Val: "{{"},
				{Typ: lexer.TokenIdentifier, Val: "surname"},
				{Typ: lexer.TokenRightExpr, Val: "}}"},
				{Typ: lexer.TokenEOF, Val: ""},
			},
		},
	}

	runTestCases(t, testCases)
}

// func TestExpressionsEdgeCases(t *testing.T) {
// 	testCases := []testCase{
// 		{
// 			name:  "Line break inside expression",
// 			input: "{{greeting\n}}\r\n{{name}}",
// 			expectedTokens: []lexer.Token{
// 				{Typ: lexer.TokenLeftExpr, Val: "{{"},
// 				{Typ: lexer.TokenIdentifier, Val: "greeting"},
// 				{Typ: lexer.TokenError, Val: "line break"},
// 				{Typ: lexer.TokenEOF, Val: ""},
// 			},
// 		},
// 		{
// 			name:  "Unclosed expression",
// 			input: "Hello, {{name",
// 			expectedTokens: []lexer.Token{
// 				{Typ: lexer.TokenText, Val: "Hello, "},
// 				{Typ: lexer.TokenLeftExpr, Val: "{{"},
// 				{Typ: lexer.TokenIdentifier, Val: "name"},
// 				{Typ: lexer.TokenError, Val: "unexpected EOF"},
// 				{Typ: lexer.TokenEOF, Val: ""},
// 			},
// 		},
// 		{
// 			name:  "Only left expr",
// 			input: "{{",
// 			expectedTokens: []lexer.Token{
// 				{Typ: lexer.TokenLeftExpr, Val: "{{"},
// 				{Typ: lexer.TokenError, Val: "unexpected EOF"},
// 				{Typ: lexer.TokenEOF, Val: ""},
// 			},
// 		},
// 		{
// 			name:  "Only right expr",
// 			input: "}}",
// 			expectedTokens: []lexer.Token{
// 				{Typ: lexer.TokenText, Val: "}}"},
// 				{Typ: lexer.TokenEOF, Val: ""},
// 			},
// 		},
// 	}

// 	runTestCases(t, testCases)
// }

func TestNumLiterals(t *testing.T) {
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

// func TestNumLiteralsEdgeCases(t *testing.T) {
// 	// TODO: add tests for edge cases
// 	testCases := []testCase{
// 		{
// 			name:  "Integer value with unclosed expression",
// 			input: "{{10} text",
// 			expectedTokens: []lexer.Token{
// 				{Typ: lexer.TokenLeftExpr, Val: "{{"},
// 				{Typ: lexer.TokenInteger, Val: "10"},
// 				{Typ: lexer.TokenError, Val: "unclosed expression"},
// 				{Typ: lexer.TokenEOF, Val: ""},
// 			},
// 		},
// 		{
// 			name:  "Multiple dots in float value",
// 			input: "{{-12.3.2}}",
// 			expectedTokens: []lexer.Token{
// 				{Typ: lexer.TokenLeftExpr, Val: "{{"},
// 				{Typ: lexer.TokenFloat, Val: "-12.3"},
// 				{Typ: lexer.TokenError, Val: "unexpected point after float"},
// 				{Typ: lexer.TokenEOF, Val: ""},
// 			},
// 		},
// 	}

// 	runTestCases(t, testCases)
// }

func TestStringLiterals(t *testing.T) {
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
	}

	runTestCases(t, testCases)
}

// func TestStringLiteralsEdgeCases(t *testing.T) {
// 	// TODO: add tests for edge cases
// 	testCases := []testCase{
// 		{
// 			name:  "String literal with unclosed expression",
// 			input: `{{"double"`,
// 			expectedTokens: []lexer.Token{
// 				{Typ: lexer.TokenLeftExpr, Val: "{{"},
// 				{Typ: lexer.TokenString, Val: `"double"`},
// 				{Typ: lexer.TokenRightExpr, Val: "}}"},
// 				{Typ: lexer.TokenEOF, Val: ""},
// 			},
// 		},
// 		{
// 			name:  "String literal unclosed quotes",
// 			input: `{{"asd}}`,
// 			expectedTokens: []lexer.Token{
// 				{Typ: lexer.TokenLeftExpr, Val: "{{"},
// 				{Typ: lexer.TokenString, Val: `"123 falseasd"`},
// 				{Typ: lexer.TokenWhitespace, Val: " "},
// 				{Typ: lexer.TokenError, Val: "unclosed expression"},
// 				{Typ: lexer.TokenEOF, Val: ""},
// 			},
// 		},
// {
// 	name:  "Multiple strings",
// 	input: `{{"123 falseasd" "bsdbq12 )_ asd" }}`,
// 	expectedTokens: []lexer.Token{
// 		{Typ: lexer.TokenLeftExpr, Val: "{{"},
// 		{Typ: lexer.TokenString, Val: `"123 falseasd"`},
// 		{Typ: lexer.TokenWhitespace, Val: " "},
// 		{Typ: lexer.TokenError, Val: "unclosed expression"},
// 		{Typ: lexer.TokenEOF, Val: ""},
// 	},
// },
// 	}

// 	runTestCases(t, testCases)
// }

func TestBooleanLiterals(t *testing.T) {
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
	}

	runTestCases(t, testCases)
}

// func TestBooleanLiteralsEdgeCases(t *testing.T) {
// 	// TODO: add tests for edge cases
// 	testCases := []testCase{
// 		{
// 			name:  "Multiple values",
// 			input: "{{false  true }}",
// 			expectedTokens: []lexer.Token{
// 				{Typ: lexer.TokenLeftExpr, Val: "{{"},
// 				{Typ: lexer.TokenBoolean, Val: "false"},
// 				{Typ: lexer.TokenWhitespace, Val: "  "},
// 				{Typ: lexer.TokenError, Val: "unclosed expression"},
// 				{Typ: lexer.TokenEOF, Val: ""},
// 			},
// 		},
// 	}

// 	runTestCases(t, testCases)
// }

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
				{Typ: lexer.TokenFilterIdentifier, Val: "upper"},
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
				{Typ: lexer.TokenFilterIdentifier, Val: "upper"},
				{Typ: lexer.TokenWhitespace, Val: " "},
				{Typ: lexer.TokenPipe, Val: "->"},
				{Typ: lexer.TokenWhitespace, Val: " "},
				{Typ: lexer.TokenFilterIdentifier, Val: "camel"},
				{Typ: lexer.TokenWhitespace, Val: " "},
				{Typ: lexer.TokenRightExpr, Val: "}}"},
				{Typ: lexer.TokenEOF, Val: ""},
			},
		},
	}

	runTestCases(t, testCases)
}
