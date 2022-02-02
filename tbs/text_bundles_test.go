package tbs

import (
	"embed"
	"testing"
)

func TestLoad(t *testing.T) {
	ts := NewTS()
	root := "testdata"

	err := ts.Load(root)
	if err != nil {
		t.Errorf(`ts.Load(%q) = %v`, root, err)
		return
	}

	testLoad(t, ts)
}

//go:embed testdata
var testdata embed.FS

func TestLoadFS(t *testing.T) {
	ts := NewTS()
	root := "testdata"

	err := ts.LoadFS(testdata, root)
	if err != nil {
		t.Errorf(`ts.LoadFS(%q) = %v`, root, err)
		return
	}

	testLoad(t, ts)
}

func testLoad(t *testing.T, ts *TS) {
	cs := []struct {
		lang string
		name string
		args []interface{}
		want string
	}{
		{"en", "title", nil, "hello world"},
		{"en", "label.welcome", []interface{}{"home"}, "welcome home"},
		{"ja-JP", "title", nil, "こんにちは世界"},
		{"ja-JP", "label.welcome", []interface{}{"ダーリン"}, "ようこそ ダーリン"},
	}

	for i, c := range cs {
		a := ts.Format(c.lang, c.name, c.args...)
		if a != c.want {
			t.Errorf("%d Foramt(%q, %q, %v) = %q, want %q", i, c.lang, c.name, c.args, a, c.want)
		}
	}
}
