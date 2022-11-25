package xin

import (
	"net/http"

	"github.com/pandafw/pango/str"
)

// RedirectToSlach is a HandlerFunc that redirect to the URL's path + "/"
// ex: /index?page=1  --> /index/?page=1
func RedirectToSlach(c *Context) {
	if str.EndsWithByte(c.Request.URL.Path, '/') {
		return
	}

	u := *c.Request.URL
	u.Path += "/"
	c.Redirect(http.StatusFound, u.String())
}

// Redirector is a HandlerFunc that redirect to the url with http status codes[0] or http.StatusFound
func Redirector(url string, codes ...int) HandlerFunc {
	code := http.StatusFound
	if len(codes) > 0 {
		code = codes[0]
	}
	return func(c *Context) {
		c.Redirect(code, url)
	}
}
