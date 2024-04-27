package binding

import (
	"fmt"
	"strings"
)

// FieldBindError bind error
type FieldBindError struct {
	Err    error
	Field  string
	Values []string
}

// Error return a string representing the bind error
func (fbe *FieldBindError) Error() string {
	return fmt.Sprintf("FieldBindError: %s: %s - %v", fbe.Field, fbe.Err, fbe.Values)
}

func (fbe *FieldBindError) Unwrap() error {
	return fbe.Err
}

// FieldBindErrors bind errors
type FieldBindErrors []*FieldBindError

// Error return a string representing the bind errors
func (fbes FieldBindErrors) Error() string {
	var sb strings.Builder
	for i, e := range fbes {
		if i > 0 {
			sb.WriteRune('\n')
		}
		sb.WriteString(e.Error())
	}
	return sb.String()
}
