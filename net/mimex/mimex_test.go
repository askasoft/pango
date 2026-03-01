package mimex

import (
	"testing"

	"github.com/askasoft/pango/test/assert"
)

func TestMediaTypeByFilename(t *testing.T) {
	cs := []struct {
		n string
		w string
	}{
		{"a.txt", "text/plain"},
		{"a.html", "text/html"},
	}

	for _, c := range cs {
		assert.Equal(t, c.w, MediaTypeByFilename(c.n))
	}
}

func TestMediaTypeByExtension(t *testing.T) {
	cs := []struct {
		n string
		w string
	}{
		{".txt", "text/plain"},
		{".html", "text/html"},
	}

	for _, c := range cs {
		assert.Equal(t, c.w, MediaTypeByExtension(c.n))
	}
}
