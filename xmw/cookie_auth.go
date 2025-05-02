package xmw

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/askasoft/pango/cpt"
	"github.com/askasoft/pango/cpt/ccpt"
	"github.com/askasoft/pango/num"
	"github.com/askasoft/pango/str"
	"github.com/askasoft/pango/xin"
)

const (
	AuthCookieName             = "X_AUTH"
	AuthRedirectOriginURLQuery = "origin"
)

// CookieAuth cookie authenticator
type CookieAuth struct {
	Cryptor        cpt.Cryptor // cryptor to encode/decode cookie, MUST concurrent safe
	FindUser       FindUserFunc
	CookieName     string
	CookieMaxAge   time.Duration
	CookieDomain   string
	CookiePath     string
	CookieSecure   bool
	CookieHttpOnly bool
	CookieSameSite http.SameSite

	AuthUserKey     string
	RedirectURL     string
	OriginURLQuery  string
	AuthPassed      func(c *xin.Context, au AuthUser)
	AuthFailed      xin.HandlerFunc
	AuthRequired    xin.HandlerFunc
	GetCookieMaxAge func(c *xin.Context) time.Duration
}

func NewCookieAuth(f FindUserFunc, secret string) *CookieAuth {
	ca := &CookieAuth{
		Cryptor:        ccpt.NewAes128CBCCryptor(secret),
		FindUser:       f,
		CookieName:     AuthCookieName,
		CookiePath:     "/",
		CookieMaxAge:   time.Minute * 30,
		CookieSecure:   true,
		CookieHttpOnly: true,
		CookieSameSite: http.SameSiteLaxMode,
		AuthUserKey:    AuthUserKey,
		RedirectURL:    "/",
		OriginURLQuery: AuthRedirectOriginURLQuery,
	}
	ca.AuthPassed = ca.Authorized
	ca.AuthFailed = ca.Unauthorized
	ca.AuthRequired = ca.Unauthorized
	ca.GetCookieMaxAge = ca.getCookieMaxAge

	return ca
}

