package middleware

import (
	"errors"
	"net/http"

	"github.com/askasoft/pango/iox"
	"github.com/askasoft/pango/xin"
)

// RequestSizeLimiter http request limit middleware
type RequestSizeLimiter struct {
	DrainBody      bool // drain request body if we are under apache, otherwise the apache will return 502 Bad Gateway
	MaxBodySize    int64
	GetMaxBodySize func(c *xin.Context) int64
	BodyTooLarge   func(c *xin.Context)
}

// NewRequestSizeLimiter create a default RequestSizeLimiter middleware
func NewRequestSizeLimiter(maxBodySize int64) *RequestSizeLimiter {
	return &RequestSizeLimiter{MaxBodySize: maxBodySize}
}

// Handle process xin request
func (rsl *RequestSizeLimiter) Handle(c *xin.Context) {
	mbs := rsl.MaxBodySize
	if gmbs := rsl.GetMaxBodySize; gmbs != nil {
		mbs = gmbs(c)
	}

	if mbs <= 0 {
		c.Next()
		return
	}

	var err error

	if c.Request.ContentLength > mbs {
		err = &http.MaxBytesError{Limit: mbs}
	} else {
		crb := c.Request.Body
		mbr := http.MaxBytesReader(c.Writer.Writer(), crb, mbs) // let http.ParseForm() check http.maxBytesReader
		c.Request.Body = mbr
		c.Next()
		c.Request.Body = crb

		_, err = mbr.Read(nil) // get last error
	}

	if err != nil {
		var mbe *http.MaxBytesError
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
