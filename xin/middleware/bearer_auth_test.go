package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/askasoft/pango/test/assert"
	"github.com/askasoft/pango/xin"
)

func testBearerAuthHeader(token string) string {
	return "Bearer " + token
}

func TestBearerAuthSucceed(t *testing.T) {
	accounts := testAccounts{"admin": {"admin", "password"}}
	router := xin.New()
	router.Use(NewBearerAuth(accounts.FindUser).Handle)
	router.GET("/login", func(c *xin.Context) {
		c.String(http.StatusOK, c.MustGet(AuthUserKey).(*testAccount).username)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/login", nil)
	req.Header.Set("Authorization", testBearerAuthHeader("admin"))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "admin", w.Body.String())
}

func TestBearerAuth401(t *testing.T) {
	called := false
	accounts := testAccounts{"foo": {"foo", "bar"}}
	router := xin.New()
	router.Use(NewBearerAuth(accounts.FindUser).Handle)
	router.GET("/login", func(c *xin.Context) {
		called = true
		c.String(http.StatusOK, c.MustGet(AuthUserKey).(*testAccount).username)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/login", nil)
	req.Header.Set("Authorization", testBearerAuthHeader("admin"))
	router.ServeHTTP(w, req)

	assert.False(t, called)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Equal(t, "Bearer realm=\"Authorization Required\"", w.Header().Get("WWW-Authenticate"))
}

func TestBearerAuth401WithCustomRealm(t *testing.T) {
	called := false
	accounts := testAccounts{"foo": {"foo", "bar"}}
	router := xin.New()
	ba := NewBearerAuth(accounts.FindUser)
	ba.Realm = `My Custom "Realm"`
	router.Use(ba.Handle)
	router.GET("/login", func(c *xin.Context) {
		called = true
		c.String(http.StatusOK, c.MustGet(AuthUserKey).(*testAccount).username)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/login", nil)
	req.Header.Set("Authorization", testBearerAuthHeader("admin"))
	router.ServeHTTP(w, req)

	assert.False(t, called)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Equal(t, "Bearer realm=\"My Custom \\\"Realm\\\"\"", w.Header().Get("WWW-Authenticate"))
}
