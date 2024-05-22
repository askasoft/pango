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

func (fbes FieldBindErrors) As(err any) bool {
	if pp, ok := err.(**FieldBindErrors); ok {
		*pp = &fbes
		return true
	}
	return false
}

func (fbes FieldBindErrors) Unwrap() []error {
	errs := make([]error, len(fbes))
	for i, fbe := range fbes {
		errs[i] = fbe
	}
	return errs
}

// Error return a string representing the bind errors
func (fbes FieldBindErrors) Error() string {
	var sb strings.Builder
	for i, fbe := range fbes {
		if i > 0 {
			sb.WriteRune('\n')
		}
		sb.WriteString(fbe.Error())
	}
	return sb.String()
}
