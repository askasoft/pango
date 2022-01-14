package render

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRenderData(t *testing.T) {
	w := httptest.NewRecorder()
	data := []byte("#!PNG some raw data")

	err := (Data{
		ContentType: "image/png",
		Data:        data,
	}).Render(w)

	assert.NoError(t, err)
	assert.Equal(t, "#!PNG some raw data", w.Body.String())
	assert.Equal(t, "image/png", w.Header().Get("Content-Type"))
}
