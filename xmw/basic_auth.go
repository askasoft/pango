package xmw

import (
	"net/http"
	"strconv"

	"github.com/askasoft/pango/str"
	"github.com/askasoft/pango/xin"
)

const (
	BasicAuthPrefix = "Basic " // Basic Authentication Prefix
)

// BasicAuth basic http authenticator
type BasicAuth struct {
	Realm       string
	FindUser    FindUserFunc
	AuthUserKey string
	AuthPassed  func(c *xin.Context, au AuthUser)
	AuthFailed  xin.HandlerFunc
}

func NewBasicAuth(f FindUserFunc) *BasicAuth {
	ba := &BasicAuth{
		AuthUserKey: AuthUserKey,
		FindUser:    f,
	}
	ba.AuthPassed = ba.authorized
	ba.AuthFailed = ba.Unauthorized

	return ba
}

// Handle process xin request
func (ba *BasicAuth) Handle(c *xin.Context) {
	next, au, err := ba.Authenticate(c)
	if err != nil {
		c.Logger.Errorf("BasicAuth: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if next {
		// already authenticated
		c.Next()
		return
	}

	if au == nil {
		ba.AuthFailed(c)
		return
	}

	ba.AuthPassed(c, au)
}

func (ba *BasicAuth) authorized(c *xin.Context, au AuthUser) {
	c.Next()
}

// Unauthorized set basic authentication WWW-Authenticate header
func (ba *BasicAuth) Unauthorized(c *xin.Context) {
	c.Header("WWW-Authenticate", "Basic realm="+strconv.Quote(str.IfEmpty(ba.Realm, "Authorization Required")))
	c.AbortWithStatus(http.StatusUnauthorized)
}

func (ba *BasicAuth) Authenticate(c *xin.Context) (next bool, au AuthUser, err error) {
	if _, ok := c.Get(ba.AuthUserKey); ok {
		// already authenticated
		next = true
		return
	}

	username, password, ok := c.Request.BasicAuth()
	if !ok {
		return
	}

	au, err = ba.FindUser(c, username, password)
	if err != nil || au == nil {
		return
	}

	// set user to context
	c.Set(ba.AuthUserKey, au)

	return
}
