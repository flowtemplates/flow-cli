package lexer_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/flowtemplates/cli/pkg/flow-go/lexer"
	"github.com/flowtemplates/cli/pkg/flow-go/token"
)

func equal(gotTokens []token.Token, expectedTokens []token.Token) error {
	l := len(gotTokens)
	if l != len(expectedTokens) {
		return errors.New("not matching length")
	}

	for i := range l {
		got, expected := gotTokens[i], expectedTokens[i]

		if got.Typ != expected.Typ {
			return fmt.Errorf("wrong type: expected %s, got %s", got.Typ, expected.Typ)
		}

		var expectedValue string
		if expected.IsValueable() {
			expectedValue = expected.Val
		} else {
			expectedValue = token.TokenString(expected.Typ)
		}

		// if expected.Pos != got.Pos {
		// 	return fmt.Errorf("wrong pos: expected %+v, got %+v", expected.Pos, got.Pos)
		// }

		if got.Val != expectedValue {
			return fmt.Errorf("wrong value: expected %q, got %q", expectedValue, got.Val)
		}
	}

	return nil
}

type testCase struct {
	name           string
	input          string
	expectedTokens []token.Token
}

func runTestCases(t *testing.T, testCases []testCase) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			l := lexer.LexString(tc.input)
			var tokens []token.Token
			for {
				tok := l.NextToken()
				if tok.Typ == token.EOF {
					break
				}
				tokens = append(tokens, tok)
			}
			if err := equal(tokens, tc.expectedTokens); err != nil {
				t.Errorf("%s\nTest Case: %s\nExpected:\n%v\nGot:\n%v",
					err, tc.name, tc.expectedTokens, tokens)
			}
		})
	}
}
