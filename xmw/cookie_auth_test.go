package xmw

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/askasoft/pango/xin"
	"github.com/stretchr/testify/assert"
)

func TestCookieAuthSucceed(t *testing.T) {
	accounts := testAccounts{"admin": {"admin", "password"}}

	router := xin.New()

	ca := NewCookieAuth(accounts.FindUser, "1234567890abced")
	router.Use(ca.Handler())

	router.GET("/login", func(c *xin.Context) {
		c.String(http.StatusOK, c.MustGet(AuthUserKey).(*testAccount).username)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/login", nil)

	val, _ := ca.encrypt("admin", "password")
	req.AddCookie(&http.Cookie{
		Name:  AuthCookieName,
		Value: val,
	})
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "admin", w.Body.String())
}

func TestCookieAuthRedirect(t *testing.T) {
	called := false
	accounts := testAccounts{"foo": {"foo", "bar"}}

	router := xin.New()

	ca := NewCookieAuth(accounts.FindUser, "1234567890abcdefg")
	ca.RedirectURL = "/redirect"
	router.Use(ca.Handler())

	router.GET("/login", func(c *xin.Context) {
		called = true
		c.String(http.StatusOK, c.MustGet(AuthUserKey).(*testAccount).username)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/login", nil)

	val, _ := ca.encrypt("admin", "password")
	req.AddCookie(&http.Cookie{
		Name:  AuthCookieName,
		Value: val,
	})
	router.ServeHTTP(w, req)

	assert.False(t, called)
	assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
	assert.Equal(t, "/redirect?origin=%2Flogin", w.Header().Get("Location"))
}
