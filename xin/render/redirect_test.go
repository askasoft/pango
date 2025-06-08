package render

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/askasoft/pango/test/assert"
)

func TestRenderRedirect(t *testing.T) {
	req, err := http.NewRequest("GET", "/test-redirect", nil)
	assert.NoError(t, err)

	data1 := Redirect{
		Code:     http.StatusMovedPermanently,
		Request:  req,
		Location: "/new/location",
	}

	w := httptest.NewRecorder()
	err = data1.Render(w)
	assert.NoError(t, err)

	data2 := Redirect{
		Code:     http.StatusOK,
		Request:  req,
		Location: "/new/location",
	}

	w = httptest.NewRecorder()
	assert.PanicsWithValue(t, "Cannot redirect with status code 200", func() {
		err := data2.Render(w)
		assert.NoError(t, err)
	})

	data3 := Redirect{
		Code:     http.StatusCreated,
		Request:  req,
		Location: "/new/location",
	}

	w = httptest.NewRecorder()
	err = data3.Render(w)
	assert.NoError(t, err)

	// only improve coverage
	data2.WriteContentType(w)
}
