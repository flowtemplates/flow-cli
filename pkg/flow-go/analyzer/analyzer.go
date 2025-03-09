package analyzer

import (
	"fmt"

	"github.com/flowtemplates/cli/pkg/flow-go/parser"
	"github.com/flowtemplates/cli/pkg/flow-go/token"
	"github.com/flowtemplates/cli/pkg/flow-go/types"
)

type Variable struct {
	Name string
	Typ  types.Type
}

type TypeMap map[string]types.Type

func addToTypeMap(v Variable, tm TypeMap) error {
	if typ, exists := tm[v.Name]; !exists || typ == types.Any {
		tm[v.Name] = v.Typ
	} else if v.Typ != types.Any && v.Typ != typ {
		return fmt.Errorf("unmatched type of %q, got %s, expected %s", v.Name, v.Typ, typ)
	}

	return nil
}

func parseExpressionTypes(expr parser.Expr, tm TypeMap, errs *[]error) types.Type {
	switch e := expr.(type) {
	case parser.Ident:
		if err := addToTypeMap(Variable{
			Name: e.Name,
			Typ:  types.Any,
		}, tm); err != nil {
			*errs = append(*errs, err)
		}
		return types.Any

	case parser.Lit:
		switch e.Typ {
		case token.INT, token.FLOAT:
			return types.Number
		case token.STR:
			return types.String
		default:
			return types.Any
		}
	case parser.BinaryExpr:
		t1 := parseExpressionTypes(e.X, tm, errs)
		t2 := parseExpressionTypes(e.Y, tm, errs)

		if e.Op == token.ADD {
			if t1 == types.String || t2 == types.String {
				// If one side is a string, enforce string type
				if ident, ok := e.X.(parser.Ident); ok {
					if err := addToTypeMap(Variable{Name: ident.Name, Typ: types.String}, tm); err != nil {
						*errs = append(*errs, err)
					}
				}
				if ident, ok := e.Y.(parser.Ident); ok {
					if err := addToTypeMap(Variable{Name: ident.Name, Typ: types.String}, tm); err != nil {
						*errs = append(*errs, err)
					}
				}
				return types.String
			} else if t1 == types.Number || t2 == types.Number {
				// If one side is a number, enforce number type
				if ident, ok := e.X.(parser.Ident); ok {
					if err := addToTypeMap(Variable{Name: ident.Name, Typ: types.Number}, tm); err != nil {
						*errs = append(*errs, err)
					}
				}
				if ident, ok := e.Y.(parser.Ident); ok {
					if err := addToTypeMap(Variable{Name: ident.Name, Typ: types.Number}, tm); err != nil {
						*errs = append(*errs, err)
					}
				}

				return types.Number
			}
		}

		// If neither inference rule applies, both variables remain Any
		return types.Any
	}

	return types.Any
}
