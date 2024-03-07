package xin

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestPanicWithAbort assert that panic has been recovered even if context.Abort was used.
func TestPanicWithAbort(t *testing.T) {
	e := "Oupps, Houston, we have a problem"

	router := New()
	router.Use(Recovery())
	router.GET("/recovery", func(c *Context) {
		c.AbortWithStatus(http.StatusBadRequest)
		panic(e)
	})

	w := performRequest(router, "GET", "/recovery")
	assert.Equal(t, http.StatusBadRequest, w.Code)

	b := fmt.Sprintf("%d %s", http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
	b += fmt.Sprintf("%d %s", http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	assert.Equal(t, b, w.Body.String())
}

// TestPanicWithBrokenPipe asserts that recovery specifically handles
// writing responses to broken pipes
func TestPanicWithBrokenPipe(t *testing.T) {
	const expectCode = 204

	expectMsgs := map[syscall.Errno]string{
		syscall.EPIPE:      "broken pipe",
		syscall.ECONNRESET: "connection reset by peer",
	}

	for errno, expectMsg := range expectMsgs {
		t.Run(expectMsg, func(t *testing.T) {
			router := New()
			router.Use(Recovery())
			router.GET("/recovery", func(c *Context) {
				c.String(expectCode, expectMsg)

				// Oops. Client connection closed
				e := &net.OpError{Err: &os.SyscallError{Err: errno}}
				panic(e)
			})

			w := performRequest(router, "GET", "/recovery")
			assert.Equal(t, expectCode, w.Code)
			//assert.Equal(t, expectMsg, w.Body.String())
		})
	}
}

func TestPanicCustomRecovery(t *testing.T) {
	e := "Oupps, Houston, we have a problem"

	router := New()
	recovery := func(c *Context, err interface{}) {
		c.String(http.StatusBadRequest, err.(string))
	}
	router.Use(CustomRecovery(recovery))
	router.GET("/recovery", func(c *Context) {
		c.Status(http.StatusNoContent)
		panic(e)
	})

	w := performRequest(router, "GET", "/recovery")
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, e, w.Body.String())
}

func TestPanicInCustomRecovery(t *testing.T) {
	e := "Oupps, Houston, we have a problem. "

	router := New()
	recovery := func(c *Context, err interface{}) {
		c.String(http.StatusBadRequest, err.(string))
		panic(e + e)
	}
	router.Use(CustomRecovery(recovery))
	router.GET("/recovery", func(c *Context) {
		c.Status(http.StatusNoContent)
		panic(e)
	})

	w := performRequest(router, "GET", "/recovery")
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, e, w.Body.String())
}
