package xmw

import (
	"net/http"
	"strconv"

	"github.com/askasoft/pango/str"
	"github.com/askasoft/pango/xin"
)

// AuthUserKey is the key for user credential authenticated saved in context
const AuthUserKey = "XMW_USER"

type User interface {
	GetUsername() string
	GetPassword() string
}

type UserProvider interface {
	FindUser(username string) User
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
	return ba.Handle
}

// Handle process xin request
func (ba *BasicAuth) Handle(c *xin.Context) {
	username, password, ok := c.Request.BasicAuth()
	if ok {
		if user := ba.UserProvider.FindUser(username); user != nil {
			if password == user.GetPassword() {
				c.Set(ba.AuthUserKey, user)
				c.Next()
				return
			}
		}
	}

	c.Header("WWW-Authenticate", "Basic realm="+strconv.Quote(str.IfEmpty(ba.Realm, "Authorization Required")))
	c.AbortWithStatus(http.StatusUnauthorized)
}
