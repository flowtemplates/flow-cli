package renderer

import (
	"fmt"
	"strings"

	"github.com/flowtemplates/cli/pkg/flow-go/lexer"
	"github.com/flowtemplates/cli/pkg/flow-go/parser"
)

func RenderAst(ast []parser.Node, context Scope) (string, error) {
	var result strings.Builder
	for _, node := range ast {
		switch n := node.(type) {
		case parser.Text:
			result.WriteString(n.Val)
		case parser.ExprBlock:
			switch body := n.Body.(type) {
			case parser.Ident:
				s, err := identToString(&body, context)
				if err != nil {
					return "", err
				}

				result.WriteString(s)
			// case parser.Lit:
			// result.WriteString(valueToString(body.Val))
			default:
				return "", fmt.Errorf("unsupported expr type: %T", body)
			}
		case parser.IfStmt:
			conditionValue, err := evaluateCondition(n.Condition, context)
			if err != nil {
				return "", err
			}

			if isFalsy(conditionValue) {
				bodyContent, err := RenderAst(n.Body, context)
				if err != nil {
					return "", err
				}

				result.WriteString(bodyContent)
			} else if n.Else != nil {
				elseContent, err := RenderAst(n.Else, context)
				if err != nil {
					return "", err
				}

				result.WriteString(elseContent)
			}
		}
	}

	return result.String(), nil
}

func RenderString(input string, tm Scope) (string, error) {
	tokens := lexer.TokensFromString(input)
	ast, errs := parser.New(tokens).Parse()
	if len(errs) != 0 {
		return "", errs[0]
	}

	res, err := RenderAst(ast, tm)
	if err != nil {
		return "", fmt.Errorf("failed to render: %w", err)
	}

	return res, nil
}
