package parser

import "fmt"

type Error struct {
	Pos     int
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("Parse error at position %d: %s", e.Pos, e.Message)
}

type ErrorList []Error
