package lexer_test

import (
	"testing"

	"github.com/flowtemplates/cli/pkg/flow-go/token"
)

func TestExpressions(t *testing.T) {
	testCases := []testCase{
		{
			name:     "Empty input",
			input:    "",
			expected: []token.Token{},
		},
		{
			name:  "Plain text",
			input: "Hello, world!",
			expected: []token.Token{
				{Typ: token.TEXT, Val: "Hello, world!"},
			},
		},
		{
			name:  "Simple expression",
			input: "{{name}}",
			expected: []token.Token{
				{Typ: token.LEXPR},
				{Typ: token.IDENT, Val: "name"},
				{Typ: token.REXPR},
			},
		},
		{
			name:  "Whitespaces inside expr",
			input: "{{ name		}}",
			expected: []token.Token{
				{Typ: token.LEXPR},
				{Typ: token.WS, Val: " "},
				{Typ: token.IDENT, Val: "name"},
				{Typ: token.WS, Val: "		"},
				{Typ: token.REXPR},
			},
		},
		{
			name:  "Multiple expressions",
			input: "{{greeting}}, {{name}}",
			expected: []token.Token{
				{Typ: token.LEXPR},
				{Typ: token.IDENT, Val: "greeting"},
				{Typ: token.REXPR},
				{Typ: token.TEXT, Val: ", "},
				{Typ: token.LEXPR},
				{Typ: token.IDENT, Val: "name"},
				{Typ: token.REXPR},
			},
		},
		{
			name:  "Var name with underscores",
			input: "{{_user_name}}",
			expected: []token.Token{
				{Typ: token.LEXPR},
				{Typ: token.IDENT, Val: "_user_name"},
				{Typ: token.REXPR},
			},
		},
		{
			name:  "Var name with digits",
			input: "{{user123}}",
			expected: []token.Token{
				{Typ: token.LEXPR},
				{Typ: token.IDENT, Val: "user123"},
				{Typ: token.REXPR},
			},
		},
		{
			name:  "Var name with leading dollar sign",
			input: "{{$name}}",
			expected: []token.Token{
				{Typ: token.LEXPR},
				{Typ: token.IDENT, Val: "$name"},
				{Typ: token.REXPR},
			},
		},
		{
			name:  "Var name with dollar sign",
			input: "{{mirco$oft}}",
			expected: []token.Token{
				{Typ: token.LEXPR},
				{Typ: token.IDENT, Val: "mirco$oft"},
				{Typ: token.REXPR},
			},
		},
		{
			name:  "Var name with non-latin symbols",
			input: "{{こんにちは}} {{🙋}}",
			expected: []token.Token{
				{Typ: token.LEXPR},
				{Typ: token.IDENT, Val: "こんにちは"},
				{Typ: token.REXPR},
				{Typ: token.TEXT, Val: " "},
				{Typ: token.LEXPR},
				{Typ: token.IDENT, Val: "🙋"},
				{Typ: token.REXPR},
			},
		},
		{
			name:  "Text before, after, and between expressions",
			input: "Hello, {{greeting}}, {{name}}! Welcome!",
			expected: []token.Token{
				{Typ: token.TEXT, Val: "Hello, "},
				{Typ: token.LEXPR},
				{Typ: token.IDENT, Val: "greeting"},
				{Typ: token.REXPR},
				{Typ: token.TEXT, Val: ", "},
				{Typ: token.LEXPR},
				{Typ: token.IDENT, Val: "name"},
				{Typ: token.REXPR},
				{Typ: token.TEXT, Val: "! Welcome!"},
			},
		},
		{
			name:  "Multiline input with several blocks",
			input: "Hello,\n {{greeting}}\r\n{{name}}{{surname}}",
			expected: []token.Token{
				{Typ: token.TEXT, Val: "Hello,\n "},
				{Typ: token.LEXPR},
				{Typ: token.IDENT, Val: "greeting"},
				{Typ: token.REXPR},
				{Typ: token.TEXT, Val: "\r\n"},
				{Typ: token.LEXPR},
				{Typ: token.IDENT, Val: "name"},
				{Typ: token.REXPR},
				{Typ: token.LEXPR},
				{Typ: token.IDENT, Val: "surname"},
				{Typ: token.REXPR},
			},
		},
	}
	runTestCases(t, testCases)
}

