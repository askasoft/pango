package xmw

import (
	"net/http"

	"github.com/askasoft/pango/cog"
	"github.com/askasoft/pango/xin"
)

// OriginAccessController Access-Control-Allow-Origin response header middleware
// see https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Origin
// if Origins contains the request header 'origin', set the Access-Control-Allow-Origin response header.
// if the request method is OPTIONS, also set the status code to 200.
type OriginAccessController struct {
	Origins *cog.HashSet[string]
	Headers string
}

// NewOriginAccessController create a default OriginAccessController
func NewOriginAccessController(origins ...string) *OriginAccessController {
	return &OriginAccessController{Origins: cog.NewHashSet(origins...)}
}

// Handler returns the xin.HandlerFunc
func (ll *OriginAccessController) Handler() xin.HandlerFunc {
	return ll.Handle
}

// SetAllowOrigins set allowed origins
func (ll *OriginAccessController) SetAllowOrigins(origins ...string) {
	ll.Origins = cog.NewHashSet(origins...)
}

// SetAllowHeaders set allowed headers
func (ll *OriginAccessController) SetAllowHeaders(headers string) {
	ll.Headers = headers
}

// Handle process xin request
func (ll *OriginAccessController) Handle(c *xin.Context) {
	acaos := ll.Origins
	if acaos.Len() > 0 {
		acao := ""

		if acaos.Contain("*") {
			acao = "*"
		} else {
			origin := c.GetHeader("Origin")
			if origin != "" && acaos.Contain(origin) {
				acao = origin
			}
		}

		if acao != "" {
			c.Header("Access-Control-Allow-Origin", acao)
			acahs := ll.Headers
			if acahs != "" {
				c.Header("Access-Control-Allow-Headers", acahs)
			}
			if c.Request.Method == http.MethodOptions {
				c.Status(http.StatusOK)
			}
		}
	}

	c.Next()
}
