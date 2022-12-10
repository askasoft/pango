package xmw

import (
	"net/http"
	"time"

	"github.com/pandafw/pango/cpt"
	"github.com/pandafw/pango/xin"
)

const (
	TokenAttrKey    = "WW_TOKEN"
	TokenParamName  = "_token_"
	TokenHeaderName = "X-WW-TOKEN" //nolint: gosec
	TokenCookieName = "WW_TOKEN"
)

// TokenProtector token protector for CSRF
type TokenProtector struct {
	Cryptor        cpt.Cryptor
	Expires        time.Duration
	AttrKey        string
	ParamName      string
	HeaderName     string
	CookieName     string
	CookieMaxAge   time.Duration
	CookieDomain   string
	CookiePath     string
	CookieSecure   bool
	CookieHttpOnly bool

	methods *stringSet
}

// NewTokenProtector create a default TokenProtector
// default methods: DELETE, PATCH, POST, PUT
func NewTokenProtector(secret string) *TokenProtector {
	t := &TokenProtector{
		Cryptor:        cpt.NewAesCBC(secret),
		Expires:        time.Hour * 24,
		AttrKey:        TokenAttrKey,
		ParamName:      TokenParamName,
		HeaderName:     TokenHeaderName,
		CookieName:     TokenCookieName,
		CookieMaxAge:   time.Hour * 24 * 30, // 30 days
		CookieHttpOnly: true,
		methods:        newStringSet(http.MethodDelete, http.MethodPatch, http.MethodPost, http.MethodPut),
	}
	return t
}

// SetMethods Set the http methods to protect
// default methods: DELETE, PATCH, POST, PUT
func (tp *TokenProtector) SetMethods(ms ...string) {
	if len(ms) == 0 {
		tp.methods = nil
		return
	}

	tp.methods = newStringSet(ms...)
}

// Handler returns the xin.HandlerFunc
func (tp *TokenProtector) Handler() xin.HandlerFunc {
	return func(c *xin.Context) {
		tp.handle(c)
	}
}

// handle process xin request
func (tp *TokenProtector) handle(c *xin.Context) {
	ms := tp.methods
	if ms != nil && ms.Contains(c.Request.Method) {
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
	c.Logger.Tracef("Source token: %v", st)

	rt := tp.getRequestToken(c)
	if rt == nil {
		return false
	}
	c.Logger.Tracef("Request token: %v", rt)

	if st.Secret() != rt.Secret() {
		c.Logger.Warnf("Invalid token secret %q, want %q", rt.Secret(), st.Secret())
		return false
	}

	if tp.Expires > 0 && rt.Timestamp().Add(tp.Expires).Before(time.Now()) {
		c.Logger.Warnf("Request token (%v) is expired for %v", rt, tp.Expires)
		return false
	}

	return true
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
		c.Logger.Errorf("Invalid Context Token: %v", av)
	}

	ts, err := c.Cookie(tp.CookieName)
	if err != nil {
		return nil
	}

	t, err := tp.parseToken(ts)
	if err != nil {
		c.Logger.Warnf("Invalid Cookie Token: %v: %q", err, ts)
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
		c.Logger.Warnf("Invalid Form Token: %v: %q", err, ts)
	}

	if ts := c.Query(tp.ParamName); ts != "" {
		t, err := tp.parseToken(ts)
		if err == nil {
			return t
		}
		c.Logger.Warnf("Invalid Query Token: %v: %q", err, ts)
	}

	if ts := c.GetHeader(tp.HeaderName); ts != "" {
		t, err := tp.parseToken(ts)
		if err == nil {
			return t
		}
		c.Logger.Warnf("Invalid Header Token: %v: %q", err, ts)
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

	ts, err := tp.Cryptor.EncryptString(t.Token())
	if err == nil {
		c.SetCookie(tp.CookieName, ts, int(tp.CookieMaxAge.Seconds()), tp.CookiePath, tp.CookieDomain, tp.CookieSecure, tp.CookieHttpOnly)
	} else {
		c.Logger.Errorf("EncryptToken: %v", err)
	}
	return ts
}