func TestExpressionsEdgeCases(t *testing.T) {
	testCases := []testCase{
		{
			name:  "Empty expression",
			input: "{{}}",
			expected: []token.Token{
				{Typ: token.LEXPR},
				{Typ: token.REXPR},
			},
		},
		{
			name:  "Line break inside expression",
			input: "{{greeting\n}}",
			expected: []token.Token{
				{Typ: token.LEXPR},
				{Typ: token.IDENT, Val: "greeting"},
				{Typ: token.TEXT, Val: "\n"},
				{Typ: token.REXPR},
			},
		},
		{
			name:  "Unclosed expression",
			input: "Hello, {{name",
			expected: []token.Token{
				{Typ: token.TEXT, Val: "Hello, "},
				{Typ: token.LEXPR},
				{Typ: token.IDENT, Val: "name"},
			},
		},
		{
			name:  "Only left expr",
			input: "{{",
			expected: []token.Token{
				{Typ: token.LEXPR},
			},
		},
		{
			name:  "Only right expr",
			input: "}}",
			expected: []token.Token{
				{Typ: token.REXPR},
			},
		},
		{
			name:  "Text after unclosed expression",
			input: "{{greeting\r\nSome text",
			expected: []token.Token{
				{Typ: token.LEXPR},
				{Typ: token.IDENT, Val: "greeting"},
				{Typ: token.TEXT, Val: "\r\nSome text"},
			},
		},
		{
			name:  "Single right expr between text",
			input: "Another text}}Some text",
			expected: []token.Token{
				{Typ: token.TEXT, Val: "Another text"},
				{Typ: token.REXPR},
				{Typ: token.TEXT, Val: "Some text"},
			},
		},
		{
			name:  "Right expr after valid expression",
			input: "{{name}}\n}}Some text",
			expected: []token.Token{
				{Typ: token.LEXPR},
				{Typ: token.IDENT, Val: "name"},
				{Typ: token.REXPR},
				{Typ: token.TEXT, Val: "\n"},
				{Typ: token.REXPR},
				{Typ: token.TEXT, Val: "Some text"},
			},
		},
	}
	runTestCases(t, testCases)
}

func TestNumLiterals(t *testing.T) {
	testCases := []testCase{
		{
			name:  "Integer value",
			input: "{{10}}",
			expected: []token.Token{
				{Typ: token.LEXPR},
				{Typ: token.INT, Val: "10"},
				{Typ: token.REXPR},
			},
		},
		{
			name:  "Negative integer value",
			input: "{{-123}}",
			expected: []token.Token{
				{Typ: token.LEXPR},
				{Typ: token.SUB, Val: "-"},
				{Typ: token.INT, Val: "123"},
				{Typ: token.REXPR},
			},
		},
		{
			name:  "Float value",
			input: "{{12.3}}",
			expected: []token.Token{
				{Typ: token.LEXPR},
				{Typ: token.FLOAT, Val: "12.3"},
				{Typ: token.REXPR},
			},
		},
		{
			name:  "Negative float value",
			input: "{{-12.3}}",
			expected: []token.Token{
				{Typ: token.LEXPR},
				{Typ: token.SUB, Val: "-"},
				{Typ: token.FLOAT, Val: "12.3"},
				{Typ: token.REXPR},
			},
		},
	}
	runTestCases(t, testCases)
}

