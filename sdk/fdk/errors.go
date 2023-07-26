package fdk

import (
	"fmt"
	"strings"

	"github.com/askasoft/pango/sdk"
)

type RateLimitedError = sdk.RateLimitedError

type Error struct {
	Code    string `json:"code,omitempty"`
	Field   string `json:"field,omitempty"`
	Message string `json:"message,omitempty"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s: %s", e.Code, e.Field, e.Message)
}

type ErrorResult struct {
	StatusCode  int      `json:"-"` // http status code
	Status      string   `json:"-"` // http status
	Code        string   `json:"code,omitempty"`
	Message     string   `json:"message,omitempty"`
	Description string   `json:"description,omitempty"`
	Errors      []*Error `json:"errors,omitempty"`
}

func (er *ErrorResult) Error() string {
	var sb strings.Builder
	sb.WriteString(er.Status)
	if er.Code != "" {
		sb.WriteString(": ")
		sb.WriteString(er.Code)
	}
	if er.Message != "" {
		sb.WriteString(": ")
		sb.WriteString(er.Message)
	}
	if er.Description != "" {
		sb.WriteString(": ")
		sb.WriteString(er.Description)
	}
	for _, e := range er.Errors {
		sb.WriteRune('\n')
		sb.WriteString(e.Error())
	}
	return sb.String()
}
