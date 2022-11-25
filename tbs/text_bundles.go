package tbs

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"time"

	"github.com/pandafw/pango/ars"
	"github.com/pandafw/pango/ini"
	"github.com/pandafw/pango/str"
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
	if !ars.ContainsString(tbs.Extensions, ext) {
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

// GetBundle get all target locale strings
func (tbs *TextBundles) GetBundle(locale string) *ini.Ini {
	if bundle, ok := tbs.bundles[locale]; ok {
		return bundle
	}

	if locale != "" {
		if bundle, ok := tbs.bundles[""]; ok {
			return bundle
		}
	}

	return nil
}

// GetAll get all target locale strings
func (tbs *TextBundles) GetAll(locale string) map[string]map[string]string {
	bundle := tbs.GetBundle(locale)
	if bundle != nil {
		return bundle.StringMap()
	}
	return nil
}

// Get target locale string
func (tbs *TextBundles) Get(locale, section, name string) (string, bool) {
	for locale != "" {
		if bundle, ok := tbs.bundles[locale]; ok {
			if sec := bundle.GetSection(section); sec != nil {
				if val := sec.Get(name); val != "" {
					return val, ok
				}
			}
		}

		locale2 := str.SubstrBeforeByte(locale, '-')
		if locale2 == locale {
			locale = ""
		} else {
			locale = locale2
		}
	}

	if bundle, ok := tbs.bundles[locale]; ok {
		if sec := bundle.GetSection(section); sec != nil {
			if val := sec.Get(name); val != "" {
				return val, ok
			}
		}
	}

	return "", false
}

// Format translate content to target language.
func (tbs *TextBundles) Format(locale, format string, args ...any) string {
	section := ""

	dot := str.LastIndexByte(format, '.')
	if dot >= 0 {
		section = format[:dot]
		format = format[dot+1:]
	}

	if value, ok := tbs.Get(locale, section, format); ok {
		format = value
	}

	if len(args) == 0 {
		return format
	}

	params := make([]any, 0, len(args))
	for _, arg := range args {
		if arg != nil {
			val := reflect.ValueOf(arg)
			if val.Kind() == reflect.Slice {
				for i := 0; i < val.Len(); i++ {
					params = append(params, val.Index(i).Interface())
				}
			} else {
				params = append(params, arg)
			}
		}
	}
	return fmt.Sprintf(format, params...)
}
