package gmw

import (
	"github.com/pandafw/pango/gin"
)

// HTTPHeader http header middleware
type HTTPHeader struct {
	ResponseHeader map[string]string
}

// NewHTTPHeader create a default HTTPHeader
func NewHTTPHeader(resHeader map[string]string) *HTTPHeader {
	return &HTTPHeader{ResponseHeader: resHeader}
}

// Handler returns the gin.HandlerFunc
func (hh *HTTPHeader) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		hh.handle(c)
	}
}

// handle process gin request
func (hh *HTTPHeader) handle(c *gin.Context) {
	for k, v := range hh.ResponseHeader {
		c.Header(k, v)
	}
	c.Next()
}
