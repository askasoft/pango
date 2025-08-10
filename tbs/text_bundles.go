package tbs

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/askasoft/pango/asg"
	"github.com/askasoft/pango/bol"
	"github.com/askasoft/pango/ini"
	"github.com/askasoft/pango/num"
	"github.com/askasoft/pango/str"
)

// TextBundles a localized text resource bundle container
type TextBundles struct {
	Extensions []string  // file extensions
	Timestamp  time.Time // modified timestamp

	bundles map[string]*ini.Ini
}

// NewTextBundles returns a new TextBundles instance
func NewTextBundles(extensions ...string) *TextBundles {
	if len(extensions) == 0 {
		extensions = []string{".ini", ".txt"}
	}
	return &TextBundles{
		Extensions: extensions,
		Timestamp:  time.Now(),
		bundles:    make(map[string]*ini.Ini),
	}
}

// Clear clear all loaded text resources
func (tbs *TextBundles) Clear() {
	tbs.bundles = make(map[string]*ini.Ini)
	tbs.Timestamp = time.Now()
}

// Load glob and parse text files under root path
func (tbs *TextBundles) Load(root string) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		return tbs.loadFile(nil, path)
	})
}

// LoadFS glob and parse template files from FS
func (tbs *TextBundles) LoadFS(fsys fs.FS, root string) error {
	return fs.WalkDir(fsys, root, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		return tbs.loadFile(fsys, path)
	})
}

// loadFile load from path or fsys
func (tbs *TextBundles) loadFile(fsys fs.FS, path string) error {
	ext := filepath.Ext(path)
	if !asg.Contains(tbs.Extensions, ext) {
		return nil
	}

	name := str.TrimSuffix(filepath.Base(path), ext)
	locale := str.SubstrAfterLastByte(name, '_')

	bundle, ok := tbs.bundles[locale]
	if !ok {
		bundle = ini.NewIni()
		tbs.bundles[locale] = bundle
	}

	if fsys != nil {
		if err := bundle.LoadFileFS(fsys, path); err != nil {
			return fmt.Errorf("TextBundles load ini %q error: %w", path, err)
		}
	} else {
		if err := bundle.LoadFile(path); err != nil {
			return fmt.Errorf("TextBundles load ini %q error: %w", path, err)
		}
	}

	tbs.Timestamp = time.Now()
	return nil
}

// GetBundle get target locale ini bundle
func (tbs *TextBundles) GetBundle(locale string) *ini.Ini {
	bundles := []*ini.Ini{}
	for locale != "" {
		if bundle, ok := tbs.bundles[locale]; ok {
			bundles = append(bundles, bundle)
		}

		if l2, _, ok := str.LastCutByte(locale, '-'); ok {
			locale = l2
		} else {
			break
		}
	}

	if bundle, ok := tbs.bundles[""]; ok {
		bundles = append(bundles, bundle)
	}

	b := ini.NewIni()
	for i := len(bundles) - 1; i >= 0; i-- {
		b.Copy(bundles[i])
	}
	return b
}

// get get the target locale string
func (tbs *TextBundles) get(locale, sect, name string) (string, bool) {
	for locale != "" {
		if bundle, ok := tbs.bundles[locale]; ok {
			if sec := bundle.GetSection(sect); sec != nil {
				if val := sec.Get(name); val != "" {
					return val, ok
				}
			}
		}

		if l2, _, ok := str.LastCutByte(locale, '-'); ok {
			locale = l2
		} else {
			break
		}
	}

	if bundle, ok := tbs.bundles[""]; ok {
		if sec := bundle.GetSection(sect); sec != nil {
			if val := sec.Get(name); val != "" {
				return val, ok
			}
		}
	}

	return "", false
}

// GetBool get a bool value of the key from the text bundle.
// if not found or convert error, returns the first non-false value from defs.
func (tbs *TextBundles) GetBool(locale, key string, defs ...bool) bool {
	return bol.Atob(tbs.GetText(locale, key), defs...)
}

// GetInt get a int value of the key from the text bundle.
// if not found or convert error, returns the first non-zero value from defs.
func (tbs *TextBundles) GetInt(locale, key string, defs ...int) int {
	return num.Atoi(tbs.GetText(locale, key), defs...)
}

// GetInt64 get a int64 value of the key from the text bundle.
// if not found or convert error, returns the first non-zero value from defs.
func (tbs *TextBundles) GetInt64(locale, key string, defs ...int64) int64 {
	return num.Atol(tbs.GetText(locale, key), defs...)
}

// GetFloat get a float64 value of the key from the text bundle.
// if not found or convert error, returns the first non-zero value from defs.
func (tbs *TextBundles) GetFloat(locale, key string, defs ...float64) float64 {
	return num.Atof(tbs.GetText(locale, key), defs...)
}

// GetText get the locale text by key.
// if not found, returns the first non-empty value from defs.
func (tbs *TextBundles) GetText(locale, key string, defs ...string) string {
	sect, name := "", key

	if i := str.LastIndexByte(key, '.'); i >= 0 {
		sect, name = key[:i], key[i+1:]
	}

	if val, ok := tbs.get(locale, sect, name); ok {
		return val
	}

	return str.NonEmpty(defs...)
}

// Format use fmt.Sprintf to format the locale text by key and args.
func (tbs *TextBundles) Format(locale, key string, args ...any) string {
	format := tbs.GetText(locale, key)

	if format == "" || len(args) == 0 {
		return format
	}

	return fmt.Sprintf(format, args...)
}

// Replace use strings.Replacer to translate content to the locale language.
func (tbs *TextBundles) Replace(locale, key string, args ...any) string {
	var defs []string

	if len(args)&1 != 0 {
		if s, ok := args[0].(string); ok {
			defs = append(defs, s)
			args = args[1:]
		} else {
			panic("tbs: invalid arguments")
		}
	}

	format := tbs.GetText(locale, key, defs...)
	if format == "" || len(args) == 0 {
		return format
	}

	sargs := make([]string, len(args))
	for i, a := range args {
		sargs[i] = fmt.Sprint(a)
	}
	repl := strings.NewReplacer(sargs...)

	return repl.Replace(format)
}

// Error create a error with the locale text by key.
func (tbs *TextBundles) Error(locale, key string, defs ...string) error {
	return errors.New(tbs.GetText(locale, key, defs...))
}

// Errorf create a error with the locale text by key and args.
func (tbs *TextBundles) Errorf(locale, key string, args ...any) error {
	return errors.New(tbs.Format(locale, key, args...))
}
