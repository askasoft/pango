package fdk

import (
	"fmt"
	"strings"
	"time"

	"github.com/askasoft/pango/str"
)

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
	RetryAfter  time.Duration
}

func (er *ErrorResult) GetRetryAfter() time.Duration {
	return er.RetryAfter
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
	es := er.Status

	if er.RetryAfter > 0 {
		es = fmt.Sprintf("%s (Retry After %s)", es, er.RetryAfter)
	}

	detail := er.Detail()
	if detail != "" {
		es = es + " - " + detail
	}

	return es
}
