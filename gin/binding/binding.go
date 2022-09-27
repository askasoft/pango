package binding

import (
	"fmt"
	"net/http"
	"strings"
)

// Content-Type MIME of the most common data formats.
const (
	MIMEJSON              = "application/json"
	MIMEHTML              = "text/html"
	MIMEXML               = "application/xml"
	MIMEXML2              = "text/xml"
	MIMEPlain             = "text/plain"
	MIMEPOSTForm          = "application/x-www-form-urlencoded"
	MIMEMultipartPOSTForm = "multipart/form-data"
)

// BindError bind error
type BindError struct {
	Name   string
	Values []string
	Cause  error
}

// Error return a string representing the bind error
func (be *BindError) Error() string {
	return fmt.Sprintf("%s: %s - %v", be.Name, be.Cause, be.Values)
}

// BindErrors bind errors
type BindErrors struct {
	Errors []*BindError
}

func (bes *BindErrors) IsEmpty() bool {
	return len(bes.Errors) == 0
}

func (bes *BindErrors) AddError(be *BindError) {
	bes.Errors = append(bes.Errors, be)
}

// Error return a string representing the bind errors
func (bes *BindErrors) Error() string {
	var sb strings.Builder
	for i, e := range bes.Errors {
		if i > 0 {
			sb.WriteRune('\n')
		}
		sb.WriteString(e.Error())
	}
	return sb.String()
}

// Binding describes the interface which needs to be implemented for binding the
// data present in the request such as JSON request body, query parameters or
// the form POST.
type Binding interface {
	Name() string
	Bind(*http.Request, any) error
}

// BodyBinding adds BindBody method to Binding. BindBody is similar with Bind,
// but it reads the body from supplied bytes instead of req.Body.
type BodyBinding interface {
	Binding
	BindBody([]byte, any) error
}

// URIBinding adds BindUri method to Binding. BindUri is similar with Bind,
// but it reads the Params.
type URIBinding interface {
	Name() string
	BindUri(map[string][]string, any) error
}

// These implement the Binding interface and can be used to bind the data
// present in the request to struct instances.
var (
	JSON          = jsonBinding{}
	XML           = xmlBinding{}
	Form          = formBinding{}
	Query         = queryBinding{}
	FormPost      = formPostBinding{}
	FormMultipart = formMultipartBinding{}
	URI           = uriBinding{}
	Header        = headerBinding{}
)

// Default returns the appropriate Binding instance based on the HTTP method
// and the content type.
func Default(method, contentType string) Binding {
	if method == http.MethodGet {
		return Form
	}

	switch contentType {
	case MIMEJSON:
		return JSON
	case MIMEXML, MIMEXML2:
		return XML
	case MIMEMultipartPOSTForm:
		return FormMultipart
	default: // case MIMEPOSTForm:
		return Form
	}
}