// Handle process xin request
func (ca *CookieAuth) Handle(c *xin.Context) {
	if _, ok := c.Get(ca.AuthUserKey); ok {
		// already authenticated
		c.Next()
		return
	}

	username, password, ok := ca.GetUserPassFromCookie(c)
	if !ok {
		ca.AuthRequired(c)
		return
	}

	au, err := ca.FindUser(c, username)
	if err != nil {
		c.Logger.Errorf("CookieAuth: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if au == nil || password != au.GetPassword() {
		ca.AuthFailed(c)
		return
	}

	ca.AuthPassed(c, au)
}

// Authorized set user to context and cookie then call c.Next()
func (ca *CookieAuth) Authorized(c *xin.Context, au AuthUser) {
	// set user to context
	c.Set(ca.AuthUserKey, au)

	// save or refresh cookie
	err := ca.SaveUserPassToCookie(c, au)
	if err != nil {
		c.Logger.Errorf("CookieAuth: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Next()
}

// Unauthorized redirect or abort with status 401
func (ca *CookieAuth) Unauthorized(c *xin.Context) {
	u := ca.BuildRedirectURL(c)
	if u != "" {
		c.Abort()
		c.Redirect(http.StatusTemporaryRedirect, u)
		return
	}

	c.AbortWithStatus(http.StatusUnauthorized)
}

func (ca *CookieAuth) BuildRedirectURL(c *xin.Context) string {
	u := ca.RedirectURL
	if u == "" || u == c.Request.URL.Path { // prevent redirect dead loop
		return ""
	}

	if xin.IsAjax(c) {
		return ""
	}

	p := ca.OriginURLQuery
	if p != "" {
		url, err := url.Parse(u)
		if err != nil {
			c.Logger.Errorf("Invalid RedirectURL %q", u)
		} else {
			q := url.Query()
			q.Set(p, c.Request.URL.String())
			url.RawQuery = q.Encode()
			u = url.String()
		}
	}

	return u
}

func (ca *CookieAuth) GetUserPassFromCookie(c *xin.Context) (username, password string, ok bool) {
	if raw, err := c.Cookie(ca.CookieName); err == nil && raw != "" {
		auth, err := ca.Cryptor.DecryptString(raw)
		if err != nil {
			c.Logger.Warnf("Invalid Cookie Auth %q: %v", raw, err)
			return
		}

		username, password, err = ca.decode(auth, ca.GetCookieMaxAge(c))
		if err != nil {
			c.Logger.Warnf("Invalid Cookie Auth %q: %v", auth, err)
			return
		}

		ok = true
	}
	return
}

func (ca *CookieAuth) SaveUserPassToCookie(c *xin.Context, au AuthUser) error {
	user, pass, age := au.GetUsername(), au.GetPassword(), ca.GetCookieMaxAge(c)

	val, err := ca.encrypt(user, pass, age)
	if err != nil {
		return err
	}

	ck := &http.Cookie{
		Name:     ca.CookieName,
		Value:    val,
		MaxAge:   int(age.Seconds()),
		Path:     ca.CookiePath,
		Domain:   ca.CookieDomain,
		Secure:   ca.CookieSecure,
		HttpOnly: ca.CookieHttpOnly,
		SameSite: ca.CookieSameSite,
	}

	c.SetCookie(ck)
	return nil
}

func (ca *CookieAuth) getCookieMaxAge(c *xin.Context) time.Duration {
	return ca.CookieMaxAge
}

func (ca *CookieAuth) DeleteCookie(c *xin.Context) {
	c.SetCookie(&http.Cookie{
		Name:     ca.CookieName,
		Value:    "",
		Expires:  time.Unix(1, 0),
		Path:     ca.CookiePath,
		Domain:   ca.CookieDomain,
		Secure:   ca.CookieSecure,
		HttpOnly: ca.CookieHttpOnly,
		SameSite: ca.CookieSameSite,
	})
}

func (ca *CookieAuth) encrypt(username, password string, maxage time.Duration) (string, error) {
	auth := ca.encode(username, password, maxage)
	return ca.Cryptor.EncryptString(auth)
}

func (ca *CookieAuth) encode(username, password string, maxage time.Duration) string {
	now := num.Ltoa(time.Now().UnixMilli())
	raw := fmt.Sprintf("%d\n%s\n%s", maxage.Milliseconds(), username, password)
	unsalt := base64.RawURLEncoding.EncodeToString(str.UnsafeBytes(raw))
	salted := cpt.Salt(cpt.SecretChars, now, unsalt)
	auth := fmt.Sprintf("%s\n%s", now, salted)
	return auth
}

func (ca *CookieAuth) decode(auth string, maxage time.Duration) (string, string, error) {
	timestamp, salted, ok := str.CutByte(auth, '\n')
	if !ok {
		return "", "", errors.New("no timestamp")
	}

	unsalt := cpt.Unsalt(cpt.SecretChars, timestamp, salted)
	bs, err := base64.RawURLEncoding.DecodeString(unsalt)
	if err != nil {
		return "", "", err
	}

	raw := str.UnsafeString(bs)

	ss := str.FieldsByte(raw, '\n')
	if len(ss) != 3 {
		return "", "", errors.New("invalid - " + raw)
	}

	duration := num.Atol(ss[0])

	// cookie maxage check
	if maxage.Milliseconds() != duration {
		return "", "", fmt.Errorf("cookie max age %d != %d", maxage.Milliseconds(), duration)
	}

	now := time.Now().UnixMilli()
	delta := time.Minute.Milliseconds()

	// -+ 1m for different time on cluster servers
	start, after := now-duration-delta, now+delta

	// timestamp check
	created := num.Atol(timestamp)
	if created < start || created > after {
		return "", "", errors.New("timestamp expired")
	}

	return ss[1], ss[2], nil
}
