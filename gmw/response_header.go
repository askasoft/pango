package gmw

import (
	"github.com/pandafw/pango/gin"
)

// ResponseHeader response header middleware
type ResponseHeader struct {
	Header map[string]string
}

// NewResponseHeader create a default ResponseHeader
func NewResponseHeader(header map[string]string) *ResponseHeader {
	return &ResponseHeader{Header: header}
}

// Handler returns the gin.HandlerFunc
func (rh *ResponseHeader) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		rh.handle(c)
	}
}

// handle process gin request
func (rh *ResponseHeader) handle(c *gin.Context) {
	for k, v := range rh.Header {
		c.Header(k, v)
	}
	c.Next()
}
