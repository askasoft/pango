package ini

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/pandafw/pango/bye"
	"github.com/pandafw/pango/col"
	"github.com/pandafw/pango/iox"
	"github.com/pandafw/pango/str"
)

// Ini INI file reader / writer
type Ini struct {
	sections *col.LinkedHashMap // Parsed sections
	EOL      string             // End of Line
	Multiple bool               // Multiple entry with same key
}

// NewIni create a Ini
func NewIni() *Ini {
	ini := &Ini{
		sections: col.NewLinkedHashMap(),
		EOL:      iox.EOL,
	}

	ini.NewSection("") // init global section
	return ini
}

// Clear clears the ini
func (ini *Ini) Clear() {
	ini.sections.Clear()
	ini.NewSection("") // init global section
}

// IsEmpty returns true if the Ini has no entry
func (ini *Ini) IsEmpty() bool {
	if ini.sections.IsEmpty() {
		return true
	}

	for it := ini.sections.Iterator(); it.Next(); {
		s := it.Value().(*Section)
		if s.name != "" {
			return false
		}
		if !s.entries.IsEmpty() {
			return false
		}
	}

	return true
}

// Map convert ini to map[string]map[string]any
func (ini *Ini) Map() map[string]map[string]any {
	m := make(map[string]map[string]any, ini.sections.Len())
	for it := ini.sections.Iterator(); it.Next(); {
		sec := it.Value().(*Section)
		m[sec.name] = sec.Map()
	}
	return m
}

// StringsMap convert ini to map[string]map[string][]string
func (ini *Ini) StringsMap() map[string]map[string][]string {
	m := make(map[string]map[string][]string, ini.sections.Len())
	for it := ini.sections.Iterator(); it.Next(); {
		sec := it.Value().(*Section)
		m[sec.name] = sec.StringsMap()
	}
	return m
}

// StringMap convert ini to map[string]map[string]string
func (ini *Ini) StringMap() map[string]map[string]string {
	m := make(map[string]map[string]string, ini.sections.Len())
	for it := ini.sections.Iterator(); it.Next(); {
		sec := it.Value().(*Section)
		m[sec.name] = sec.StringMap()
	}
	return m
}

// SectionNames returns the section array
func (ini *Ini) SectionNames() []string {
	ss := make([]string, ini.sections.Len())
	for it := ini.sections.Iterator(); it.Next(); {
		ss = append(ss, it.Key().(string))
	}
	return ss
}

// Sections returns the section array
func (ini *Ini) Sections() []*Section {
	ss := make([]*Section, ini.sections.Len())
	for it := ini.sections.Iterator(); it.Next(); {
		ss = append(ss, it.Value().(*Section))
	}
	return ss
}

// Section return a section with the specified name.
// if the section does not exists, create and add it to ini.
func (ini *Ini) Section(name string) (sec *Section) {
	if sec = ini.GetSection(name); sec == nil {
		sec = ini.NewSection(name)
		ini.AddSection(sec)
	}
	return
}

// GetSection return a section with the specified name or nil if section not exists
func (ini *Ini) GetSection(name string) *Section {
	if sec, ok := ini.sections.Get(name); ok {
		return sec.(*Section)
	}
	return nil
}

// NewSection create a section to INI, overwrite existing section
func (ini *Ini) NewSection(name string, comments ...string) *Section {
	section := NewSection(name, comments...)
	ini.sections.Set(section.name, section)
	return section
}

// AddSection add a section to INI, overwrite existing section
func (ini *Ini) AddSection(section *Section) {
	ini.sections.Set(section.name, section)
}

// RemoveSection remove a section from INI
func (ini *Ini) RemoveSection(name string) *Section {
	if name == "" {
		sec, _ := ini.sections.Set("", NewSection(""))
		if sec == nil {
			return nil
		}
		return sec.(*Section)
	}

	sec, _ := ini.sections.Delete(name)
	if sec == nil {
		return nil
	}
	return sec.(*Section)
}

// emptySection empty section for internal use
var emptySection = NewSection("")

// esection get a section from ini, or a empty section if it does not exists
func (ini *Ini) esection(name string) *Section {
	if sec, ok := ini.sections.Get(name); ok {
		return sec.(*Section)
	}
	return emptySection
}

// Set set a key/value entry to the section sec
func (ini *Ini) Set(sec, key, value string, comments ...string) *Entry {
	return ini.Section(sec).Set(key, value, comments...)
}

// Get get a value of the key from the section sec
func (ini *Ini) Get(sec, key string) string {
	return ini.esection(sec).Get(key)
}

// GetString get a string value of the key from the section sec
// if not found, returns the default defs[0] string value
func (ini *Ini) GetString(sec, key string, defs ...string) string {
	return ini.esection(sec).GetString(key, defs...)
}

// GetInt get a int value of the key from the section sec
// if not found or convert error, returns the default defs[0] int value
func (ini *Ini) GetInt(sec, key string, defs ...int) int {
	return ini.esection(sec).GetInt(key, defs...)
}

// GetInt64 get a int64 value of the key from the section sec
// if not found or convert error, returns the default defs[0] int64 value
func (ini *Ini) GetInt64(sec, key string, defs ...int64) int64 {
	return ini.esection(sec).GetInt64(key, defs...)
}

// GetFloat get a float value of the key from the section sec
// if not found or convert error, returns the default defs[0] float value
func (ini *Ini) GetFloat(sec, key string, defs ...float64) float64 {
	return ini.esection(sec).GetFloat(key, defs...)
}

