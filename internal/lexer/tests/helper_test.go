package lexer_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/templatesflow/cli/internal/lexer"
)

func equal(tokens1 []lexer.Token, tokens2 []lexer.Token) error {
	l := len(tokens1)
	if l != len(tokens2) {
		return errors.New("not matching length")
	}

	for i := range l {
		t1, t2 := tokens1[i], tokens2[i]

		if t1.Typ != t2.Typ {
			return fmt.Errorf("wrong type: expected %d, got %d", t1.Typ, t2.Typ)
		}

		if t1.Val != t2.Val {
			return fmt.Errorf("wrong value: expected %s, got %s", t1.Val, t2.Val)
		}
	}

	return nil
}

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
			if err := equal(tokens, tc.expectedTokens); err != nil {
				t.Errorf("%s\nTest Case: %s\nExpected:\n%v\nGot:\n%v",
					err, tc.name, tc.expectedTokens, tokens)
			}
		})
	}
}
