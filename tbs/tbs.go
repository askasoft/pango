package tbs

import (
	"io/fs"

	"github.com/pandafw/pango/ini"
)

var _tbs = NewTextBundles()

// Default returns the default TextBundles instance used by the package-level functions.
func Default() *TextBundles {
	return _tbs
}

// SetDefault set the default TextBundles instance used by the package-level functions.
func SetDefault(tbs *TextBundles) {
	_tbs = tbs
}

// Clear clear all loaded text resources
func Clear() {
	_tbs.Clear()
}

// Load glob and parse text files under root path
func Load(root string) error {
	return _tbs.Load(root)
}

// LoadFS glob and parse template files from FS
func LoadFS(fsys fs.FS, root string) error {
	return _tbs.LoadFS(fsys, root)
}

// GetBundle get all target locale strings
func GetBundle(locale string) *ini.Ini {
	return _tbs.GetBundle(locale)
}

// GetAll get all target locale strings
func GetAll(locale string) map[string]map[string]string {
	return _tbs.GetAll(locale)
}

// Get target locale string
func Get(locale, section, name string) (string, bool) {
	return _tbs.Get(locale, section, name)
}

// Format translate content to target language.
func Format(locale, format string, args ...any) string {
	return _tbs.Format(locale, format, args...)
}
