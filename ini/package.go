package ini

import (
	"io"
	"io/fs"
	"time"
)

var _ini = NewIni()

// Default returns the default Ini instance used by the package-level functions.
func Default() *Ini {
	return _ini
}

// SetDefault set the default Ini instance used by the package-level functions.
func SetDefault(ini *Ini) {
	_ini = ini
}

// Clear clears the ini
func Clear() {
	_ini.Clear()
}

// IsEmpty returns true if the Ini has no entry
func IsEmpty() bool {
	return _ini.IsEmpty()
}

// StringMap convert ini to map[string]map[string]string
func StringMap() map[string]map[string]string {
	return _ini.StringMap()
}

// StringsMap convert ini to map[string]map[string][]string
func StringsMap() map[string]map[string][]string {
	return _ini.StringsMap()
}

// Map convert ini to map[string]map[string]any
func Map() map[string]map[string]any {
	return _ini.Map()
}

// SectionNames returns the section array
func SectionNames() []string {
	return _ini.SectionNames()
}

// Sections returns the section array
func Sections() []*Section {
	return _ini.Sections()
}

// SectionN return a section with the specified name.
// if the section does not exists, create and add it to ini.
func SectionN(name string) *Section {
	return _ini.Section(name)
}

// GetSection return a section with the specified name or nil if section not exists
func GetSection(name string) *Section {
	return _ini.GetSection(name)
}

// NewSection create a section to INI, overwrite existing section
func NewSection(name string, comments ...string) *Section {
	return _ini.NewSection(name, comments...)
}

// AddSection add a section to INI, overwrite existing section
func AddSection(sec *Section) {
	_ini.AddSection(sec)
}

// RemoveSection remove a section from INI
func RemoveSection(name string) *Section {
	return _ini.RemoveSection(name)
}

// Add add a key/value entry to the section sec
func Add(sec, key, value string, comments ...string) *Entry {
	return _ini.Add(sec, key, value, comments...)
}

// Set set a key/value entry to the section sec
func Set(sec, key, value string, comments ...string) *Entry {
	return _ini.Set(sec, key, value, comments...)
}

// Get get a value of the key from the section sec
func Get(sec, key string) string {
	return _ini.Get(sec, key)
}

// GetString get a string value of the key from the section sec
// if not found, returns the default defs[0] string value
func GetString(sec, key string, defs ...string) string {
	return _ini.GetString(sec, key, defs...)
}

// GetInt get a int value of the key from the section sec
// if not found or convert error, returns the first non-zero value from defs.
func GetInt(sec, key string, defs ...int) int {
	return _ini.GetInt(sec, key, defs...)
}

// GetInt64 get a int64 value of the key from the section sec
// if not found or convert error, returns the first non-zero value from defs.
func GetInt64(sec, key string, defs ...int64) int64 {
	return _ini.GetInt64(sec, key, defs...)
}

// GetFloat get a float value of the key from the section sec
// if not found or convert error, returns the first non-zero value from defs.
func GetFloat(sec, key string, defs ...float64) float64 {
	return _ini.GetFloat(sec, key, defs...)
}

// GetBool get a bool value of the key from the section sec
// if not found or convert error, returns the first non-false value from defs.
func GetBool(sec, key string, defs ...bool) bool {
	return _ini.GetBool(sec, key, defs...)
}

// GetSize get a int64 size value of the key from the section sec
// if not found or convert error, returns the first non-zero value from defs.
func GetSize(sec, key string, defs ...int64) int64 {
	return _ini.GetSize(sec, key, defs...)
}

// GetDuration get a time.Duration value of the key from the section sec
// if not found or convert error, returns the first non-zero value from defs.
func GetDuration(sec, key string, defs ...time.Duration) time.Duration {
	return _ini.GetDuration(sec, key, defs...)
}

// Copy copy section/entries from src ini, overwrite existing entries
func Copy(src *Ini) {
	_ini.Copy(src)
}

// Merge merge entries from src ini, existing section/entries will be merged
func Merge(src *Ini) {
	_ini.Merge(src)
}

// LoadFile load INI from file
// if multiple = true, allow multiple entry with same key
func LoadFile(path string, multiple ...bool) error {
	return _ini.LoadFile(path, multiple...)
}

// LoadFileFS load INI from file
// if multiple = true, allow multiple entry with same key
func LoadFileFS(fsys fs.FS, path string, multiple ...bool) error {
	return _ini.LoadFileFS(fsys, path, multiple...)
}

// LoadData load INI from io.Reader
// if multiple = true, allow multiple entry with same key
func LoadData(r io.Reader, multiple ...bool) error {
	return _ini.LoadData(r, multiple...)
}

// WriteFile write INI to the file
func WriteFile(path string) error {
	return _ini.WriteFile(path)
}

// WriteData write INI to io.Writer
func WriteData(w io.Writer) error {
	return _ini.WriteData(w)
}

// String write INI to the string
func String() string {
	return _ini.String()
}
