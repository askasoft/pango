package xin

import (
	"errors"
	"net"
	"net/http"
	"os"

	"github.com/pandafw/pango/str"
)

// IsBrokenPipeError Check for a broken connection error
func IsBrokenPipeError(err any) bool {
	if err != nil {
		// Check for a broken connection, as it is not really a
		// condition that warrants a panic stack trace.
		if ne, ok := err.(*net.OpError); ok {
			var se *os.SyscallError
			if errors.As(ne, &se) {
				if str.ContainsFold(se.Error(), "broken pipe") || str.ContainsFold(se.Error(), "connection reset by peer") {
					return true
				}
			}
		}
	}
	return false
}

// Recovery returns a middleware that recovers from any panics and writes a 500 if there was one.
func Recovery() HandlerFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				if IsBrokenPipeError(err) {
					c.Logger().Debug("Abort: %v", err)
					// If the connection is dead, we can't write a status to it.
					c.Abort()
					return
				}

				c.Logger().Errorf("Panic: %v", err)
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()

		c.Next()
	}
}
