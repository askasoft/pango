package binding

import (
	"testing"

	"github.com/askasoft/pango/test/assert"
	"github.com/askasoft/pango/test/require"
)

func TestXMLBindingBindBody(t *testing.T) {
	var s struct {
		Foo string `xml:"foo"`
	}
	xmlBody := `<?xml version="1.0" encoding="UTF-8"?>
<root>
   <foo>FOO</foo>
</root>`
	err := xmlBinding{}.BindBody([]byte(xmlBody), &s)
	require.NoError(t, err)
	assert.Equal(t, "FOO", s.Foo)
}
