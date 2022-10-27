package gmw

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/pandafw/pango/gin"
)

func TestRequestLimiter(t *testing.T) {
	req, _ := http.NewRequest("POST", "/", strings.NewReader(`{"a": "1345678901"}`))

	w := httptest.NewRecorder()
	router := gin.New()
	router.Use(NewRequestLimiter(10).Handler())
	router.POST("/", func(c *gin.Context) {
		m := map[string]string{}
		if err := c.BindJSON(&m); err == nil {
			c.String(200, "OK")
		}
	})

	router.ServeHTTP(w, req)

	val := w.Result().StatusCode
	exp := http.StatusBadRequest
	if val != exp {
		t.Errorf("%v = %q, want %q", req.URL.String(), val, exp)
	}
}
