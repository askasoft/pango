package xmw

import (
	"net/http"
	"strconv"

	"github.com/askasoft/pango/str"
	"github.com/askasoft/pango/xin"
)

// AuthUserKey is the key for user credential authenticated saved in context
const AuthUserKey = "XMW_USER"

type AuthUser interface {
	GetUsername() string
	GetPassword() string
}

type FindUserFunc func(c *xin.Context, username string) (AuthUser, error)

// BasicAuth basic http authenticator
type BasicAuth struct {
	Realm       string
	FindUser    FindUserFunc
	AuthUserKey string
}

func NewBasicAuth(f FindUserFunc) *BasicAuth {
	return &BasicAuth{
		AuthUserKey: AuthUserKey,
		FindUser:    f,
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
		user, err := ba.FindUser(c, username)
		if err != nil {
			c.Logger.Errorf("BasicAuth: %v", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if user != nil && password == user.GetPassword() {
			c.Set(ba.AuthUserKey, user)
			c.Next()
			return
		}
	}

	c.Header("WWW-Authenticate", "Basic realm="+strconv.Quote(str.IfEmpty(ba.Realm, "Authorization Required")))
	c.AbortWithStatus(http.StatusUnauthorized)
}
