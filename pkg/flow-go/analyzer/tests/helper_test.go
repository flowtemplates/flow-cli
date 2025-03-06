package analyzer_test

import (
	"encoding/json"
	"slices"
	"testing"

	"github.com/flowtemplates/cli/pkg/flow-go/analyzer"
	"github.com/flowtemplates/cli/pkg/flow-go/parser"
)

type testCase struct {
	name        string
	str         string
	input       []parser.Node
	expected    analyzer.TypeMap
	errExpected bool
}

func runTestCases(t *testing.T, testCases []testCase) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, errs := analyzer.GetTypeMap(tc.input)
			if (len(errs) != 0) != tc.errExpected {
				t.Errorf("Input: %q\nUnexpected error: %v", tc.expected, errs)
				return
			}

			a, _ := json.MarshalIndent(tc.expected, "", "  ")
			b, _ := json.MarshalIndent(got, "", "  ")
			if !slices.Equal(a, b) {
				t.Errorf("Input: %q\nTypeMap mismatch.\nExpected:\n%s\nGot:\n%s", tc.str, a, b)
			}
		})
	}
}
