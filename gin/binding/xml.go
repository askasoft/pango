package binding

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
)

type XmlBindError struct {
	Cause error
}

// Error return a string representing the bind error
func (jbe *XmlBindError) Error() string {
	return fmt.Sprintf("XmlBindError: %v", jbe.Cause)
}

type xmlBinding struct{}

func (xmlBinding) Name() string {
	return "xml"
}

func (xmlBinding) Bind(req *http.Request, obj any) error {
	return decodeXML(req.Body, obj)
}

func (xmlBinding) BindBody(body []byte, obj any) error {
	return decodeXML(bytes.NewReader(body), obj)
}

func decodeXML(r io.Reader, obj any) error {
	if err := xml.NewDecoder(r).Decode(obj); err != nil {
		return &XmlBindError{err}
	}
	return nil
}