func TestOperations(t *testing.T) {
	testCases := []testCase{
		{
			name:  "Addittion",
			input: "{{seconds+1}}",
			expected: []token.Token{
				{Typ: token.LEXPR},
				{Typ: token.IDENT, Val: "seconds"},
				{Typ: token.ADD},
				{Typ: token.INT, Val: "1"},
				{Typ: token.REXPR},
			},
		},
		{
			name:  "Subtraction",
			input: "{{age-123.2}}",
			expected: []token.Token{
				{Typ: token.LEXPR},
				{Typ: token.IDENT, Val: "age"},
				{Typ: token.SUB},
				{Typ: token.FLOAT, Val: "123.2"},
				{Typ: token.REXPR},
			},
		},
		{
			name:  "Negative number subtraction",
			input: "{{age- -123.2}}",
			expected: []token.Token{
				{Typ: token.LEXPR},
				{Typ: token.IDENT, Val: "age"},
				{Typ: token.SUB},
				{Typ: token.WS, Val: " "},
				{Typ: token.SUB},
				{Typ: token.FLOAT, Val: "123.2"},
				{Typ: token.REXPR},
			},
		},
		{
			name:  "Multiply",
			input: "{{age*30}}",
			expected: []token.Token{
				{Typ: token.LEXPR},
				{Typ: token.IDENT, Val: "age"},
				{Typ: token.MUL},
				{Typ: token.INT, Val: "30"},
				{Typ: token.REXPR},
			},
		},
		{
			name:  "Multiply by negative number",
			input: "{{age*-30}}",
			expected: []token.Token{
				{Typ: token.LEXPR},
				{Typ: token.IDENT, Val: "age"},
				{Typ: token.MUL},
				{Typ: token.SUB},
				{Typ: token.INT, Val: "30"},
				{Typ: token.REXPR},
			},
		},
		{
			name:  "Multiply by negative number",
			input: "{{age*-30}}",
			expected: []token.Token{
				{Typ: token.LEXPR},
				{Typ: token.IDENT, Val: "age"},
				{Typ: token.MUL},
				{Typ: token.SUB},
				{Typ: token.INT, Val: "30"},
				{Typ: token.REXPR},
			},
		},
		{
			name:  "Division",
			input: "{{age/30}}",
			expected: []token.Token{
				{Typ: token.LEXPR},
				{Typ: token.IDENT, Val: "age"},
				{Typ: token.DIV},
				{Typ: token.INT, Val: "30"},
				{Typ: token.REXPR},
			},
		},
		{
			name:  "Division by negative number",
			input: "{{age/-30}}",
			expected: []token.Token{
				{Typ: token.LEXPR},
				{Typ: token.IDENT, Val: "age"},
				{Typ: token.DIV},
				{Typ: token.SUB},
				{Typ: token.INT, Val: "30"},
				{Typ: token.REXPR},
			},
		},
		{
			name:  "Single parens",
			input: "{{(12/2)+age}}",
			expected: []token.Token{
				{Typ: token.LEXPR},
				{Typ: token.LPAREN},
				{Typ: token.INT, Val: "12"},
				{Typ: token.DIV},
				{Typ: token.INT, Val: "2"},
				{Typ: token.RPAREN},
				{Typ: token.ADD},
				{Typ: token.IDENT, Val: "age"},
				{Typ: token.REXPR},
			},
		},
		{
			name:  "Two operations with parens",
			input: "{{(age/-30)+(12-2.2)}}",
			expected: []token.Token{
				{Typ: token.LEXPR},
				{Typ: token.LPAREN},
				{Typ: token.IDENT, Val: "age"},
				{Typ: token.DIV},
				{Typ: token.SUB},
				{Typ: token.INT, Val: "30"},
				{Typ: token.RPAREN},
				{Typ: token.ADD},
				{Typ: token.LPAREN},
				{Typ: token.INT, Val: "12"},
				{Typ: token.SUB},
				{Typ: token.FLOAT, Val: "2.2"},
				{Typ: token.RPAREN},
				{Typ: token.REXPR},
			},
		},
	}
	runTestCases(t, testCases)
}

func TestOperationsEdgeCases(t *testing.T) {
	testCases := []testCase{
		{
			name:  "Unclosed addition",
			input: "{{1+}}",
			expected: []token.Token{
				{Typ: token.LEXPR},
				{Typ: token.INT, Val: "1"},
				{Typ: token.ADD},
				{Typ: token.REXPR},
			},
		},
		{
			name:  "Unclosed expression with addition",
			input: "{{1+",
			expected: []token.Token{
				{Typ: token.LEXPR},
				{Typ: token.INT, Val: "1"},
				{Typ: token.ADD},
			},
		},
	}
	runTestCases(t, testCases)
}

func TestNumLiteralsEdgeCases(t *testing.T) {
	testCases := []testCase{
		{
			name:  "Integer value with unclosed expression",
			input: "{{10} name",
			expected: []token.Token{
				{Typ: token.LEXPR},
				{Typ: token.INT, Val: "10"},
				{Typ: token.EXPECTED_EXPR, Val: "}"},
				{Typ: token.WS, Val: " "},
				{Typ: token.IDENT, Val: "name"},
			},
		},
		{
			name:  "Multiple points in float value",
			input: "{{-12.3.2}}",
			expected: []token.Token{
				{Typ: token.LEXPR},
				{Typ: token.SUB, Val: "-"},
				{Typ: token.FLOAT, Val: "12.3"},
				{Typ: token.EXPECTED_EXPR, Val: "."},
				{Typ: token.INT, Val: "2"},
				{Typ: token.REXPR},
			},
		},
	}
	runTestCases(t, testCases)
}

