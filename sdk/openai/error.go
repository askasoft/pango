package openai

import (
	"fmt"
	"strings"

	"github.com/askasoft/pango/sdk"
)

type RateLimitedError = sdk.RateLimitedError

type ErrorDetail struct {
	Type    string `json:"type,omitempty"`
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Param   any    `json:"param,omitempty"`
}

func (ed *ErrorDetail) String() string {
	var sb strings.Builder
	if ed.Type != "" {
		sb.WriteString(ed.Type)
	}
	if ed.Code != "" {
		if sb.Len() > 0 {
			sb.WriteString(": ")
		}
		sb.WriteString(ed.Code)
	}
	if ed.Message != "" {
		if sb.Len() > 0 {
			sb.WriteString(": ")
		}
		sb.WriteString(ed.Message)
	}
	if ed.Param != nil {
		if sb.Len() > 0 {
			sb.WriteString(": ")
		}
		sb.WriteString(fmt.Sprint(ed.Param))
	}
	return sb.String()
}

type ErrorResult struct {
	StatusCode int          `json:"-"` // http status code
	Status     string       `json:"-"` // http status
	Detail     *ErrorDetail `json:"error,omitempty"`
}

func (er *ErrorResult) Error() string {
	detail := ""
	if er.Detail != nil {
		detail = er.Detail.String()
	}

	if detail != "" {
		return er.Status + " - " + detail
	}

	return er.Status
}
