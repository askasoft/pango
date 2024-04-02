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
	Realm        string
	FindUser     FindUserFunc
	AuthUserKey  string
	AuthPassed   xin.HandlerFunc
	AuthFailed   xin.HandlerFunc
	AuthRequired xin.HandlerFunc
}

func NewBasicAuth(f FindUserFunc) *BasicAuth {
	ba := &BasicAuth{
		AuthUserKey: AuthUserKey,
		FindUser:    f,
	}
	ba.AuthPassed = ba.Authorized
	ba.AuthFailed = ba.Unauthorized
	ba.AuthRequired = ba.Unauthorized

	return ba
}

// Handler returns the xin.HandlerFunc
func (ba *BasicAuth) Handler() xin.HandlerFunc {
	return ba.Handle
}

// Handle process xin request
func (ba *BasicAuth) Handle(c *xin.Context) {
	if _, ok := c.Get(ba.AuthUserKey); ok {
		// already authenticated
		c.Next()
		return
	}

	username, password, ok := c.Request.BasicAuth()
	if !ok {
		ba.AuthRequired(c)
		return
	}

	user, err := ba.FindUser(c, username)
	if err != nil {
		c.Logger.Errorf("BasicAuth: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if user == nil || password != user.GetPassword() {
		ba.AuthFailed(c)
		return
	}

	c.Set(ba.AuthUserKey, user)
	ba.AuthPassed(c)
}

// Authorized just call c.Next()
func (ba *BasicAuth) Authorized(c *xin.Context) {
	c.Next()
}

// Unauthorized set basic authentication WWW-Authenticate header
func (ba *BasicAuth) Unauthorized(c *xin.Context) {
	c.Header("WWW-Authenticate", "Basic realm="+strconv.Quote(str.IfEmpty(ba.Realm, "Authorization Required")))
	c.AbortWithStatus(http.StatusUnauthorized)
}
