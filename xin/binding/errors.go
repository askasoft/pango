package binding

import (
	"fmt"
	"strings"
)

// FieldBindError bind error
type FieldBindError struct {
	Field  string
	Cause  error
	Values []string
}

// Error return a string representing the bind error
func (be *FieldBindError) Error() string {
	return fmt.Sprintf("FieldBindError: %s: %s - %v", be.Field, be.Cause, be.Values)
}

// FieldBindErrors bind errors
type FieldBindErrors []*FieldBindError

// Error return a string representing the bind errors
func (bes FieldBindErrors) Error() string {
	var sb strings.Builder
	for i, e := range bes {
		if i > 0 {
			sb.WriteRune('\n')
		}
		sb.WriteString(e.Error())
	}
	return sb.String()
}
