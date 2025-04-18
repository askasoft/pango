package xmw

import (
	"errors"
	"net/http"

	"github.com/askasoft/pango/iox"
	"github.com/askasoft/pango/xin"
)

// RequestSizeLimiter http request limit middleware
type RequestSizeLimiter struct {
	MaxBodySize  int64
	DrainBody    bool // drain request body if we are under apache, otherwise the apache will return 502 Bad Gateway
	BodyTooLarge func(c *xin.Context)
}

// NewRequestSizeLimiter create a default RequestSizeLimiter middleware
func NewRequestSizeLimiter(maxBodySize int64) *RequestSizeLimiter {
	return &RequestSizeLimiter{MaxBodySize: maxBodySize}
}

// Handle process xin request
func (rsl *RequestSizeLimiter) Handle(c *xin.Context) {
	if rsl.MaxBodySize <= 0 {
		c.Next()
		return
	}

	var err error

	if c.Request.ContentLength > rsl.MaxBodySize {
		err = &iox.MaxBytesError{Limit: rsl.MaxBodySize}
	} else {
		crb := c.Request.Body
		mbr := iox.NewMaxBytesReader(crb, rsl.MaxBodySize)
		c.Request.Body = mbr
		c.Next()
		c.Request.Body = crb

		err = mbr.Error()
	}

	if err != nil {
		var mbe *iox.MaxBytesError
		if ok := errors.As(err, &mbe); ok {
			if rsl.DrainBody {
				iox.Drain(c.Request.Body)
			}

			if btl := rsl.BodyTooLarge; btl != nil {
				btl(c)
			} else {
				c.String(http.StatusRequestEntityTooLarge, mbe.Error())
				c.Abort()
			}
		}
	}
}
