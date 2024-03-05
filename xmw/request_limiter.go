package xmw

import (
	"errors"
	"net/http"

	"github.com/askasoft/pango/iox"
	"github.com/askasoft/pango/xin"
)

// RequestLimiter http request limit middleware
type RequestLimiter struct {
	MaxBodySize  int64
	DrainBody    bool // drain request body if we are under apache, otherwise the apache will return 502 Bad Gateway
	BodyTooLarge func(c *xin.Context, maxBodySize int64)
}

// NewRequestLimiter create a default RequestLimiter middleware
func NewRequestLimiter(maxBodySize int64, bodyTooLarge ...func(c *xin.Context, maxBodySize int64)) *RequestLimiter {
	rl := &RequestLimiter{MaxBodySize: maxBodySize}
	if len(bodyTooLarge) > 0 {
		rl.BodyTooLarge = bodyTooLarge[0]
	}
	return rl
}

// Handler returns the xin.HandlerFunc
func (rl *RequestLimiter) Handler() xin.HandlerFunc {
	return rl.Handle
}

// Handle process xin request
func (rl *RequestLimiter) Handle(c *xin.Context) {
	if rl.MaxBodySize <= 0 {
		c.Next()
		return
	}

	var err error

	if c.Request.ContentLength > rl.MaxBodySize {
		err = &iox.MaxBytesError{Limit: rl.MaxBodySize}
	} else {
		crb := c.Request.Body
		mbr := iox.NewMaxBytesReader(crb, rl.MaxBodySize)
		c.Request.Body = mbr
		c.Next()
		c.Request.Body = crb

		err = mbr.Error()
	}

	if err != nil {
		var mbe *iox.MaxBytesError
		if ok := errors.As(err, &mbe); ok {
			if rl.DrainBody {
				iox.Drain(c.Request.Body)
			}

			btl := rl.BodyTooLarge
			if btl != nil {
				btl(c, rl.MaxBodySize)
			} else {
				c.String(http.StatusBadRequest, mbe.Error())
				c.Abort()
			}
		}
	}
}
