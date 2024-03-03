package xmw

import (
	"crypto/md5" //nolint: gosec
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/askasoft/pango/bye"
	"github.com/askasoft/pango/cpt"
	"github.com/askasoft/pango/net/httpx"
	"github.com/askasoft/pango/str"
	"github.com/askasoft/pango/xin"
)

// DigestAuth digest http authenticator
type DigestAuth struct {
	UserProvider UserProvider
	AuthUserKey  string
	Realm        string
	Opaque       string
	NonceExpires time.Duration
}

func NewDigestAuth(up UserProvider) *DigestAuth {
	return &DigestAuth{
		UserProvider: up,
		AuthUserKey:  AuthUserKey,
		NonceExpires: time.Minute * 5,
	}
}

// Handler returns the xin.HandlerFunc
func (da *DigestAuth) Handler() xin.HandlerFunc {
	return da.Handle
}

// Handle process xin request
func (da *DigestAuth) Handle(c *xin.Context) {
	auth := c.Request.Header.Get("Authorization")
	if auth != "" {
		dap := DigestAuthParams(auth)
		if dap != nil {
			if user, ok := da.checkAuthorization(c, dap); ok {
				c.Set(AuthUserKey, user)
				c.Next()
				return
			}
		}
	}

	cip := base64.RawURLEncoding.EncodeToString(str.UnsafeBytes(c.ClientIP()))
	nonce := noncer.NewToken(cip)
	wa := fmt.Sprintf(`Digest realm="%s", nonce="%s", opaque="%s", algorithm=MD5, qop="auth"`, da.Realm, nonce.Token(), da.Opaque)

	c.Header("WWW-Authenticate", wa)
	c.AbortWithStatus(http.StatusUnauthorized)
}

var noncer = cpt.NewTokener(8, 16)

func (da *DigestAuth) checkAuthorization(c *xin.Context, auth map[string]string) (any, bool) {
	if !da.checkAlgorithm(c, auth) {
		return nil, false
	}
	if !da.checkURI(c, auth["uri"]) {
		return nil, false
	}
	if !da.checkNonce(c, auth["nonce"]) {
		return nil, false
	}
	return da.checkUserPass(c, auth)
}

func (da *DigestAuth) checkAlgorithm(c *xin.Context, auth map[string]string) bool {
	// RFC2617 Section 3.2.1 specifies that unset value of algorithm in
	// WWW-Authenticate Response header should be treated as
	// "MD5". According to section 3.2.2 the "algorithm" value in
	// subsequent Request Authorization header must be set to whatever
	// was supplied in the WWW-Authenticate Response header. This
	// implementation always returns an algorithm in WWW-Authenticate
	// header, however there seems to be broken clients in the wild
	// which do not set the algorithm. Assume the unset algorithm in
	// Authorization header to be equal to MD5.
	if _, ok := auth["algorithm"]; !ok {
		auth["algorithm"] = "MD5"
	}
	if da.Opaque != auth["opaque"] {
		c.Logger.Debugf("Digest auth opaque %q is invalid", auth["opaque"])
		return false
	}
	if auth["algorithm"] != "MD5" {
		c.Logger.Debugf("Digest auth algorithm %q is invalid", auth["algorithm"])
		return false
	}
	if auth["qop"] != "auth" {
		c.Logger.Debugf("Digest auth qop %q is invalid", auth["algorithm"])
		return false
	}
	return true
}

func (da *DigestAuth) checkURI(c *xin.Context, uri string) bool {
	// Check if the requested URI matches auth header
	if c.Request.RequestURI != uri {
		// We allow auth["uri"] to be a full path prefix of request-uri
		// for some reason lost in history, which is probably wrong, but
		// used to be like that for quite some time
		// (https://tools.ietf.org/html/rfc2617#section-3.2.2 explicitly
		// says that auth["uri"] is the request-uri).
		switch u, err := url.Parse(uri); {
		case err != nil:
			return false
		case c.Request.URL == nil:
			return false
		case len(u.Path) > len(c.Request.URL.Path):
			return false
		case !strings.HasPrefix(c.Request.URL.Path, u.Path):
			return false
		}
	}

	return true
}

func (da *DigestAuth) checkNonce(c *xin.Context, nonce string) bool {
	t, err := noncer.ParseToken(nonce)
	if err != nil {
		c.Logger.Debugf("Digest auth nonce %q is invalid", nonce)
		return false
	}
	bs, err := base64.RawURLEncoding.DecodeString(t.Secret())
	if err != nil {
		c.Logger.Debugf("Digest auth nonce %q contains invalid secret", t.String())
		return false
	}
	if bye.UnsafeString(bs) != c.ClientIP() {
		c.Logger.Debugf("Digest auth %q does not match request client IP %s", t.String(), c.ClientIP())
		return false
	}
	if da.NonceExpires.Seconds() > 0 && t.Timestamp().Add(da.NonceExpires).Before(time.Now()) {
		c.Logger.Debugf("Digest auth nonce %q is expired", t.String())
		return false
	}
	return true
}

func (da *DigestAuth) checkUserPass(c *xin.Context, auth map[string]string) (any, bool) {
	digest := md5.New() //nolint: gosec

	// find user
	user := da.UserProvider.FindUser(auth["username"])
	if user == nil {
		return nil, false
	}
	pass := user.GetPassword()

	// hash password
	ha1 := auth["username"] + ":" + da.Realm + ":" + pass

	digest.Reset()
	digest.Write(str.UnsafeBytes(ha1))
	ha1 = fmt.Sprintf("%x", digest.Sum(nil))

	ha2 := c.Request.Method + ":" + auth["uri"]
	digest.Reset()
	digest.Write(str.UnsafeBytes(ha2))
	ha2 = fmt.Sprintf("%x", digest.Sum(nil))

	kd := strings.Join([]string{ha1, auth["nonce"], auth["nc"], auth["cnonce"], auth["qop"], ha2}, ":")
	digest.Reset()
	digest.Write(str.UnsafeBytes(kd))
	kd = fmt.Sprintf("%x", digest.Sum(nil))

	// check password
	if kd != auth["response"] {
		c.Logger.Debugf("Digest auth response %q is invalid", auth["response"])
		return nil, false
	}

	return user, true
}

// DigestAuthParams parses Authorization header from the
// http.Request. Returns a map of auth parameters or nil if the header
// is not a valid parsable Digest auth header.
func DigestAuthParams(auth string) map[string]string {
	const prefix = "Digest "
	if !str.StartsWithFold(auth, prefix) {
		return nil
	}

	return httpx.ParsePairs(auth[len(prefix):])
}
