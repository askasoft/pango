package gmw

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pandafw/pango/gin"
)

func TestHTTPHeader(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)

	w := httptest.NewRecorder()
	router := gin.New()
	router.Use(NewHTTPHeader(map[string]string{"Access-Control-Allow-Origin": "*"}).Handler())
	router.GET("/", func(c *gin.Context) {
		c.String(200, "OK")
	})

	router.ServeHTTP(w, req)

	val := w.Header().Get("Access-Control-Allow-Origin")
	exp := "*"
	if val != exp {
		t.Errorf("%v = %q, want %q", req.URL.String(), val, exp)
	}
}
