package xin

import (
	"net"
	"net/http"

	"github.com/askasoft/pango/str"
)

// RecoveryFunc defines the function passable to CustomRecovery.
type RecoveryFunc func(c *Context, err any)

// Recovery returns a middleware that recovers from any panics and writes a 500 if there was one.
func Recovery() HandlerFunc {
	return CustomRecovery(defaultRecover)
}

// CustomRecovery returns a middleware that recovers from any panics and calls the provided handle func to handle it.
func CustomRecovery(r RecoveryFunc) HandlerFunc {
	return func(c *Context) {
		defer Recover(c, r)

		c.Next()
	}
}

func defaultRecover(c *Context, err any) {
	if IsBrokenPipeError(err) {
		c.Logger.Warnf("Broken (//%s%s): %v", c.Request.Host, c.Request.URL, err)

		// connection is dead, we can't write a status to it.
		c.Abort()
		return
	}

	c.Logger.Errorf("Panic (//%s%s): %v", c.Request.Host, c.Request.URL, err)
	c.AbortWithStatus(http.StatusInternalServerError)
}

func Recover(c *Context, r RecoveryFunc) {
	defer func() {
		if err := recover(); err != nil {
			defaultRecover(c, err)
		}
	}()

	if err := recover(); err != nil {
		r(c, err)
	}
}

var BrokenPipeErrors = []string{
	"broken pipe",
	"connection reset by peer",
	"i/o timeout",
}

// IsBrokenPipeError Check for a broken connection error
func IsBrokenPipeError(err any) bool {
	if err != nil {
		// Check for a broken connection, as it is not really a
		// condition that warrants a panic stack trace.
		if ne, ok := err.(*net.OpError); ok {
			se := ne.Unwrap().Error()
			for _, s := range BrokenPipeErrors {
				if str.ContainsFold(s, se) {
					return true
				}
			}
		}
	}
	return false
}
