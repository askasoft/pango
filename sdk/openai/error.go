package openai

import (
	"fmt"
	"strings"

	"github.com/askasoft/pango/sdk"
)

type RateLimitedError = sdk.RateLimitedError

type ErrorResult struct {
	StatusCode int    `json:"-"` // http status code
	Status     string `json:"-"` // http status
	Code       string `json:"code,omitempty"`
	Message    string `json:"message,omitempty"`
	Type       string `json:"type,omitempty"`
	Param      any    `json:"param,omitempty"`
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
	if er.Param != nil {
		sb.WriteString(": ")
		sb.WriteString(fmt.Sprint(er.Param))
	}
	return sb.String()
}
