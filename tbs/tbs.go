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

// GetBool get a bool value of the key from the text bundle.
// if not found or convert error, returns the default defs[0] value.
func GetBool(locale, key string, defs ...bool) bool {
	return _tbs.GetBool(locale, key, defs...)
}

// GetInt get a int value of the key from the text bundle.
// if not found or convert error, returns the default defs[0] value.
func GetInt(locale, key string, defs ...int) int {
	return _tbs.GetInt(locale, key, defs...)
}

// GetInt64 get a int64 value of the key from the text bundle.
// if not found or convert error, returns the default defs[0] value.
func GetInt64(locale, key string, defs ...int64) int64 {
	return _tbs.GetInt64(locale, key, defs...)
}

// GetFloat get a float64 value of the key from the text bundle.
// if not found or convert error, returns the default defs[0] value.
func GetFloat(locale, key string, defs ...float64) float64 {
	return _tbs.GetFloat(locale, key, defs...)
}

// GetText get the locale text by key.
func GetText(locale, key string, defs ...string) string {
	return _tbs.GetText(locale, key, defs...)
}

// Format use fmt.Sprintf to format the locale text by key and args.
func Format(locale, format string, args ...any) string {
	return _tbs.Format(locale, format, args...)
}

// Replace use strings.Replacer to translate content to the locale language.
func Replace(locale, format string, args ...any) string {
	return _tbs.Replace(locale, format, args...)
}

// Error create a error with the locale text by key.
func Error(locale, key string, defs ...string) error {
	return _tbs.Error(locale, key, defs...)
}

// Errorf create a error with the locale text by key and args.
func Errorf(locale, key string, args ...any) error {
	return _tbs.Errorf(locale, key, args...)
}