func TestStringLiterals(t *testing.T) {
	testCases := []testCase{
		{
			name:  "Simple string literal in double quotes",
			input: `{{"double"}}`,
			expected: []token.Token{
				{Typ: token.LEXPR},
				{Typ: token.STR, Val: `"double"`},
				{Typ: token.REXPR},
			},
		},
		{
			name:  "Simple string literal in single quotes",
			input: `{{'single'}}`,
			expected: []token.Token{
				{Typ: token.LEXPR},
				{Typ: token.STR, Val: `'single'`},
				{Typ: token.REXPR},
			},
		},
		{
			name:  "Empty string literal",
			input: `{{""}}`,
			expected: []token.Token{
				{Typ: token.LEXPR},
				{Typ: token.STR, Val: `""`},
				{Typ: token.REXPR},
			},
		},
		{
			name:  "String literal with whitespaces",
			input: `{{"word1 word2  	word3"}}`,
			expected: []token.Token{
				{Typ: token.LEXPR},
				{Typ: token.STR, Val: `"word1 word2  	word3"`},
				{Typ: token.REXPR},
			},
		},
		{
			name:  "String literal with numbers and booleans",
			input: `{{"123 false -22.0"}}`,
			expected: []token.Token{
				{Typ: token.LEXPR},
				{Typ: token.STR, Val: `"123 false -22.0"`},
				{Typ: token.REXPR},
			},
		},
	}
	runTestCases(t, testCases)
}

func TestStringLiteralsEdgeCases(t *testing.T) {
	testCases := []testCase{
		{
			name:  "String not terminated",
			input: `{{"double`,
			expected: []token.Token{
				{Typ: token.LEXPR},
				{Typ: token.NOT_TERMINATED_STR, Val: `"double`},
			},
		},
		{
			name:  "Multiple strings",
			input: `{{"123 falseasd" 'bsdbq12 )_ asd' }}`,
			expected: []token.Token{
				{Typ: token.LEXPR},
				{Typ: token.STR, Val: `"123 falseasd"`},
				{Typ: token.WS, Val: " "},
				{Typ: token.STR, Val: `'bsdbq12 )_ asd'`},
				{Typ: token.WS, Val: " "},
				{Typ: token.REXPR},
			},
		},
	}
	runTestCases(t, testCases)
}

func TestFilters(t *testing.T) {
	testCases := []testCase{
		{
			name:  "Simple filter",
			input: "{{name->upper}}",
			expected: []token.Token{
				{Typ: token.LEXPR},
				{Typ: token.IDENT, Val: "name"},
				{Typ: token.RARR},
				{Typ: token.IDENT, Val: "upper"},
				{Typ: token.REXPR},
			},
		},
		{
			name:  "Nested filters",
			input: "{{name -> upper -> camel}}",
			expected: []token.Token{
				{Typ: token.LEXPR},
				{Typ: token.IDENT, Val: "name"},
				{Typ: token.WS, Val: " "},
				{Typ: token.RARR},
				{Typ: token.WS, Val: " "},
				{Typ: token.IDENT, Val: "upper"},
				{Typ: token.WS, Val: " "},
				{Typ: token.RARR},
				{Typ: token.WS, Val: " "},
				{Typ: token.IDENT, Val: "camel"},
				{Typ: token.REXPR},
			},
		},
	}
	runTestCases(t, testCases)
}

func TestFiltersEdgeCases(t *testing.T) {
	testCases := []testCase{
		{
			name:  "Empty filter",
			input: "{{name->}}",
			expected: []token.Token{
				{Typ: token.LEXPR},
				{Typ: token.IDENT, Val: "name"},
				{Typ: token.RARR},
				{Typ: token.REXPR},
			},
		},
		{
			name:  "Filter in expression",
			input: "{{name->upper=='UP'}}",
			expected: []token.Token{
				{Typ: token.LEXPR},
				{Typ: token.IDENT, Val: "name"},
				{Typ: token.RARR},
				{Typ: token.IDENT, Val: "upper"},
				{Typ: token.EQL},
				{Typ: token.STR, Val: "'UP'"},
				{Typ: token.REXPR},
			},
		},
	}
	runTestCases(t, testCases)
}