// GetBool get a bool value of the key from the section sec
// if not found or convert error, returns the default defs[0] bool value
func (ini *Ini) GetBool(sec, key string, defs ...bool) bool {
	return ini.esection(sec).GetBool(key, defs...)
}

// GetSize get a int64 size value of the key from the section sec
// if not found or convert error, returns the default defs[0] int value
func (ini *Ini) GetSize(sec, key string, defs ...int64) int64 {
	return ini.esection(sec).GetSize(key, defs...)
}

// GetDuration get a time.Duration value of the key from the section sec
// if not found or convert error, returns the default defs[0] Duration value
func (ini *Ini) GetDuration(sec, key string, defs ...time.Duration) time.Duration {
	return ini.esection(sec).GetDuration(key, defs...)
}

// LoadFile load INI from file
func (ini *Ini) LoadFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return ini.LoadData(f)
}

// LoadFileFS load INI from file
func (ini *Ini) LoadFileFS(fsys fs.FS, path string) error {
	f, err := fsys.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return ini.LoadData(f)
}

// LoadData load INI from io.Reader
func (ini *Ini) LoadData(r io.Reader) (err error) {
	if r, err = iox.SkipBOM(r); err != nil {
		return err
	}

	lineContinue := false         // line continue flag
	section := ini.GetSection("") // last section
	var comments []string         // last comments
	var key string                // last key
	var val bytes.Buffer          // last value

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		bs := bytes.TrimSpace(scanner.Bytes())

		// line continuation
		if lineContinue {
			if bye.EndsWithByte(bs, '\\') {
				if len(bs) == 1 {
					// a single '\\' line means EOL
					if bye.StartsWithByte(val.Bytes(), '"') {
						qs := quote(ini.EOL)
						val.WriteString(qs[1 : len(qs)-1])
					} else {
						val.WriteString(ini.EOL)
					}
				} else {
					val.Write(bs[:len(bs)-1])
				}
			} else {
				val.Write(bs)

				s, err := unquote(val.String())
				if err != nil {
					return err
				}

				section.Add(key, s, comments...)

				key = ""
				val.Reset()
				comments = nil
				lineContinue = false
			}
			continue
		}

		// empty line
		if len(bs) == 0 {
			if len(comments) > 0 {
				if ini.IsEmpty() {
					global := ini.GetSection("") // global section / no name section
					if len(global.comments) == 0 {
						global.comments = comments
					} else {
						global.comments = append(global.comments, "")
						global.comments = append(global.comments, comments...)
					}
					comments = nil
				} else {
					comments = append(comments, "")
				}
			}
			continue
		}

		// first char
		c := bs[0]

		// comment
		if c == ';' || c == '#' {
			comments = append(comments, string(bs))
			continue
		}

		// section
		if c == '[' {
			if bs[len(bs)-1] != ']' {
				return fmt.Errorf("Invalid section: %s", string(bs))
			}

			sn := string(bs[1 : len(bs)-1])
			section = ini.GetSection(sn)
			if section == nil {
				section = ini.NewSection(sn, comments...)
			}
			comments = nil
			continue
		}

		// entry
		d := bytes.IndexByte(bs, '=')
		if d < 1 {
			return fmt.Errorf("Invalid entry: %s", string(bs))
		}

		// entry key
		k := string(bytes.TrimSpace(bs[:d]))
		if k == "" {
			return fmt.Errorf("Invalid entry: %s", string(bs))
		}

		// entry value
		v := bytes.TrimSpace(bs[d+1:])

		if bye.EndsWithByte(v, '\\') { // line continuation
			val.Write(v[:len(v)-1])
			key = k
			lineContinue = true
			continue
		}

		// add entry
		s, err := unquote(string(v))
		if err != nil {
			return err
		}

		if ini.Multiple {
			section.Add(k, s, comments...)
		} else {
			section.Set(k, s, comments...)
		}
		comments = nil
	}

	return scanner.Err()
}

// WriteFile write INI to the file
func (ini *Ini) WriteFile(path string) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(0660))
	if err != nil {
		return err
	}
	defer f.Close()

	return ini.WriteData(f)
}

// WriteData write INI to io.Writer
func (ini *Ini) WriteData(w io.Writer) (err error) {
	bw := bufio.NewWriter(w)

	for it := ini.sections.Iterator(); it.Next(); {
		sec := it.Value().(*Section)

		if err := sec.Write(bw, ini.EOL); err != nil {
			return err
		}
	}

	return bw.Flush()
}

// String write INI to the string
func (ini *Ini) String() string {
	sb := &strings.Builder{}
	ini.WriteData(sb) //nolint: errcheck
	return sb.String()
}

func isQuoted(s string) bool {
	return str.StartsWithByte(s, '"') && str.EndsWithByte(s, '"')
}

func unquote(s string) (string, error) {
	if isQuoted(s) {
		us, err := strconv.Unquote(s)
		if err != nil {
			err = errors.New("invalid quoted value: " + s)
		}
		return us, err
	}

	return s, nil
}

func isNotPrint(r rune) bool {
	return !unicode.IsPrint(r)
}

func needQuote(s string) bool {
	ss := []byte(s)
	bs := bytes.TrimSpace(ss)
	if len(bs) != len(ss) {
		return true
	}

	return bytes.IndexFunc(bs, isNotPrint) >= 0
}

func quote(s string) string {
	if needQuote(s) {
		return strconv.Quote(s)
	}
	return s
}
