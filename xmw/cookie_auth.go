package xmw

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/askasoft/pango/bye"
	"github.com/askasoft/pango/cpt"
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
	Cryptor        cpt.Cryptor
	FindUser       FindUserFunc
	CookieName     string
	CookieMaxAge   time.Duration
	CookieDomain   string
	CookiePath     string
	CookieSecure   bool
	CookieHttpOnly bool
	CookieSameSite http.SameSite
	AuthUserKey    string
	RedirectURL    string
	OriginURLQuery string
}

func NewCookieAuth(f FindUserFunc, secret string) *CookieAuth {
	return &CookieAuth{
		Cryptor:        cpt.NewAesCBC(secret),
		FindUser:       f,
		CookieName:     AuthCookieName,
		CookieMaxAge:   time.Minute * 30,
		CookieSecure:   true,
		CookieHttpOnly: true,
		CookieSameSite: http.SameSiteStrictMode,
		AuthUserKey:    AuthUserKey,
		RedirectURL:    "/",
		OriginURLQuery: AuthRedirectOriginURLQuery,
	}
}

// Handler returns the xin.HandlerFunc
func (ca *CookieAuth) Handler() xin.HandlerFunc {
	return ca.Handle
}

// Handle process xin request
func (ca *CookieAuth) Handle(c *xin.Context) {
	username, password, ok := ca.GetUserPassFromCookie(c)
	if ok {
		user, err := ca.FindUser(c, username)
		if err != nil {
			c.Logger.Errorf("CookieAuth: %v", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if user != nil && password == user.GetPassword() {
			// save or refresh cookie
			err = ca.SaveUserPassToCookie(c, username, password)
			if err != nil {
				c.Logger.Errorf("CookieAuth: %v", err)
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}

			// set user to context
			c.Set(ca.AuthUserKey, user)

			c.Next()
			return
		}
	}

	u := ca.buildRedirectURL(c)
	if u != "" {
		c.Abort()
		c.Redirect(http.StatusTemporaryRedirect, u)
		return
	}

	c.AbortWithStatus(http.StatusUnauthorized)
}

func (ca *CookieAuth) buildRedirectURL(c *xin.Context) string {
	u := ca.RedirectURL
	if u == "" {
		return ""
	}

	if str.EqualFold(c.GetHeader("X-Requested-With"), "XMLHttpRequest") {
		return ""
	}

	p := ca.OriginURLQuery
	if p != "" {
		url, err := url.Parse(u)
		if err != nil {
			c.Logger.Errorf("Invalid RedirectURL %q", u)
		} else {
			q := url.Query()
			if q.Get(p) == "" {
				q.Add(p, c.Request.URL.String())
				url.RawQuery = q.Encode()
				u = url.String()
			}
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

		username, password, ok = ca.decode(auth)
		if !ok {
			c.Logger.Warnf("Invalid Cookie Auth %q", auth)
		}

	}
	return
}

func (ca *CookieAuth) SaveUserPassToCookie(c *xin.Context, username, password string) error {
	val, err := ca.encrypt(username, password)
	if err != nil {
		return err
	}

	c.SetCookie(&http.Cookie{
		Name:     ca.CookieName,
		Value:    val,
		MaxAge:   int(ca.CookieMaxAge.Seconds()),
		Path:     ca.CookiePath,
		Domain:   ca.CookieDomain,
		Secure:   ca.CookieSecure,
		HttpOnly: ca.CookieHttpOnly,
		SameSite: ca.CookieSameSite,
	})
	return nil
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

func (ca *CookieAuth) encrypt(username, password string) (string, error) {
	auth := ca.encode(username, password)
	return ca.Cryptor.EncryptString(auth)
}

func (ca *CookieAuth) encode(username, password string) string {
	now := num.Ltoa(time.Now().UnixMilli())
	raw := fmt.Sprintf("%d\n%s\n%s", ca.CookieMaxAge.Milliseconds(), username, password)
	unsalt := base64.RawURLEncoding.EncodeToString(str.UnsafeBytes(raw))
	salted := cpt.Salt(cpt.SecretChars, now, unsalt)
	auth := fmt.Sprintf("%s\n%s", now, salted)
	return auth
}

func (ca *CookieAuth) decode(auth string) (username, password string, ok bool) {
	timestamp, salted, ok := str.CutByte(auth, '\n')
	if !ok {
		return
	}

	unsalt := cpt.Unsalt(cpt.SecretChars, timestamp, salted)
	bs, err := base64.RawURLEncoding.DecodeString(unsalt)
	if err != nil {
		return
	}

	raw := bye.UnsafeString(bs)

	ss := str.FieldsByte(raw, '\n')
	if len(ss) != 3 {
		return
	}

	duration := num.Atol(ss[0])

	// cookie maxage check
	if ca.CookieMaxAge.Milliseconds() != duration {
		return
	}

	now := time.Now().UnixMilli()
	delta := time.Minute.Milliseconds()

	// -+ 1m for different time on cluster servers
	start, after := now-duration-delta, now+delta

	// timestamp check
	created := num.Atol(timestamp)
	if created >= start && created <= after {
		username, password, ok = ss[1], ss[2], true
	}
	return
}
