package gmw

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/pandafw/pango/gin"
)

func doLocalizerTest(t *testing.T, req *http.Request, want string) {
	w := httptest.NewRecorder()
	router := gin.New()
	router.Use(NewLocalizer("ja", "zh").Handler())
	router.Any("/", func(c *gin.Context) {
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
	req.Header.Add(LocaleHeaderName, "zh;ja")

	doLocalizerTest(t, req, "zh")
}

func TestLocalizerQueryString(t *testing.T) {
	req, _ := http.NewRequest("GET", "/?__locale=zh", nil)

	doLocalizerTest(t, req, "zh")
}

func TestLocalizerPostForm(t *testing.T) {
	req, _ := http.NewRequest("POST", "/", nil)
	req.PostForm = url.Values{}
	req.PostForm.Add("__locale", "zh")

	doLocalizerTest(t, req, "zh")
}

func TestLocalizerCookie(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{
		Name:  LocaleCookieName,
		Value: "zh",
	})

	doLocalizerTest(t, req, "zh")
}