func TestOperators(t *testing.T) {
	testCases := []testCase{
		{
			name:  "Equal",
			input: "{{age==3}}",
			expected: []token.Token{
				{Typ: token.LEXPR},
				{Typ: token.IDENT, Val: "age"},
				{Typ: token.EQL},
				{Typ: token.INT, Val: "3"},
				{Typ: token.REXPR},
			},
		},
		{
			name:  "Not",
			input: "{{not flag}}",
			expected: []token.Token{
				{Typ: token.LEXPR},
				{Typ: token.NOT},
				{Typ: token.WS, Val: " "},
				{Typ: token.IDENT, Val: "flag"},
				{Typ: token.REXPR},
			},
		},
		{
			name:  "Not equal",
			input: "{{age!=3}}",
			expected: []token.Token{
				{Typ: token.LEXPR},
				{Typ: token.IDENT, Val: "age"},
				{Typ: token.NEQL},
				{Typ: token.INT, Val: "3"},
				{Typ: token.REXPR},
			},
		},
		{
			name:  "Is",
			input: "{{age is 3}}",
			expected: []token.Token{
				{Typ: token.LEXPR},
				{Typ: token.IDENT, Val: "age"},
				{Typ: token.WS, Val: " "},
				{Typ: token.IS},
				{Typ: token.WS, Val: " "},
				{Typ: token.INT, Val: "3"},
				{Typ: token.REXPR},
			},
		},
		{
			name:  "is not",
			input: "{{age is not 3}}",
			expected: []token.Token{
				{Typ: token.LEXPR},
				{Typ: token.IDENT, Val: "age"},
				{Typ: token.WS, Val: " "},
				{Typ: token.IS},
				{Typ: token.WS, Val: " "},
				{Typ: token.NOT},
				{Typ: token.WS, Val: " "},
				{Typ: token.INT, Val: "3"},
				{Typ: token.REXPR},
			},
		},
		{
			name:  "Excl",
			input: "{{!flag}}",
			expected: []token.Token{
				{Typ: token.LEXPR},
				{Typ: token.EXCL},
				{Typ: token.IDENT, Val: "flag"},
				{Typ: token.REXPR},
			},
		},
		{
			name:  "Greater",
			input: "{{var>1}}",
			expected: []token.Token{
				{Typ: token.LEXPR},
				{Typ: token.IDENT, Val: "var"},
				{Typ: token.GTR},
				{Typ: token.INT, Val: "1"},
				{Typ: token.REXPR},
			},
		},
		{
			name:  "Less",
			input: "{{var<1}}",
			expected: []token.Token{
				{Typ: token.LEXPR},
				{Typ: token.IDENT, Val: "var"},
				{Typ: token.LESS},
				{Typ: token.INT, Val: "1"},
				{Typ: token.REXPR},
			},
		},
		{
			name:  "Less or equal",
			input: "{{var<=1}}",
			expected: []token.Token{
				{Typ: token.LEXPR},
				{Typ: token.IDENT, Val: "var"},
				{Typ: token.LEQ},
				{Typ: token.INT, Val: "1"},
				{Typ: token.REXPR},
			},
		},
		{
			name:  "Greater or equal",
			input: "{{var>=1}}",
			expected: []token.Token{
				{Typ: token.LEXPR},
				{Typ: token.IDENT, Val: "var"},
				{Typ: token.GEQ},
				{Typ: token.INT, Val: "1"},
				{Typ: token.REXPR},
			},
		},
		{
			name:  "And",
			input: "{{var and 1}}",
			expected: []token.Token{
				{Typ: token.LEXPR},
				{Typ: token.IDENT, Val: "var"},
				{Typ: token.WS, Val: " "},
				{Typ: token.AND},
				{Typ: token.WS, Val: " "},
				{Typ: token.INT, Val: "1"},
				{Typ: token.REXPR},
			},
		},
		{
			name:  "or",
			input: "{{var or 1}}",
			expected: []token.Token{
				{Typ: token.LEXPR},
				{Typ: token.IDENT, Val: "var"},
				{Typ: token.WS, Val: " "},
				{Typ: token.OR},
				{Typ: token.WS, Val: " "},
				{Typ: token.INT, Val: "1"},
				{Typ: token.REXPR},
			},
		},
		{
			name:  "&&",
			input: "{{var&&1}}",
			expected: []token.Token{
				{Typ: token.LEXPR},
				{Typ: token.IDENT, Val: "var"},
				{Typ: token.LAND},
				{Typ: token.INT, Val: "1"},
				{Typ: token.REXPR},
			},
		},
		{
			name:  "||",
			input: "{{var||1}}",
			expected: []token.Token{
				{Typ: token.LEXPR},
				{Typ: token.IDENT, Val: "var"},
				{Typ: token.LOR},
				{Typ: token.INT, Val: "1"},
				{Typ: token.REXPR},
			},
		},
	}
	runTestCases(t, testCases)
}
