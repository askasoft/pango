package xmw

import (
	"github.com/askasoft/pango/xin"
)

// ResponseHeader response header middleware
type ResponseHeader struct {
	Header map[string]string
}

// NewResponseHeader create a default ResponseHeader
func NewResponseHeader(header map[string]string) *ResponseHeader {
	return &ResponseHeader{Header: header}
}

// Handle process xin request
func (rh *ResponseHeader) Handle(c *xin.Context) {
	for k, v := range rh.Header {
		c.Header(k, v)
	}
	c.Next()
}
