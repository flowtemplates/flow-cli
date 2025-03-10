package lexer_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/flowtemplates/cli/pkg/flow-go/lexer"
	"github.com/flowtemplates/cli/pkg/flow-go/token"
)

func TestPosition(t *testing.T) {
	testCases := []testCase{
		{
			name:  "Empty if block",
			input: "{%if name%}\n{%end%}",
			expected: []token.Token{
				{Typ: token.LSTMT, Pos: token.Position{
					Line:   1,
					Column: 1,
					Offset: 0,
				}},
				{Typ: token.IF, Pos: token.Position{
					Line:   1,
					Column: 3,
					Offset: 2,
				}},
				{Typ: token.WS, Val: " ", Pos: token.Position{
					Line:   1,
					Column: 5,
					Offset: 4,
				}},
				{Typ: token.IDENT, Val: "name", Pos: token.Position{
					Line:   1,
					Column: 6,
					Offset: 5,
				}},
				{Typ: token.RSTMT, Pos: token.Position{
					Line:   1,
					Column: 10,
					Offset: 9,
				}},
				{Typ: token.TEXT, Val: "\n", Pos: token.Position{
					Line:   1,
					Column: 12,
					Offset: 11,
				}},
				{Typ: token.LSTMT, Pos: token.Position{
					Line:   2,
					Column: 1,
					Offset: 12,
				}},
				{Typ: token.END, Pos: token.Position{
					Line:   2,
					Column: 3,
					Offset: 14,
				}},
				{Typ: token.RSTMT, Pos: token.Position{
					Line:   2,
					Column: 6,
					Offset: 17,
				}},
			},
		},
		{
			name:  "Multiple expressions on many lines",
			input: "Hello {{name}}!\nFrom {{ flow }}templates\n\n{{1 + 3}}",
			expected: []token.Token{
				{Typ: token.TEXT, Val: "Hello ", Pos: token.Position{
					Line:   1,
					Column: 1,
					Offset: 0,
				}},
				{Typ: token.LEXPR, Pos: token.Position{
					Line:   1,
					Column: 7,
					Offset: 6,
				}},
				{Typ: token.IDENT, Val: "name", Pos: token.Position{
					Line:   1,
					Column: 9,
					Offset: 8,
				}},
				{Typ: token.REXPR, Pos: token.Position{
					Line:   1,
					Column: 13,
					Offset: 12,
				}},
				{Typ: token.TEXT, Val: "!\nFrom ", Pos: token.Position{
					Line:   1,
					Column: 15,
					Offset: 14,
				}},
				{Typ: token.LEXPR, Pos: token.Position{
					Line:   2,
					Column: 6,
					Offset: 21,
				}},
				{Typ: token.WS, Pos: token.Position{
					Line:   2,
					Column: 8,
					Offset: 23,
				}},
				{Typ: token.IDENT, Val: "flow", Pos: token.Position{
					Line:   2,
					Column: 9,
					Offset: 24,
				}},
				{Typ: token.WS, Pos: token.Position{
					Line:   2,
					Column: 13,
					Offset: 28,
				}},
				{Typ: token.REXPR, Pos: token.Position{
					Line:   2,
					Column: 14,
					Offset: 29,
				}},
				{Typ: token.TEXT, Val: "templates\n\n", Pos: token.Position{
					Line:   2,
					Column: 16,
					Offset: 31,
				}},
				{Typ: token.LEXPR, Pos: token.Position{
					Line:   4,
					Column: 1,
					Offset: 42,
				}},
				{Typ: token.INT, Val: "1", Pos: token.Position{
					Line:   4,
					Column: 3,
					Offset: 44,
				}},
				{Typ: token.WS, Pos: token.Position{
					Line:   4,
					Column: 4,
					Offset: 45,
				}},
				{Typ: token.ADD, Pos: token.Position{
					Line:   4,
					Column: 5,
					Offset: 46,
				}},
				{Typ: token.WS, Pos: token.Position{
					Line:   4,
					Column: 6,
					Offset: 47,
				}},
				{Typ: token.INT, Val: "1", Pos: token.Position{
					Line:   4,
					Column: 7,
					Offset: 48,
				}},
				{Typ: token.REXPR, Pos: token.Position{
					Line:   4,
					Column: 8,
					Offset: 49,
				}},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			l := lexer.FromString(tc.input)
			var tokens []token.Token
			for {
				tok := l.NextToken()
				if tok.Typ == token.EOF {
					break
				}
				tokens = append(tokens, tok)
			}

			eq := func(gotTokens []token.Token, expectedTokens []token.Token) error {
				l := len(gotTokens)
				if l != len(expectedTokens) {
					return errors.New("not matching length")
				}

				for i := range l {
					got, expected := gotTokens[i], expectedTokens[i]

					if expected.Pos != got.Pos {
						return fmt.Errorf("wrong pos: expected %+v, got %+v", expected.Pos, got.Pos)
					}
				}

				return nil
			}

			if err := eq(tokens, tc.expected); err != nil {
				t.Errorf("%s\nTest Case: %s\nExpected:\n%v\nGot:\n%v",
					err, tc.name, tc.expected, tokens)
			}
		})
	}
}
