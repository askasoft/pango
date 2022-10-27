package gmw

import (
	"net/http"

	"github.com/pandafw/pango/gin"
)

// RequestLimiter http request limit middleware
type RequestLimiter struct {
	MaxBodySize int64
}

// NewRequestLimiter create a default RequestLimiter
func NewRequestLimiter(maxBodySize int64) *RequestLimiter {
	return &RequestLimiter{MaxBodySize: maxBodySize}
}

// Handler returns the gin.HandlerFunc
func (rl *RequestLimiter) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		rl.handle(c)
	}
}

// handle process gin request
func (rl *RequestLimiter) handle(c *gin.Context) {
	if rl.MaxBodySize > 0 {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, rl.MaxBodySize)
	}
	c.Next()
}
