package lexer_test

import (
	"testing"

	"github.com/templatesflow/cli/pkg/token"
)

func TestComments(t *testing.T) {
	testCases := []testCase{
		{
			name:  "Empty comment",
			input: "{##}",
			expectedTokens: []token.Token{
				{Typ: token.LCOMM},
				{Typ: token.RCOMM},
			},
		},
		{
			name:  "Single comment",
			input: "{# no comments.. #}",
			expectedTokens: []token.Token{
				{Typ: token.LCOMM},
				{Typ: token.COMM_TEXT, Val: ` no comments.. `},
				{Typ: token.RCOMM},
			},
		},
		{
			name:  "Multiline comment",
			input: "{# line 1\nline 2\r\n\nline 3 #}",
			expectedTokens: []token.Token{
				{Typ: token.LCOMM},
				{Typ: token.COMM_TEXT, Val: " line 1\nline 2\r\n\nline 3 "},
				{Typ: token.RCOMM},
			},
		},
		{
			name:  "Multiple comments",
			input: "{# line 1 #}\n{# line 2 #}\n{# line 3 #}",
			expectedTokens: []token.Token{
				{Typ: token.LCOMM},
				{Typ: token.COMM_TEXT, Val: " line 1 "},
				{Typ: token.RCOMM},
				{Typ: token.TEXT, Val: "\n"},
				{Typ: token.LCOMM},
				{Typ: token.COMM_TEXT, Val: " line 2 "},
				{Typ: token.RCOMM},
				{Typ: token.TEXT, Val: "\n"},
				{Typ: token.LCOMM},
				{Typ: token.COMM_TEXT, Val: " line 3 "},
				{Typ: token.RCOMM},
			},
		},
	}
	runTestCases(t, testCases)
}

func TestCommentsEdgeCases(t *testing.T) {
	testCases := []testCase{
		{
			name:  "Comment not terminated",
			input: "{#",
			expectedTokens: []token.Token{
				{Typ: token.LCOMM},
			},
		},
	}
	runTestCases(t, testCases)
}
