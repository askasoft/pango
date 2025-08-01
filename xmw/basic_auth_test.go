package xmw

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/askasoft/pango/str"
	"github.com/askasoft/pango/test/assert"
	"github.com/askasoft/pango/xin"
)

func testBasicAuthHeader(user, password string) string {
	base := user + ":" + password
	return "Basic " + base64.StdEncoding.EncodeToString(str.UnsafeBytes(base))
}

func TestBasicAuthSucceed(t *testing.T) {
	accounts := testAccounts{"admin": {"admin", "password"}}
	router := xin.New()
	router.Use(NewBasicAuth(accounts.FindUser).Handle)
	router.GET("/login", func(c *xin.Context) {
		c.String(http.StatusOK, c.MustGet(AuthUserKey).(*testAccount).username)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/login", nil)
	req.Header.Set("Authorization", testBasicAuthHeader("admin", "password"))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "admin", w.Body.String())
}

func TestBasicAuth401(t *testing.T) {
	called := false
	accounts := testAccounts{"foo": {"foo", "bar"}}
	router := xin.New()
	router.Use(NewBasicAuth(accounts.FindUser).Handle)
	router.GET("/login", func(c *xin.Context) {
		called = true
		c.String(http.StatusOK, c.MustGet(AuthUserKey).(*testAccount).username)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/login", nil)
	req.Header.Set("Authorization", testBasicAuthHeader("admin", "password"))
	router.ServeHTTP(w, req)

	assert.False(t, called)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Equal(t, "Basic realm=\"Authorization Required\"", w.Header().Get("WWW-Authenticate"))
}

func TestBasicAuth401WithCustomRealm(t *testing.T) {
	called := false
	accounts := testAccounts{"foo": {"foo", "bar"}}
	router := xin.New()
	ba := NewBasicAuth(accounts.FindUser)
	ba.Realm = `My Custom "Realm"`
	router.Use(ba.Handle)
	router.GET("/login", func(c *xin.Context) {
		called = true
		c.String(http.StatusOK, c.MustGet(AuthUserKey).(*testAccount).username)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/login", nil)
	req.Header.Set("Authorization", testBasicAuthHeader("admin", "password"))
	router.ServeHTTP(w, req)

	assert.False(t, called)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Equal(t, "Basic realm=\"My Custom \\\"Realm\\\"\"", w.Header().Get("WWW-Authenticate"))
}
