package ini

import (
	"bytes"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/askasoft/pango/ars"
	"github.com/askasoft/pango/bye"
	"github.com/askasoft/pango/iox"
)

func TestLoadEmpty(t *testing.T) {
	ini := NewIni()

	buf := bytes.NewBuffer([]byte{})
	if err := ini.LoadData(buf); err != nil {
		t.Errorf("ini.LoadData('') = %v", err)
	}
}

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

func testReadString(filename string) string {
	bs, _ := os.ReadFile(filename)
	return bye.UnsafeString(bs)
}

func testLoadFile(t *testing.T, fsrc, fexp, fout string) {
	defer os.Remove(fout)

	ini := NewIni()
	ini.EOL = iox.CRLF

	// load
	if ini.LoadFile(fsrc, true) != nil {
		t.Errorf("ini.LoadFile(%q) != nil", fsrc)
	}

	// empty
	global := ini.GetSection("")
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
	other := ini.GetSection("other")
	if other == nil {
		t.Error(`ini.GetSection("other") == nil`)
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
	sexp := testReadString(fexp)

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

		sout := testReadString(fout)
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

	// Map
	{
		sec := ini.Section("multi-2")
		sm := sec.StringMap()
		if len(sm) != 1 || sm["test"] != "a" {
			t.Errorf(`sec.StringMap() = %v`, sm)
		}

		ssm := sec.StringsMap()
		if len(ssm) != 1 || !ars.EqualStrings(ssm["test"], []string{"a", "b", "<tab>\t<tab>"}) {
			t.Errorf(`sec.StringsMap() = %v`, ssm)
		}

		om := sec.Map()
		if len(om) != 1 || !ars.EqualStrings(ssm["test"], []string{"a", "b", "<tab>\t<tab>"}) {
			t.Errorf(`sec.Map() = %v`, ssm)
		}
	}
}

func TestIniCopy(t *testing.T) {
	f1 := "testdata/cm1.ini"
	f2 := "testdata/cm2.ini"
	fexp := "testdata/copy-expect.ini"
	fout := "testdata/copy-output.ini"

	i1 := NewIni()
	i1.EOL = iox.CRLF
	if err := i1.LoadFile(f1); err != nil {
		t.Fatalf(`Failed to load %s: %v`, f1, err)
	}

	i2 := NewIni()
	i2.EOL = iox.CRLF
	if err := i2.LoadFile(f2); err != nil {
		t.Fatalf(`Failed to load %s: %v`, f2, err)
	}

	i1.Copy(i2)

	// expected file
	sexp := testReadString(fexp)

	// write data
	{
		sout := &strings.Builder{}
		if err := i1.WriteData(sout); err != nil {
			t.Fatalf(`ini.WriteData(sout) = %v`, err)
		}
		if sexp != sout.String() {
			i1.WriteFile(fout)
			t.Fatalf(`ini.WriteData(sout)\n actual: %v\n   want: %v`, sout.String(), sexp)
		}
	}
}

func TestIniMerge(t *testing.T) {
	f1 := "testdata/cm1.ini"
	f2 := "testdata/cm2.ini"
	fexp := "testdata/merge-expect.ini"
	fout := "testdata/merge-output.ini"

	i1 := NewIni()
	i1.EOL = iox.CRLF
	if err := i1.LoadFile(f1, true); err != nil {
		t.Fatalf(`Failed to load %s: %v`, f1, err)
	}

	i2 := NewIni()
	i2.EOL = iox.CRLF
	if err := i2.LoadFile(f2, true); err != nil {
		t.Fatalf(`Failed to load %s: %v`, f2, err)
	}

	i1.Merge(i2)

	// expected file
	sexp := testReadString(fexp)

	// write data
	{
		sout := &strings.Builder{}
		if err := i1.WriteData(sout); err != nil {
			t.Fatalf(`ini.WriteData(sout) = %v`, err)
		}
		if sexp != sout.String() {
			i1.WriteFile(fout)
			t.Fatalf(`ini.WriteData(sout)\n actual: %v\n   want: %v`, sout.String(), sexp)
		}
	}
}
