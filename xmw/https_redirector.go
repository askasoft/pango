package xmw

import (
	"net/http"

	"github.com/askasoft/pango/xin"
)

// HTTPSRedirector is a middleware that helps setup a https redirect features.
type HTTPSRedirector struct {
	disabled bool

	// If TemporaryRedirect is true, the a 302 will be used while redirecting. Default is false (301).
	TemporaryRedirect bool

	// SSLHost is the host name that is used to redirect http requests to https. Default is "", which indicates to use the same host.
	SSLHost string
}

func NewHTTPSRedirector() *HTTPSRedirector {
	return &HTTPSRedirector{}
}

// Disable disable the secure handler or not
func (sh *HTTPSRedirector) Disable(disabled bool) {
	sh.disabled = disabled
}

// Handler returns the xin.HandlerFunc
func (sh *HTTPSRedirector) Handler() xin.HandlerFunc {
	return sh.Handle
}

// Handle process xin request
func (sh *HTTPSRedirector) Handle(c *xin.Context) {
	if sh.disabled {
		c.Next()
		return
	}

	if c.IsSecure() {
		c.Next()
		return
	}

	r := c.Request

	url := r.URL
	url.Scheme = "https"
	url.Host = r.Host

	sslHost := sh.SSLHost
	if sslHost != "" {
		url.Host = sslHost
	}

	status := http.StatusMovedPermanently
	if sh.TemporaryRedirect {
		status = http.StatusTemporaryRedirect
	}

	c.Redirect(status, url.String())
	c.Abort()
}
