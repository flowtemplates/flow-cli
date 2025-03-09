package renderer_test

import (
	"testing"

	"github.com/flowtemplates/cli/pkg/flow-go/parser"
	"github.com/flowtemplates/cli/pkg/flow-go/renderer"
)

type testCase struct {
	name        string
	str         string
	input       []parser.Node
	context     renderer.Scope
	expected    string
	errExpected bool
}

func runTestCases(t *testing.T, testCases []testCase) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := renderer.RenderAst(tc.input, tc.context)
			if (err != nil) != tc.errExpected {
				t.Errorf("Input: %q\nUnexpected error: %v", tc.expected, err)
				return
			}

			if got != tc.expected {
				t.Errorf("Input: %q\nMismatch.\nExpected:\n%s\nGot:\n%s", tc.str, tc.expected, got)
			}
		})
	}
}
