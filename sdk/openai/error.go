package openai

import (
	"fmt"
	"strings"

	"github.com/askasoft/pango/sdk"
)

type RateLimitedError = sdk.RateLimitedError

type ErrorDetail struct {
	Type    string `json:"type,omitempty"`
	Code    any    `json:"code,omitempty"`
	Param   any    `json:"param,omitempty"`
	Message string `json:"message,omitempty"`
}

func (ed *ErrorDetail) String() string {
	var sb strings.Builder
	if ed.Type != "" {
		sb.WriteString(ed.Type)
	}
	if ed.Code != nil {
		s := fmt.Sprint(ed.Code)
		if s != "" {
			if sb.Len() > 0 {
				sb.WriteByte('/')
			}
			sb.WriteString(s)
		}
	}
	if ed.Param != nil {
		s := fmt.Sprint(ed.Param)
		if s != "" {
			if sb.Len() > 0 {
				sb.WriteByte('/')
			}
			sb.WriteString(s)
		}
	}
	if ed.Message != "" {
		if sb.Len() > 0 {
			sb.WriteString(": ")
		}
		sb.WriteString(ed.Message)
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
