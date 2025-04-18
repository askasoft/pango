package xmw

import (
	"net/http"

	"github.com/askasoft/pango/cog/hashset"
	"github.com/askasoft/pango/num"
	"github.com/askasoft/pango/xin"
)

// OriginAccessController Access-Control-Allow-Origin response header middleware
// see
// - https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Origin
// - https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Credentials
// - https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Headers
// - https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Methods
// - https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Expose-Headers
// - https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Max-Age
// if Origins contains the request header 'origin', set the Access-Control-Allow-Origin response header.
// if the request method is OPTIONS, also set the status code to 200.
type OriginAccessController struct {
	AllowOrigins     *hashset.HashSet[string]
	AllowCredentials bool
	AllowMethods     string
	AllowHeaders     string
	ExposeHeaders    string
	MaxAge           int
}

// NewOriginAccessController create a default OriginAccessController
func NewOriginAccessController(origins ...string) *OriginAccessController {
	return &OriginAccessController{AllowOrigins: hashset.NewHashSet(origins...)}
}

// SetAllowOrigins set allow origins
func (ll *OriginAccessController) SetAllowOrigins(origins ...string) {
	ll.AllowOrigins = hashset.NewHashSet(origins...)
}

// SetAllowMethods set Access-Control-Allow-Methods
func (ll *OriginAccessController) SetAllowMethods(methods string) {
	ll.AllowMethods = methods
}

// SetAllowCredentials set allow Credentials
func (ll *OriginAccessController) SetAllowCredentials(allow bool) {
	ll.AllowCredentials = allow
}

// SetAllowHeaders set allow headers
func (ll *OriginAccessController) SetAllowHeaders(headers string) {
	ll.AllowHeaders = headers
}

// SetExposeHeaders set expose headers
func (ll *OriginAccessController) SetExposeHeaders(headers string) {
	ll.ExposeHeaders = headers
}

// SetMaxAge set Access-Control-Max-Age
func (ll *OriginAccessController) SetMaxAge(sec int) {
	ll.MaxAge = sec
}

// Handle process xin request
func (ll *OriginAccessController) Handle(c *xin.Context) {
	acaos := ll.AllowOrigins
	if acaos.Len() > 0 {
		acao := ""

		if acaos.Contains("*") {
			acao = "*"
		} else {
			origin := c.GetHeader("Origin")
			if origin != "" && acaos.Contains(origin) {
				acao = origin
			}
		}

		if acao != "" {
			c.Header("Access-Control-Allow-Origin", acao)

			acac := ll.AllowCredentials
			if acac {
				c.Header("Access-Control-Allow-Credentials", "true")
			}

			acams := ll.AllowMethods
			if acams != "" {
				c.Header("Access-Control-Allow-Methods", acams)
			}

			acahs := ll.AllowHeaders
			if acahs != "" {
				c.Header("Access-Control-Allow-Headers", acahs)
			}

			acehs := ll.ExposeHeaders
			if acehs != "" {
				c.Header("Access-Control-Expose-Headers", acehs)
			}

			acma := ll.MaxAge
			if acma > 0 {
				c.Header("Access-Control-Max-Age", num.Itoa(acma))
			}

			if c.Request.Method == http.MethodOptions {
				c.Status(http.StatusOK)
			}
		}
	}

	c.Next()
}
