package tbs

import (
	"io/fs"

	"github.com/askasoft/pango/ini"
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

// GetBundle get target locale ini bundle
func GetBundle(locale string) *ini.Ini {
	return _tbs.GetBundle(locale)
}

// GetText get the locale text by key.
func GetText(locale, key string, defs ...string) string {
	return _tbs.GetText(locale, key, defs...)
}

// Format translate content to target language.
func Format(locale, format string, args ...any) string {
	return _tbs.Format(locale, format, args...)
}

// Replace use strings.Replacer to translate content to the locale language.
func Replace(locale, format string, args ...any) string {
	return _tbs.Replace(locale, format, args...)
}
