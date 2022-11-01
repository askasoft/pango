package xmw

import (
	"github.com/pandafw/pango/xin"
)

// ResponseHeader response header middleware
type ResponseHeader struct {
	Header map[string]string
}

// NewResponseHeader create a default ResponseHeader
func NewResponseHeader(header map[string]string) *ResponseHeader {
	return &ResponseHeader{Header: header}
}

// Handler returns the xin.HandlerFunc
func (rh *ResponseHeader) Handler() xin.HandlerFunc {
	return func(c *xin.Context) {
		rh.handle(c)
	}
}

// handle process xin request
func (rh *ResponseHeader) handle(c *xin.Context) {
	for k, v := range rh.Header {
		c.Header(k, v)
	}
	c.Next()
}
