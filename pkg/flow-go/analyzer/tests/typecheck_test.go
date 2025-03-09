package analyzer_test

import (
	"errors"
	"slices"
	"testing"

	"github.com/flowtemplates/cli/pkg/flow-go/analyzer"
	"github.com/flowtemplates/cli/pkg/flow-go/renderer"
	"github.com/flowtemplates/cli/pkg/flow-go/types"
)

func TestTypecheck(t *testing.T) {
	testCases := []struct {
		name        string
		scope       renderer.Scope
		tm          analyzer.TypeMap
		expectedErr []analyzer.TypeError
	}{
		{
			name:        "Empty inputs",
			scope:       renderer.Scope{},
			tm:          analyzer.TypeMap{},
			expectedErr: nil,
		},
		{
			name: "String",
			scope: renderer.Scope{
				"name": "some_str",
			},
			tm: analyzer.TypeMap{
				"name": types.String,
			},
			expectedErr: nil,
		},
		{
			name: "Number",
			scope: renderer.Scope{
				"age": "123",
			},
			tm: analyzer.TypeMap{
				"age": types.Number,
			},
			expectedErr: nil,
		},
		{
			name: "Boolean",
			scope: renderer.Scope{
				"falsy_flag":  "false",
				"truthy_flag": "true",
			},
			tm: analyzer.TypeMap{
				"falsy_flag":  types.Boolean,
				"truthy_flag": types.Boolean,
			},
			expectedErr: nil,
		},
		{
			name: "Boolean",
			scope: renderer.Scope{
				"falsy_flag":  "false",
				"truthy_flag": "true",
			},
			tm: analyzer.TypeMap{
				"falsy_flag":  types.Boolean,
				"truthy_flag": types.Boolean,
			},
			expectedErr: nil,
		},
		{
			name: "Boolean expected, got some string",
			scope: renderer.Scope{
				"flag": "asd",
			},
			tm: analyzer.TypeMap{
				"flag": types.Boolean,
			},
			expectedErr: []analyzer.TypeError{
				{
					ExpectedType: types.Boolean,
					Name:         "flag",
					Val:          "asd",
				},
			},
		},
		{
			name: "Number expected, got some string",
			scope: renderer.Scope{
				"num": "asd",
			},
			tm: analyzer.TypeMap{
				"num": types.Number,
			},
			expectedErr: []analyzer.TypeError{
				{
					ExpectedType: types.Number,
					Name:         "num",
					Val:          "asd",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := analyzer.Typecheck(tc.scope, tc.tm)
			var typeErrors analyzer.TypeErrors
			if errors.As(err, &typeErrors) {
				if !slices.Equal(typeErrors, tc.expectedErr) {
					t.Fatalf("expected %s, got %s", tc.expectedErr, err)
				}
			} else {
				t.Fatalf("unexpected error: %s", err)
			}
		})
	}
}
