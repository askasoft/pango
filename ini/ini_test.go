package ini

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/pandafw/pango/iox"
	"github.com/stretchr/testify/assert"
)

func TestLoadFile(t *testing.T) {
	fin := "testdata/input.ini"
	fexp := "testdata/except.ini"
	fout := "testdata/output.ini"

	defer os.Remove(fout)

	ini := NewIni()
	ini.EOL = iox.CRLF
	ini.Multiple = true

	// load
	assert.Nil(t, ini.LoadFile(fin))

	// value
	other := ini.Section("other")
	assert.NotNil(t, other)
	assert.Equal(t, 42, other.GetInt("dec"))
	assert.Equal(t, 42, other.GetInt("hex"))
	assert.True(t, other.GetBool("true"))
	assert.False(t, other.GetBool("false"))

	// expected file
	bexp, _ := ioutil.ReadFile(fexp)
	sexp := string(bexp)

	// write data
	{
		sout := &strings.Builder{}
		assert.Nil(t, ini.WriteData(sout))
		assert.Equal(t, sexp, sout.String())
	}

	// write file
	{
		assert.Nil(t, ini.WriteFile(fout))
		bout, _ := ioutil.ReadFile(fout)
		sout := string(bout)
		assert.Equal(t, sexp, sout)
	}

	// remove section
	{
		assert.NotNil(t, ini.RemoveSection(""))
		assert.Nil(t, ini.RemoveSection("not exist"))
	}
}
