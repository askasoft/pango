package xmw

import (
	"net/http"
	"strconv"

	"github.com/askasoft/pango/str"
	"github.com/askasoft/pango/xin"
)

const (
	BearerAuthPrefix = "Bearer " // Bearer Authentication Prefix
)

// BearerAuth bearer http authenticator
type BearerAuth struct {
	Realm       string
	FindUser    FindUserFunc
	AuthUserKey string
	AuthPassed  func(c *xin.Context, au AuthUser)
	AuthFailed  xin.HandlerFunc
}

func NewBearerAuth(f FindUserFunc) *BearerAuth {
	ba := &BearerAuth{
		AuthUserKey: AuthUserKey,
		FindUser:    f,
	}
	ba.AuthPassed = ba.authorized
	ba.AuthFailed = ba.Unauthorized

	return ba
}

// Handle process xin request
func (ba *BearerAuth) Handle(c *xin.Context) {
	next, au, err := ba.Authenticate(c)
	if err != nil {
		c.Logger.Errorf("BearerAuth: %v", err)
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

func (ba *BearerAuth) authorized(c *xin.Context, au AuthUser) {
	c.Next()
}

// Unauthorized set bearer authentication WWW-Authenticate header
func (ba *BearerAuth) Unauthorized(c *xin.Context) {
	c.Header("WWW-Authenticate", "Bearer realm="+strconv.Quote(str.IfEmpty(ba.Realm, "Authorization Required")))
	c.AbortWithStatus(http.StatusUnauthorized)
}

func (ba *BearerAuth) Authenticate(c *xin.Context) (next bool, au AuthUser, err error) {
	if _, ok := c.Get(ba.AuthUserKey); ok {
		// already authenticated
		next = true
		return
	}

	token := ba.bearerRequestAuthh(c)
	if token == "" {
		return
	}

	au, err = ba.FindUser(c, token, "")
	if err != nil || au == nil {
		return
	}

	// set user to context
	c.Set(ba.AuthUserKey, au)

	return
}

func (ba *BearerAuth) bearerRequestAuthh(c *xin.Context) string {
	auth := c.Request.Header.Get("Authorization")
	if auth == "" {
		return ""
	}

	if !str.StartsWithFold(auth, BearerAuthPrefix) {
		return ""
	}

	return auth[len(BearerAuthPrefix):]
}
