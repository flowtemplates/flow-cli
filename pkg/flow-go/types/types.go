package types

import "strconv"

type Type string

const (
	Number  Type = "number"
	String  Type = "string"
	Boolean Type = "boolean"
	Any     Type = "any"
	// TypeArray   Type = "array"
	// TypeObject  Type = "object"
)

func (t Type) IsValid(val string) bool {
	switch t {
	case Number:
		_, err := strconv.Atoi(val)
		return err == nil
	case Boolean:
		return val == "true" || val == "false"
	default:
		return true
	}
}

func (t Type) GetDefaultValue() string {
	switch t {
	case Number:
		return "0"
	case Boolean:
		return "false"
	default:
		return ""
	}
}
