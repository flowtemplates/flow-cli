package parser_test

import (
	"encoding/json"
	"slices"
	"testing"

	"github.com/templatesflow/cli/internal/parser"
	"github.com/templatesflow/cli/internal/token"
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
				t.Errorf("Unexpected error: %v", err)
				return
			}

			a, _ := json.MarshalIndent(tc.expected, "", "  ")
			b, _ := json.MarshalIndent(got, "", "  ")
			if !slices.Equal(a, b) {
				t.Errorf("AST mismatch.\nExpected:\n%s\nGot:\n%s", a, b)
			}
		})
	}
}
