package middleware

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/askasoft/pango/xin"
)

func doLocalizerTest(t *testing.T, req *http.Request, want string) {
	w := httptest.NewRecorder()
	router := xin.New()
	router.Use(NewLocalizer("ja", "zh").Handle)
	router.Any("/", func(c *xin.Context) {
		c.String(200, c.Locale)
	})

	router.ServeHTTP(w, req)

	if w.Body.String() != want {
		t.Errorf("%v = %q, want %q", req.URL.String(), w.Body.String(), want)
	}
}

func TestLocalizerAcceptLanguages1(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Add("Accept-Languages", "en;en-US")

	doLocalizerTest(t, req, "ja")
}

func TestLocalizerAcceptLanguages2(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Add("Accept-Languages", "ja;zh")

	doLocalizerTest(t, req, "ja")
}

func TestLocalizerHttpHeader(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Add(LocaleHeaderNames[0], "zh;ja")

	doLocalizerTest(t, req, "zh")
}

func TestLocalizerQueryString(t *testing.T) {
	req, _ := http.NewRequest("GET", "/?_locale_=zh", nil)

	doLocalizerTest(t, req, "zh")
}

func TestLocalizerPostForm(t *testing.T) {
	req, _ := http.NewRequest("POST", "/", nil)
	req.PostForm = url.Values{}
	req.PostForm.Add("_locale_", "zh")

	doLocalizerTest(t, req, "zh")
}

func TestLocalizerCookie(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{
		Name:  LocaleCookieNames[0],
		Value: "zh",
	})

	doLocalizerTest(t, req, "zh")
}
