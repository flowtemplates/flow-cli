package renderer

import (
	"fmt"
	"strconv"

	"github.com/flowtemplates/cli/pkg/flow-go/parser"
)

type Scope map[string]string

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
		return strconv.Itoa(v)
	case float64:
		return fmt.Sprintf("%f", v)
	case bool:
		return ""
	default:
		return fmt.Sprintf("%s", v)
	}
}

func isFalsy(value string) bool {
	switch value {
	case "", "false", "0":
		return false
	default:
		return true
	}
}
