package middleware

import (
	"net/http"
	"net/url"

	"github.com/askasoft/pango/xin"
)

const (
	AuthUserKey     = "X_AUTH_USER" // Key for authenticated user object saved in context
	AuthRedirectURL = "/login"      // default redirect URL
	AuthOriginQuery = "origin"      // default redirect origin URL query param
)

// AuthUser a authenticated user interface
type AuthUser interface {
	GetUsername() string
	GetPassword() string
}

// FindUserFunc user lookup function
type FindUserFunc func(c *xin.Context, username, password string) (AuthUser, error)

func AuthFailedRedirector(redirect string, query string) xin.HandlerFunc {
	return func(c *xin.Context) {
		if url := BuildRedirectURL(c, redirect, query); url != "" {
			c.Redirect(http.StatusTemporaryRedirect, url)
			c.Abort()
			return
		}

		c.AbortWithStatus(http.StatusUnauthorized)
	}
}

func BuildRedirectURL(c *xin.Context, redirect, query string) string {
	if redirect == "" || redirect == c.Request.URL.Path { // prevent redirect dead loop
		return ""
	}

	if xin.IsAjax(c) {
		return ""
	}

	url, err := url.Parse(redirect)
	if err != nil {
		c.Logger.Errorf("Invalid Redirect URL %q", redirect)
	} else {
		q := url.Query()
		q.Set(query, c.Request.URL.String())
		url.RawQuery = q.Encode()
		redirect = url.String()
	}

	return redirect
}
