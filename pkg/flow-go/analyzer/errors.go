package analyzer

import (
	"fmt"
	"strings"

	"github.com/flowtemplates/cli/pkg/flow-go/types"
)

type TypeError struct {
	ExpectedType types.Type
	Name         string
	Val          string
}

func (e *TypeError) Error() string {
	return fmt.Sprintf("TypeError: Variable '%s' expected type '%s'", e.Name, e.ExpectedType)
}

type TypeErrors []TypeError

func (te TypeErrors) Error() string {
	messages := make([]string, len(te))
	for i, err := range te {
		messages[i] = err.Error()
	}

	return strings.Join(messages, "\n")
}
