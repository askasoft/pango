package binding

import (
	"fmt"
	"strings"
)

// FieldBindError bind error
type FieldBindError struct {
	Name   string
	Cause  error
	Values []string
}

// Error return a string representing the bind error
func (be *FieldBindError) Error() string {
	return fmt.Sprintf("FieldBindError: %s: %s - %v", be.Name, be.Cause, be.Values)
}

// FieldBindErrors bind errors
type FieldBindErrors struct {
	Errors []*FieldBindError
}

func (bes *FieldBindErrors) IsEmpty() bool {
	return len(bes.Errors) == 0
}

func (bes *FieldBindErrors) AddError(be *FieldBindError) {
	bes.Errors = append(bes.Errors, be)
}

// Error return a string representing the bind errors
func (bes *FieldBindErrors) Error() string {
	var sb strings.Builder
	for i, e := range bes.Errors {
		if i > 0 {
			sb.WriteRune('\n')
		}
		sb.WriteString(e.Error())
	}
	return sb.String()
}
