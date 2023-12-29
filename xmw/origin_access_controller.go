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
}

// NewOriginAccessController create a default OriginAccessController
func NewOriginAccessController(origins ...string) *OriginAccessController {
	return &OriginAccessController{Origins: cog.NewHashSet(origins...)}
}

// Handler returns the xin.HandlerFunc
func (ll *OriginAccessController) Handler() xin.HandlerFunc {
	return ll.Handle
}

// SetOrigins set allowed origins
func (ll *OriginAccessController) SetOrigins(origins ...string) {
	ll.Origins = cog.NewHashSet(origins...)
}

// Handle process xin request
func (ll *OriginAccessController) Handle(c *xin.Context) {
	if ll.Origins.Len() > 0 {
		acao := ""

		if ll.Origins.Contain("*") {
			acao = "*"
		} else {
			origin := c.GetHeader("Origin")
			if origin != "" && ll.Origins.Contain(origin) {
				acao = origin
			}
		}

		if acao != "" {
			c.Header("Access-Control-Allow-Origin", "*")
			if c.Request.Method == http.MethodOptions {
				c.Status(http.StatusOK)
			}
		}
	}

	c.Next()
}
