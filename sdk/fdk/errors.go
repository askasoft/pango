package fdk

import (
	"fmt"
	"strings"

	"github.com/askasoft/pango/sdk"
	"github.com/askasoft/pango/str"
)

type RateLimitedError = sdk.RateLimitedError

type Error struct {
	Code    string `json:"code,omitempty"`
	Field   string `json:"field,omitempty"`
	Message string `json:"message,omitempty"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("(%s: %s: %s)", e.Code, e.Field, e.Message)
}

type ErrorResult struct {
	StatusCode  int      `json:"-"` // http status code
	Status      string   `json:"-"` // http status
	Code        string   `json:"code,omitempty"`
	Message     string   `json:"message,omitempty"`
	Description string   `json:"description,omitempty"`
	Errors      []*Error `json:"errors,omitempty"`
}

func (er *ErrorResult) Detail() string {
	var sb strings.Builder

	if er.Code != "" {
		sb.WriteString(er.Code)
	}
	if er.Message != "" {
		if sb.Len() > 0 {
			sb.WriteByte('/')
		}
		sb.WriteString(er.Message)
	}
	if er.Description != "" {
		if sb.Len() > 0 {
			sb.WriteByte('/')
		}
		sb.WriteString(er.Description)
	}
	for i, e := range er.Errors {
		sb.WriteString(str.If(i == 0, ": ", ", "))
		sb.WriteString(e.Error())
	}

	return sb.String()
}

func (er *ErrorResult) Error() string {
	detail := er.Detail()

	if detail != "" {
		return er.Status + " - " + detail
	}

	return er.Status
}
