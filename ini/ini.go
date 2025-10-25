package ini

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/askasoft/pango/asg"
	"github.com/askasoft/pango/bye"
	"github.com/askasoft/pango/cog/linkedhashmap"
	"github.com/askasoft/pango/iox"
	"github.com/askasoft/pango/str"
)

type sections = linkedhashmap.LinkedHashMap[string, *Section]

// Ini INI file reader / writer
type Ini struct {
	sections sections // Parsed sections
	EOL      string   // End of Line
}

// NewIni create a Ini
func NewIni() *Ini {
	ini := &Ini{
		EOL: iox.EOL,
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
		sec := it.Value()
		if sec.name != "" {
			return false
		}
		if !sec.entries.IsEmpty() {
			return false
		}
	}

	return true
}

// StringMap convert ini to map[string]map[string]string
func (ini *Ini) StringMap() map[string]map[string]string {
	m := make(map[string]map[string]string, ini.sections.Len())
	for it := ini.sections.Iterator(); it.Next(); {
		sec := it.Value()
		m[sec.name] = sec.StringMap()
	}
	return m
}

// StringsMap convert ini to map[string]map[string][]string
func (ini *Ini) StringsMap() map[string]map[string][]string {
	m := make(map[string]map[string][]string, ini.sections.Len())
	for it := ini.sections.Iterator(); it.Next(); {
		sec := it.Value()
		m[sec.name] = sec.StringsMap()
	}
	return m
}

// Map convert ini to map[string]map[string]any
func (ini *Ini) Map() map[string]map[string]any {
	m := make(map[string]map[string]any, ini.sections.Len())
	for it := ini.sections.Iterator(); it.Next(); {
		sec := it.Value()
		m[sec.name] = sec.Map()
	}
	return m
}

// SectionNames returns the section array
func (ini *Ini) SectionNames() []string {
	return ini.sections.Keys()
}

// Sections returns the section array
func (ini *Ini) Sections() []*Section {
	return ini.sections.Values()
}

// Section return a section with the specified name.
// if the section does not exists, create and add it to ini.
func (ini *Ini) Section(name string) (sec *Section) {
	if sec = ini.GetSection(name); sec == nil {
		sec = ini.NewSection(name)
	}
	return
}

// GetSection return a section with the specified name or nil if section not exists
func (ini *Ini) GetSection(name string) *Section {
	if sec, ok := ini.sections.Get(name); ok {
		return sec
	}
	return nil
}

// NewSection create a section to INI, overwrite existing section
func (ini *Ini) NewSection(name string, comments ...string) *Section {
	sec := &Section{name: name, comments: comments}
	ini.sections.Set(sec.name, sec)
	return sec
}

// AddSection add a section to INI, overwrite existing section
func (ini *Ini) AddSection(section *Section) {
	ini.sections.Set(section.name, section)
}

// RemoveSection remove a section from INI
func (ini *Ini) RemoveSection(name string) (sec *Section) {
	if name == "" {
		sec, _ = ini.sections.Set("", &Section{})
		return
	}

	sec, _ = ini.sections.Remove(name)
	return
}

// emptySection empty section for internal use
var emptySection = &Section{}

// section get a section from ini, or a empty section if it does not exists
func (ini *Ini) section(name string) *Section {
	if sec, ok := ini.sections.Get(name); ok {
		return sec
	}
	return emptySection
}

// Add add a key/value entry to the section sec
func (ini *Ini) Add(sec, key, value string, comments ...string) *Entry {
	return ini.Section(sec).Add(key, value, comments...)
}

// Set set a key/value entry to the section sec
func (ini *Ini) Set(sec, key, value string, comments ...string) *Entry {
	return ini.Section(sec).Set(key, value, comments...)
}

// Get get a value of the key from the section sec
func (ini *Ini) Get(sec, key string) string {
	return ini.section(sec).Get(key)
}

// GetString get a string value of the key from the section sec
// if not found, returns the default defs[0] value.
func (ini *Ini) GetString(sec, key string, defs ...string) string {
	return ini.section(sec).GetString(key, defs...)
}

// GetInt get a int value of the key from the section sec
// if not found or convert error, returns the default defs[0] value.
func (ini *Ini) GetInt(sec, key string, defs ...int) int {
	return ini.section(sec).GetInt(key, defs...)
}

// GetInt64 get a int64 value of the key from the section sec
// if not found or convert error, returns the default defs[0] value.
func (ini *Ini) GetInt64(sec, key string, defs ...int64) int64 {
	return ini.section(sec).GetInt64(key, defs...)
}

// GetFloat get a float value of the key from the section sec
// if not found or convert error, returns the default defs[0] value.
func (ini *Ini) GetFloat(sec, key string, defs ...float64) float64 {
	return ini.section(sec).GetFloat(key, defs...)
}

// GetBool get a bool value of the key from the section sec
// if not found or convert error, returns the default defs[0] value.
func (ini *Ini) GetBool(sec, key string, defs ...bool) bool {
	return ini.section(sec).GetBool(key, defs...)
}

