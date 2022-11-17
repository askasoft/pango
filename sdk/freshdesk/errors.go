package freshdesk

import (
	"fmt"
	"strings"
)

type RateLimitedError struct {
	// retry after seconds
	RetryAfter int
}

func (e *RateLimitedError) Error() string {
	return fmt.Sprintf("Retry-After: %d seconds", e.RetryAfter)
}

type Error struct {
	Code    string `json:"code,omitempty"`
	Field   string `json:"field,omitempty"`
	Message string `json:"message,omitempty"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s: %s", e.Code, e.Field, e.Message)
}

type ErrorResult struct {
	Description string   `json:"description,omitempty"`
	Errors      []*Error `json:"errors,omitempty"`
}

func (er *ErrorResult) Error() string {
	var sb strings.Builder
	sb.WriteString(er.Description)
	for _, e := range er.Errors {
		sb.WriteRune('\n')
		sb.WriteString(e.Error())
	}
	return sb.String()
}
