package xmw

import (
	"net/http"

	"github.com/askasoft/pango/xin"
)

// AuthUserKey is the key for user credential authenticated saved in context
const AuthUserKey = "WW_USER"

type UserProvider interface {
	FindUser(username string) any
	GetPassword(user any) string
}

// BasicAuth basic http authenticator
type BasicAuth struct {
	UserProvider UserProvider
	AuthUserKey  string
	Realm        string
}

func NewBasicAuth(up UserProvider) *BasicAuth {
	return &BasicAuth{
		UserProvider: up,
		AuthUserKey:  AuthUserKey,
	}
}

// Handler returns the xin.HandlerFunc
func (ba *BasicAuth) Handler() xin.HandlerFunc {
	return func(c *xin.Context) {
		ba.handle(c)
	}
}

// handle process xin request
func (ba *BasicAuth) handle(c *xin.Context) {
	username, password, ok := c.Request.BasicAuth()
	if ok {
		if user := ba.UserProvider.FindUser(username); user != nil {
			if password == ba.UserProvider.GetPassword(user) {
				c.Set(ba.AuthUserKey, user)
				c.Next()
				return
			}
		}
	}

	c.Header("WWW-Authenticate", "Basic "+ba.Realm)
	c.AbortWithStatus(http.StatusUnauthorized)
}
