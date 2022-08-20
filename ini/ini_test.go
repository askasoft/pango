package ini

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/pandafw/pango/iox"
)

func TestLoadFile(t *testing.T) {
	fsrc := "testdata/source.ini"
	fexp := "testdata/expect.ini"
	fout := "testdata/output.ini"
	testLoadFile(t, fsrc, fexp, fout)
}

func TestLoadFileBom(t *testing.T) {
	fsrc := "testdata/source-bom.ini"
	fexp := "testdata/expect.ini"
	fout := "testdata/output-bom.ini"
	testLoadFile(t, fsrc, fexp, fout)
}

func testLoadFile(t *testing.T, fsrc, fexp, fout string) {
	defer os.Remove(fout)

	ini := NewIni()
	ini.EOL = iox.CRLF
	ini.Multiple = true

	// load
	if ini.LoadFile(fsrc) != nil {
		t.Errorf("ini.LoadFile(%q) != nil", fsrc)
	}

	// empty
	global := ini.Section("")
	if v := global.GetString("empty", "def"); v != "def" {
		t.Errorf(`global.GetString("empty", "def") = %v, want "def"`, v)
	}
	if v := global.GetInt("empty", 1); v != 1 {
		t.Errorf(`global.GetInt("empty", 1) = %v, want 1`, v)
	}
	if v := global.GetInt64("empty", 10); v != 10 {
		t.Errorf(`global.GetInt64("empty", 10) = %v, want 10`, v)
	}
	if v := global.GetBool("empty", true); v != true {
		t.Errorf(`global.GetBool("empty", true) = %v, want true`, v)
	}
	if v := global.GetFloat("empty", 1.1); v != 1.1 {
		t.Errorf(`global.GetFloat("empty", 1.1) = %v, want 1.1`, v)
	}
	if v := global.GetDuration("empty", time.Minute); v != time.Minute {
		t.Errorf(`global.GetDuration("empty", time.Minute) = %v, want 1m`, v)
	}

	// value
	other := ini.Section("other")
	if other == nil {
		t.Error(`ini.Section("other") == nil`)
	} else {
		dec := other.GetInt("dec")
		if 42 != dec {
			t.Errorf(`other.GetInt("dec") = %v, want %v`, dec, 42)
		}
		hex := other.GetInt("hex")
		if 42 != hex {
			t.Errorf(`other.GetInt("hex") = %v, want %v`, dec, 42)
		}
		vtrue := other.GetBool("true")
		if !vtrue {
			t.Error(`other.GetBool("true") != true`)
		}
		vfalse := other.GetBool("false")
		if vfalse {
			t.Error(`other.GetBool("false") != false`)
		}
	}

	// expected file
	bexp, _ := ioutil.ReadFile(fexp)
	sexp := string(bexp)

	// write data
	{
		sout := &strings.Builder{}
		werr := ini.WriteData(sout)
		if werr != nil {
			t.Errorf(`ini.WriteData(sout) = %v`, werr)
		}
		if sexp != sout.String() {
			t.Errorf(`ini.WriteData(sout)\n actual: %v\n   want: %v`, sout.String(), sexp)
		}
	}

	// write file
	{
		werr := ini.WriteFile(fout)
		if werr != nil {
			t.Errorf(`ini.WriteFile(fout) = %v`, werr)
		}

		bout, _ := ioutil.ReadFile(fout)
		sout := string(bout)
		if sexp != sout {
			t.Errorf(`ini.WriteFile(fout)\n actual: %v\n   want: %v`, fout, sexp)
		}
	}

	// remove section
	{
		sec := ini.RemoveSection("")
		if sec == nil {
			t.Errorf(`ini.RemoveSection("") = %v`, sec)
		}
		sec = ini.RemoveSection("not exist")
		if sec != nil {
			t.Errorf(`ini.RemoveSection("not exist") = %v`, sec)
		}
	}
}
