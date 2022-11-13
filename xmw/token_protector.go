package xmw

import (
	"net/http"
	"time"

	"github.com/pandafw/pango/col"
	"github.com/pandafw/pango/cpt"
	"github.com/pandafw/pango/xin"
)

// TokenProtector token protector for CSRF
type TokenProtector struct {
	Cryptor        cpt.Cryptor
	Expires        time.Duration
	Methods        col.HashSet
	AttrKey        string
	HeaderKey      string
	ParamName      string
	CookieName     string
	CookieMaxAge   time.Duration
	CookieDomain   string
	CookiePath     string
	CookieSecure   bool
	CookieHttpOnly bool
}

// NewTokenProtector create a default TokenProtector
func NewTokenProtector(secret string) *TokenProtector {
	t := &TokenProtector{
		Cryptor:        cpt.NewAesCBC(secret),
		Expires:        time.Hour * 24,
		AttrKey:        "WW_TOKEN",
		HeaderKey:      "X-WW-TOKEN",
		ParamName:      "_token_",
		CookieName:     "WW_TOKEN",
		CookieMaxAge:   time.Hour * 24 * 30,
		CookieHttpOnly: true,
	}
	t.Methods.Add(http.MethodDelete, http.MethodPatch, http.MethodPost, http.MethodPost)
	return t
}

// Handler returns the xin.HandlerFunc
func (tp *TokenProtector) Handler() xin.HandlerFunc {
	return func(c *xin.Context) {
		tp.handle(c)
	}
}

// handle process xin request
func (tp *TokenProtector) handle(c *xin.Context) {
	if tp.Methods.Contains(c.Request.Method) {
		if !tp.validate(c) {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
	}

	c.Next()
}

func (tp *TokenProtector) validate(c *xin.Context) bool {
	st := tp.getSourceToken(c)
	if st == nil {
		return false
	}

	rt := tp.getRequestToken(c)
	if rt == nil {
		return false
	}

	if tp.Expires > 0 && rt.Timestamp.Add(tp.Expires).Before(time.Now()) {
		return false
	}

	return st.Secret == rt.Secret
}

func (tp *TokenProtector) parseToken(ts string) (*cpt.Token, error) {
	s, err := tp.Cryptor.DecryptString(ts)
	if err != nil {
		return nil, err
	}

	t, err := cpt.ParseToken(s)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (tp *TokenProtector) getSourceToken(c *xin.Context) *cpt.Token {
	if av, ok := c.Get(tp.AttrKey); ok {
		if t, ok := av.(*cpt.Token); ok {
			return t
		}
		c.Logger().Errorf("Invalid Context Token: %v", av)
	}

	ts, err := c.Cookie(tp.CookieName)
	if err != nil {
		return nil
	}

	t, err := tp.parseToken(ts)
	if err != nil {
		c.Logger().Warnf("Invalid Cookie Token: %v: %q", err, ts)
		return nil
	}

	c.Set(tp.AttrKey, t)
	return t
}

func (tp *TokenProtector) getRequestToken(c *xin.Context) *cpt.Token {
	if ts := c.PostForm(tp.ParamName); ts != "" {
		t, err := tp.parseToken(ts)
		if err == nil {
			return t
		}
		c.Logger().Warnf("Invalid Form Token: %v: %q", err, ts)
	}

	if ts := c.Query(tp.ParamName); ts != "" {
		t, err := tp.parseToken(ts)
		if err == nil {
			return t
		}
		c.Logger().Warnf("Invalid Query Token: %v: %q", err, ts)
	}

	if ts := c.GetHeader(tp.HeaderKey); ts != "" {
		t, err := tp.parseToken(ts)
		if err == nil {
			return t
		}
		c.Logger().Warnf("Invalid Header Token: %v: %q", err, ts)
	}
	return nil
}

func (tp *TokenProtector) RefreshToken(c *xin.Context) string {
	t := tp.getSourceToken(c)
	if t == nil {
		t = cpt.NewToken()
		c.Set(tp.AttrKey, t)
	} else {
		t.Refresh()
	}

	ts, err := tp.Cryptor.EncryptString(t.Token)
	if err == nil {
		c.SetCookie(tp.CookieName, ts, int(tp.CookieMaxAge.Seconds()), tp.CookiePath, tp.CookieDomain, tp.CookieSecure, tp.CookieHttpOnly)
	} else {
		c.Logger().Errorf("EncryptToken: %v", err)
	}
	return ts
}
