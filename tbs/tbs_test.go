package tbs

import (
	"testing"
)

func TestPkgLoad(t *testing.T) {
	Clear()
	err := Load(testroot)
	if err != nil {
		t.Errorf(`tbs.Load(%q) = %v`, testroot, err)
		return
	}

	testFormat(t, func(locale, format string, args ...any) string {
		return Format(locale, format, args...)
	})
}

func TestPkgLoadFS(t *testing.T) {
	Clear()
	err := LoadFS(testdata, testroot)
	if err != nil {
		t.Errorf(`tbs.LoadFS(%q) = %v`, testroot, err)
		return
	}

	testFormat(t, func(locale, format string, args ...any) string {
		return Format(locale, format, args...)
	})
}
