package renderer

import (
	"fmt"
	"strings"

	"github.com/flowtemplates/cli/pkg/flow-go/parser"
)

type Scope map[string]any

func Render(ast []parser.Node, context Scope) (string, error) {
	var result strings.Builder
	for _, node := range ast {
		switch n := node.(type) {
		case parser.Text:
			result.WriteString(n.Val)
		case parser.ExprBlock:
			switch body := n.Body.(type) {
			case parser.Ident:
				return identToString(&body, context)
			case *parser.Ident:
				return identToString(body, context)
			default:
				return "", fmt.Errorf("unsupported expr type: %T", body)
			}
		case parser.IfStmt:
			conditionValue, err := evaluateCondition(n.Condition, context)
			if err != nil {
				return "", err
			}

			if conditionValue != "" {
				bodyContent, err := Render(n.Body, context)
				if err != nil {
					return "", err
				}

				result.WriteString(bodyContent)
			} else if n.Else != nil {
				elseContent, err := Render(n.Else, context)
				if err != nil {
					return "", err
				}

				result.WriteString(elseContent)
			}
		}
	}

	return result.String(), nil
}

func identToString(ident *parser.Ident, context Scope) (string, error) {
	value, exists := context[ident.Name]
	if !exists {
		return "", fmt.Errorf("%s not declared", ident.Name)
	}

	return valueToString(value), nil
}

func evaluateCondition(cond parser.Expr, context Scope) (string, error) {
	switch n := cond.(type) {
	case parser.Ident:
		value, exists := context[n.Name]
		if !exists {
			return "", fmt.Errorf("%s not declared", n.Name)
		}

		return valueToString(value), nil
	default:
		return "", fmt.Errorf("unsupported condition type: %T", cond)
	}
}

func valueToString(value any) string {
	switch v := value.(type) {
	case int:
		return fmt.Sprintf("%d", v)
	case float64:
		return fmt.Sprintf("%f", v)
	case bool:
		return ""
	default:
		return v.(string)
	}
}
