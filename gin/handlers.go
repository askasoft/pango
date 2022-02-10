package gin

import (
	"net/http"

	"github.com/pandafw/pango/str"
)

// RediretToSlach is a HandlerFunc that redirect to the URL's path + "/"
// ex: /index?page=1  --> /index/?page=1
func RediretToSlach(c *Context) {
	if str.EndsWithByte(c.Request.URL.Path, '/') {
		return
	}

	u := *c.Request.URL
	u.Path += "/"
	c.Redirect(http.StatusFound, u.String())
}
