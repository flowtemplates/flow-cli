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

func GetTypeMap(ast []parser.Node) (TypeMap, []error) {
	tm := make(TypeMap)
	errs := []error{}

	for _, node := range ast {
		switch n := node.(type) {
		case parser.ExprBlock:
			switch e := n.Body.(type) {
			case parser.Ident:
				addToTypeMap(Variable{
					Name: e.Name,
					Typ:  types.String,
				}, tm, &errs)
			}
			handleExpression(n.Body, tm, &errs)
		case parser.IfStmt:
		}
	}

	return tm, errs
}

func addToTypeMap(v Variable, tm TypeMap, errs *[]error) {
	if typ, exists := tm[v.Name]; !exists || typ == types.Any {
		tm[v.Name] = v.Typ
	} else if v.Typ != types.Any && v.Typ != typ {
		*errs = append(*errs, fmt.Errorf("unmatched type of %q, got %s, expected %s", v.Name, v.Typ, typ))
	}
}

func handleExpression(expr parser.Expr, tm TypeMap, errs *[]error) types.Type {
	switch e := expr.(type) {
	case parser.Ident:
		addToTypeMap(Variable{
			Name: e.Name,
			Typ:  types.Any,
		}, tm, errs)
		return types.Any

	case parser.Lit:
		// Determine type from literal token type
		if e.Typ == token.INT {
			return types.Number
		} else if e.Typ == token.STR {
			return types.String
		}
		return types.Any

	case parser.BinaryExpr:
		t1 := handleExpression(e.X, tm, errs)
		t2 := handleExpression(e.Y, tm, errs)

		if e.Op == token.ADD {
			if t1 == types.String || t2 == types.String {
				// If one side is a string, enforce string type
				if ident, ok := e.X.(parser.Ident); ok {
					addToTypeMap(Variable{Name: ident.Name, Typ: types.String}, tm, errs)
				}
				if ident, ok := e.Y.(parser.Ident); ok {
					addToTypeMap(Variable{Name: ident.Name, Typ: types.String}, tm, errs)
				}
				return types.String
			} else if t1 == types.Number || t2 == types.Number {
				// If one side is a number, enforce number type
				if ident, ok := e.X.(parser.Ident); ok {
					addToTypeMap(Variable{Name: ident.Name, Typ: types.Number}, tm, errs)
				}
				if ident, ok := e.Y.(parser.Ident); ok {
					addToTypeMap(Variable{Name: ident.Name, Typ: types.Number}, tm, errs)
				}
				return types.Number
			}
		}

		// If neither inference rule applies, both variables remain Any
		return types.Any
	}
	return types.Any
}
