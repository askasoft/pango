package tbs

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"

	"github.com/pandafw/pango/ars"
	"github.com/pandafw/pango/ini"
	"github.com/pandafw/pango/str"
)

// TS a localized text resource bundle container
type TS struct {
	Extensions []string // file extensions

	bundles map[string]*ini.Ini
}

// NewTS new txt resource bundle container
func NewTS(extensions ...string) *TS {
	if len(extensions) == 0 {
		extensions = []string{".ini", ".txt"}
	}
	return &TS{
		Extensions: extensions,
		bundles:    map[string]*ini.Ini{},
	}
}

// Load glob and parse text files under root path
func (ts *TS) Load(root string) error {
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		return ts.loadFile(nil, path)
	})

	if err != nil {
		return err
	}

	return nil
}

// loadFile load template file
func (ts *TS) loadFile(fsys fs.FS, path string) error {
	ext := filepath.Ext(path)
	if !ars.ContainsString(ts.Extensions, ext) {
		return nil
	}

	name := str.TrimSuffix(filepath.Base(path), ext)
	locale := str.SubstrAfterLastByte(name, '_')

	bundle, ok := ts.bundles[locale]
	if !ok {
		bundle = ini.NewIni()
		ts.bundles[locale] = bundle
	}

	if fsys != nil {
		if err := bundle.LoadFileFS(fsys, path); err != nil {
			return fmt.Errorf("TS load ini %q error: %w", path, err)
		}
	} else {
		if err := bundle.LoadFile(path); err != nil {
			return fmt.Errorf("TS load ini %q error: %w", path, err)
		}
	}

	return nil
}

// LoadFS glob and parse template files from FS
func (ts *TS) LoadFS(fsys fs.FS, root string) error {
	err := fs.WalkDir(fsys, root, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		return ts.loadFile(fsys, path)
	})

	if err != nil {
		return err
	}

	return nil
}

// Get target locale string
func (ts *TS) Get(locale, section, name string) (string, bool) {
	for locale != "" {
		if bundle, ok := ts.bundles[locale]; ok {
			if sec := bundle.Section(section); sec != nil {
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

	if bundle, ok := ts.bundles[locale]; ok {
		if sec := bundle.Section(section); sec != nil {
			if val := sec.Get(name); val != "" {
				return val, ok
			}
		}
	}

	return "", false
}

// Format translate content to target language.
func (ts *TS) Format(locale, format string, args ...any) string {
	section := ""

	dot := str.LastIndexByte(format, '.')
	if dot >= 0 {
		section = format[:dot]
		format = format[dot+1:]
	}

	if value, ok := ts.Get(locale, section, format); ok {
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
