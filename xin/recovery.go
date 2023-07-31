package xin

import (
	"errors"
	"net"
	"net/http"
	"os"

	"github.com/askasoft/pango/str"
)

// RecoveryFunc defines the function passable to CustomRecovery.
type RecoveryFunc func(c *Context, err interface{})

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
	loggerRecover(c, err)
	c.AbortWithStatus(http.StatusInternalServerError)
}

func loggerRecover(c *Context, err any) {
	c.Logger.Errorf("Panic: %v", err)
}

func doRecovery(c *Context, err any, r RecoveryFunc) {
	if IsBrokenPipeError(err) {
		c.Logger.Debugf("Abort: %v", err)

		// connection is dead, we can't write a status to it.
		c.Abort()
		return
	}

	r(c, err)
}

func Recover(c *Context, r RecoveryFunc) {
	defer func() {
		if err := recover(); err != nil {
			doRecovery(c, err, loggerRecover)
		}
	}()

	if err := recover(); err != nil {
		doRecovery(c, err, r)
	}
}

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