// GetSize get a int64 size value of the key from the section sec
// if not found or convert error, returns the default defs[0] value.
func (ini *Ini) GetSize(sec, key string, defs ...int64) int64 {
	return ini.section(sec).GetSize(key, defs...)
}

// GetEntry get a entry of key from the section sec
func (ini *Ini) GetEntry(sec, key string) *Entry {
	return ini.section(sec).GetEntry(key)
}

// GetDuration get a time.Duration value of the key from the section sec
// if not found or convert error, returns the default defs[0] value.
func (ini *Ini) GetDuration(sec, key string, defs ...time.Duration) time.Duration {
	return ini.section(sec).GetDuration(key, defs...)
}

// Copy copy section/entries from src ini, overwrite existing entries
func (ini *Ini) Copy(src *Ini) {
	for it := src.sections.Iterator(); it.Next(); {
		if ss, ok := ini.sections.Get(it.Key()); ok {
			ss.Copy(it.Value())
		} else {
			ini.sections.Set(it.Key(), it.Value())
		}
	}
}

// Merge merge entries from src ini, existing section/entries will be merged
func (ini *Ini) Merge(src *Ini) {
	for it := src.sections.Iterator(); it.Next(); {
		if ss, ok := ini.sections.Get(it.Key()); ok {
			ss.Merge(it.Value())
		} else {
			ini.sections.Set(it.Key(), it.Value())
		}
	}
}

// LoadFile load INI from file
// if multiple = true, allow multiple entry with same key
func (ini *Ini) LoadFile(path string, multiple ...bool) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return ini.LoadData(f, multiple...)
}

// LoadFileFS load INI from file
// if multiple = true, allow multiple entry with same key
func (ini *Ini) LoadFileFS(fsys fs.FS, path string, multiple ...bool) error {
	f, err := fsys.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return ini.LoadData(f, multiple...)
}

// LoadData load INI from io.Reader
// if multiple = true, allow multiple entry with same key
func (ini *Ini) LoadData(r io.Reader, multiple ...bool) (err error) {
	if r, err = iox.SkipBOM(r); err != nil {
		return err
	}

	multi := asg.First(multiple)
	section := ini.GetSection("") // last section
	linenum := 0                  // current line number
	linecon := 0                  // line continue counter

	var (
		buf      []byte       // line bytes
		key      string       // last key
		val      bytes.Buffer // last value
		comments []string     // last comments
	)

	scanner := bufio.NewScanner(r)
	for {
		eof := !scanner.Scan()
		if eof {
			buf = nil
		} else {
			buf = scanner.Bytes()
		}

		// line continuation
		if linecon > 0 {
			if bye.EndsWithByte(buf, '\\') {
				if len(buf) == 1 {
					// a single '\\' line means EOL
					if bye.StartsWithByte(val.Bytes(), '"') {
						qs := strconv.Quote(ini.EOL)
						val.WriteString(qs[1 : len(qs)-1])
					} else {
						val.WriteString(ini.EOL)
					}
				} else {
					val.Write(buf[:len(buf)-1])
				}
				linecon++
			} else {
				val.Write(buf)

				s, err := unquote(val.String())
				if err != nil {
					return fmt.Errorf("line %d: invalid quoted value", linenum)
				}

				section.Add(key, s, comments...)

				key = ""
				val.Reset()
				comments = nil

				linenum += linecon
				linecon = 0
			}
			continue
		}

		if eof {
			break
		}

		buf = bytes.TrimSpace(buf)
		linenum++

		// empty line
		if len(buf) == 0 {
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
		c := buf[0]

		// comment
		if c == ';' || c == '#' {
			comments = append(comments, string(buf))
			continue
		}

		// section
		if c == '[' {
			if buf[len(buf)-1] != ']' {
				return fmt.Errorf("line %d: invalid section", linenum)
			}

			sn := string(buf[1 : len(buf)-1])
			section = ini.GetSection(sn)
			if section == nil {
				section = ini.NewSection(sn, comments...)
			}
			comments = nil
			continue
		}

		// entry
		d := bytes.IndexByte(buf, '=')
		if d < 0 {
			return fmt.Errorf("line %d: invalid entry", linenum)
		}

		// entry key
		k := string(bytes.TrimSpace(buf[:d]))
		if k == "" {
			return fmt.Errorf("line %d: empty key", linenum)
		}

		// entry value
		v := bytes.TrimSpace(buf[d+1:])

		if bye.EndsWithByte(v, '\\') { // line continuation
			val.Write(v[:len(v)-1])
			key = k
			linecon++
			continue
		}

		// add entry
		s, err := unquote(string(v))
		if err != nil {
			return fmt.Errorf("line %d: invalid quoted value", linenum)
		}

		if multi {
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
		sec := it.Value()
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
		return strconv.Unquote(s)
	}
	return s, nil
}

func isNotPrint(r rune) bool {
	return !unicode.IsPrint(r)
}

func needQuote(s string) bool {
	ss := str.UnsafeBytes(s)
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
