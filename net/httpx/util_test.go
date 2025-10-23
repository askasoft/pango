package httpx

import (
	"net/http"
	"testing"

	"github.com/askasoft/pango/test/assert"
)

func TestBodyAllowedForStatus(t *testing.T) {
	assert.False(t, false, BodyAllowedForStatus(http.StatusProcessing))
	assert.False(t, false, BodyAllowedForStatus(http.StatusNoContent))
	assert.False(t, false, BodyAllowedForStatus(http.StatusNotModified))
	assert.True(t, true, BodyAllowedForStatus(http.StatusInternalServerError))
}
