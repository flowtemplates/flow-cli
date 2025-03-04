package parser_test

import (
	"encoding/json"
	"slices"
	"testing"

	"github.com/templatesflow/cli/pkg/parser"
	"github.com/templatesflow/cli/pkg/token"
)

type testCase struct {
	name        string
	str         string
	input       []token.Token
	expected    parser.Expr
	errExpected bool
}

func runTestCases(t *testing.T, testCases []testCase) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			parser := parser.New(tc.input)
			got, err := parser.Parse()

			if (err != nil) != tc.errExpected {
				t.Errorf("Input: %q\nUnexpected error: %v", tc.str, err)
				return
			}

			a, _ := json.MarshalIndent(tc.expected, "", "  ")
			b, _ := json.MarshalIndent(got, "", "  ")
			if !slices.Equal(a, b) {
				t.Errorf("Input: %q\nAST mismatch.\nExpected:\n%s\nGot:\n%s", tc.str, a, b)
			}
		})
	}
}
