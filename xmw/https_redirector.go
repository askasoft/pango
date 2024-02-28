package xmw

import (
	"net/http"
	"strings"

	"github.com/askasoft/pango/xin"
)

// HTTPSRedirector is a middleware that helps setup a https redirect features.
type HTTPSRedirector struct {
	disabled bool

	// If TemporaryRedirect is true, the a 302 will be used while redirecting. Default is false (301).
	TemporaryRedirect bool

	// SSLHost is the host name that is used to redirect http requests to https. Default is "", which indicates to use the same host.
	SSLHost string

	// ProxyHeaders is set of header keys with associated values that would indicate a valid https request.
	// Useful when behind a Proxy Server(Apache, Nginx).
	// Default is `map[string]string{"X-Forwarded-Proto": "https"}`.
	ProxyHeaders map[string]string
}

func NewHTTPSRedirector() *HTTPSRedirector {
	return &HTTPSRedirector{
		ProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
	}
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

	r := c.Request

	// SSL check.
	isSSL := false
	if strings.EqualFold(r.URL.Scheme, "https") || r.TLS != nil {
		isSSL = true
	} else {
		sslProxyHeaders := sh.ProxyHeaders
		for k, v := range sslProxyHeaders {
			if c.GetHeader(k) == v {
				isSSL = true
				break
			}
		}
	}

	if isSSL {
		c.Next()
		return
	}

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
