package xmw

import (
	"crypto/md5" //nolint: gosec
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/askasoft/pango/cpt"
	"github.com/askasoft/pango/net/httpx"
	"github.com/askasoft/pango/str"
	"github.com/askasoft/pango/xin"
)

const (
	DigestAuthCtxKey = "X_DIGEST" // Key for digest parameters saved in context
	DigestAuthPrefix = "Digest "  // Digest Authentication Prefix
)

var (
	errDigestParamMissing = errors.New("da: digest parameter missing")
	errDigestParamInvalid = errors.New("da: digest parameter invalid")
)

// DigestAuth digest http authenticator
type DigestAuth struct {
	Realm        string
	Opaque       string
	NonceExpires time.Duration
	FindUser     FindUserFunc
	AuthUserKey  string
	AuthPassed   func(c *xin.Context, au AuthUser)
	AuthFailed   xin.HandlerFunc

	noncer *cpt.Tokener
}

func NewDigestAuth(f FindUserFunc) *DigestAuth {
	da := &DigestAuth{
		AuthUserKey:  AuthUserKey,
		NonceExpires: time.Minute * 5,
		FindUser:     f,
	}
	da.AuthPassed = da.authorized
	da.AuthFailed = da.Unauthorized
	da.noncer = cpt.NewTokener(8, 16)

	return da
}

// Handle process xin request
func (da *DigestAuth) Handle(c *xin.Context) {
	next, au, err := da.Authenticate(c)
	if err != nil {
		c.Logger.Errorf("DigestAuth: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if next {
		// already authenticated
		c.Next()
		return
	}

	if au == nil {
		da.AuthFailed(c)
		return
	}

	da.AuthPassed(c, au)
}

func (da *DigestAuth) authorized(c *xin.Context, au AuthUser) {
	c.Next()
}

// Unauthorized set WWW-Authenticate header
func (da *DigestAuth) Unauthorized(c *xin.Context) {
	cip := base64.RawURLEncoding.EncodeToString(str.UnsafeBytes(c.ClientIP()))
	nonce := da.noncer.SecretToken(cip)
	wa := fmt.Sprintf(`Digest realm="%s", nonce="%s", opaque="%s", algorithm=MD5, qop="auth"`, da.Realm, nonce.Token(), da.Opaque)

	c.Header("WWW-Authenticate", wa)
	c.AbortWithStatus(http.StatusUnauthorized)
}

func (da *DigestAuth) Authenticate(c *xin.Context) (next bool, au AuthUser, err error) {
	if _, ok := c.Get(da.AuthUserKey); ok {
		// already authenticated
		next = true
		return
	}

	dap := da.digestRequestAuth(c)
	if dap == nil {
		da.AuthFailed(c)
		return
	}

	c.Set(DigestAuthCtxKey, dap)

	au, err = da.FindUser(c, dap["username"], dap["response"])
	if err != nil || au == nil {
		return
	}

	// set user to context
	c.Set(da.AuthUserKey, au)

	return
}

// digestRequestAuth parses Authorization header from the
// http.Request. Returns a map of auth parameters or nil if the header
// is not a valid parsable Digest auth header.
func (da *DigestAuth) digestRequestAuth(c *xin.Context) map[string]string {
	auth := c.Request.Header.Get("Authorization")
	if auth == "" {
		return nil
	}

	if !str.StartsWithFold(auth, DigestAuthPrefix) {
		return nil
	}

	dap := httpx.ParsePairs(auth[len(DigestAuthPrefix):])
	if len(dap) == 0 {
		return nil
	}

	if !da.checkAlgorithm(c, dap) || !da.checkURI(c, dap["uri"]) || !da.checkNonce(c, dap["nonce"]) {
		return nil
	}

	return dap
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
	t, err := da.noncer.ParseToken(nonce)
	if err != nil {
		c.Logger.Debugf("DigestAuth nonce %q is invalid", nonce)
		return false
	}
	bs, err := base64.RawURLEncoding.DecodeString(t.Secret())
	if err != nil {
		c.Logger.Debugf("DigestAuth nonce %q contains invalid secret", t.String())
		return false
	}
	if str.UnsafeString(bs) != c.ClientIP() {
		c.Logger.Debugf("DigestAuth %q does not match request client IP %s", t.String(), c.ClientIP())
		return false
	}
	if da.NonceExpires.Seconds() > 0 && t.Timestamp().Add(da.NonceExpires).Before(time.Now()) {
		c.Logger.Debugf("DigestAuth nonce %q is expired", t.String())
		return false
	}
	return true
}

func (da *DigestAuth) VerifyPassword(c *xin.Context, password string) (bool, error) {
	var dap map[string]string

	val, ok := c.Get(DigestAuthCtxKey)
	if !ok {
		return false, errDigestParamMissing
	}
	if dap, ok = val.(map[string]string); !ok {
		return false, errDigestParamInvalid
	}

	digest := md5.New() //nolint: gosec

	// hash password
	ha1 := dap["username"] + ":" + da.Realm + ":" + password
	digest.Reset()
	digest.Write(str.UnsafeBytes(ha1))
	ha1 = fmt.Sprintf("%x", digest.Sum(nil))

	ha2 := c.Request.Method + ":" + dap["uri"]
	digest.Reset()
	digest.Write(str.UnsafeBytes(ha2))
	ha2 = fmt.Sprintf("%x", digest.Sum(nil))

	kd := strings.Join([]string{ha1, dap["nonce"], dap["nc"], dap["cnonce"], dap["qop"], ha2}, ":")
	digest.Reset()
	digest.Write(str.UnsafeBytes(kd))
	kd = fmt.Sprintf("%x", digest.Sum(nil))

	// verify password
	if kd == dap["response"] {
		return true, nil
	}

	c.Logger.Debugf("DigestAuth response %q is invalid", dap["response"])
	return false, nil
}
