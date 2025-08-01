package ini

import (
	"io"
	"strings"
	"time"

	"github.com/askasoft/pango/bol"
	"github.com/askasoft/pango/cog/linkedhashmap"
	"github.com/askasoft/pango/cog/linkedlist"
	"github.com/askasoft/pango/iox"
	"github.com/askasoft/pango/num"
	"github.com/askasoft/pango/tmu"
)

// Entry ini entry
type Entry struct {
	Value    string
	Comments []string
}

type EntryList = linkedlist.LinkedList[*Entry]

type EntriesMap = linkedhashmap.LinkedHashMap[string, *EntryList]

// Section ini section
type Section struct {
	name     string     // Name for this section.
	comments []string   // Comment for this section.
	entries  EntriesMap // Entries for this section.
}

// Name return the section's name
func (sec *Section) Name() string {
	return sec.name
}

// Comments return the section's comment string array
func (sec *Section) Comments() []string {
	return sec.comments
}

// Keys return the section's key string array
func (sec *Section) Keys() []string {
	ks := make([]string, sec.entries.Len())
	for i, it := 0, sec.entries.Iterator(); it.Next(); {
		ks[i] = it.Key()
		i++
	}
	return ks
}

// StringMap return the section's entries key.(string)/value.(string) map
func (sec *Section) StringMap() map[string]string {
	m := make(map[string]string, sec.entries.Len())
	for it := sec.entries.Iterator(); it.Next(); {
		m[it.Key()] = it.Value().Head().Value
	}
	return m
}

// StringsMap return the section's entries key.(string)/value.([]string) map
func (sec *Section) StringsMap() map[string][]string {
	m := make(map[string][]string, sec.entries.Len())
	for it := sec.entries.Iterator(); it.Next(); {
		m[it.Key()] = sec.toStrings(it.Value())
	}
	return m
}

// Map return the section's entries key.(string)/value.(any) map
func (sec *Section) Map() map[string]any {
	m := make(map[string]any, sec.entries.Len())
	for it := sec.entries.Iterator(); it.Next(); {
		var v any
		es := it.Value()
		if es.Len() > 1 {
			v = sec.toStrings(es)
		} else {
			v = es.Head().Value
		}
		m[it.Key()] = v
	}
	return m
}

// Add add a key/value entry to the section
func (sec *Section) Add(key string, value string, comments ...string) *Entry {
	if es, ok := sec.entries.Get(key); ok {
		e := &Entry{Value: value, Comments: comments}
		es.Add(e)
		return e
	}

	return sec.Set(key, value, comments...)
}

// Set set a key/value entry to the section
func (sec *Section) Set(key string, value string, comments ...string) *Entry {
	e := &Entry{Value: value, Comments: comments}
	es := linkedlist.NewLinkedList(e)
	sec.entries.Set(key, es)
	return e
}

// Delete delete entries with key form the section
func (sec *Section) Delete(key string) {
	sec.entries.Remove(key)
}

// Remove remove a key/value entry from the section
func (sec *Section) Remove(key string, value string) (e *Entry) {
	if es, ok := sec.entries.Get(key); ok {
		for it := es.Iterator(); it.Next(); {
			en := it.Value()
			if en.Value == value {
				e = en
				it.Remove()
			}
		}
		return
	}

	return
}

// Get get a value of the key from the section
func (sec *Section) Get(key string) string {
	e := sec.GetEntry(key)
	if e != nil {
		return e.Value
	}
	return ""
}

// GetString get a string value of the key from the section
// if not found, returns the default defs[0] string value
func (sec *Section) GetString(key string, defs ...string) string {
	e := sec.GetEntry(key)
	if e != nil && e.Value != "" {
		return e.Value
	}
	if len(defs) > 0 {
		return defs[0]
	}
	return ""
}

// GetBool get a bool value of the key from the section
// if not found or convert error, returns the first non-false value from defs.
func (sec *Section) GetBool(key string, defs ...bool) bool {
	return bol.Atob(sec.GetString(key), defs...)
}

// GetInt get a int value of the key from the section
// if not found or convert error, returns the first non-zero value from defs.
func (sec *Section) GetInt(key string, defs ...int) int {
	return num.Atoi(sec.GetString(key), defs...)
}

// GetInt64 get a int64 value of the key from the section
// if not found or convert error, returns the first non-zero value from defs.
func (sec *Section) GetInt64(key string, defs ...int64) int64 {
	return num.Atol(sec.GetString(key), defs...)
}

// GetFloat get a float value of the key from the section
// if not found or convert error, returns the first non-zero value from defs.
func (sec *Section) GetFloat(key string, defs ...float64) float64 {
	return num.Atof(sec.GetString(key), defs...)
}

// GetSize get a int64 size value of the key from the section
// if not found or convert error, returns the first non-zero value from defs.
func (sec *Section) GetSize(key string, defs ...int64) int64 {
	e := sec.GetEntry(key)
	if e != nil && e.Value != "" {
		if sz, err := num.ParseSize(e.Value); err == nil {
			return sz
		}
	}

	for _, d := range defs {
		if d != 0 {
			return d
		}
	}
	return 0
}

// GetDuration get a time.Duration value of the key from the section.
// if not found or convert error, returns the first non-zero value from defs.
func (sec *Section) GetDuration(key string, defs ...time.Duration) time.Duration {
	return tmu.Atod(sec.GetString(key), defs...)
}

func (sec *Section) toStrings(es *EntryList) []string {
	ss := make([]string, 0, es.Len())
	for it := es.Iterator(); it.Next(); {
		ss = append(ss, it.Value().Value)
	}
	return ss
}

// GetValues get the key's values from the section
func (sec *Section) GetValues(key string) []string {
	if es, ok := sec.entries.Get(key); ok {
		return sec.toStrings(es)
	}
	return nil
}

// GetEntry get the key's entry from the section
func (sec *Section) GetEntry(key string) *Entry {
	if es, ok := sec.entries.Get(key); ok {
		return es.Head()
	}
	return nil
}

// GetEntries get the key's entries from the section
func (sec *Section) GetEntries(key string) []*Entry {
	if es, ok := sec.entries.Get(key); ok {
		return es.Values()
	}
	return nil
}

// Clear clear the entries and comments
func (sec *Section) Clear() {
	sec.comments = nil
	sec.entries.Clear()
}

// Copy copy entries from src section, overwrite existing entries
func (sec *Section) Copy(src *Section) {
	if len(src.comments) > 0 {
		sec.comments = src.comments
	}
	sec.entries.Copy(&src.entries)
}

// Merge merge entries from src section, existing entries will be merged
func (sec *Section) Merge(src *Section) {
	sec.comments = append(sec.comments, src.comments...)
	for it := src.entries.Iterator(); it.Next(); {
		if es, ok := sec.entries.Get(it.Key()); ok {
			es.AddCol(it.Value())
		} else {
			sec.entries.Set(it.Key(), it.Value())
		}
	}
}

// String write section to string
func (sec *Section) String() string {
	sb := &strings.Builder{}
	sec.Write(sb, iox.EOL) //nolint: errcheck
	return sb.String()
}

// Write output section to the writer
func (sec *Section) Write(w io.Writer, eol string) error {
	// comments
	if err := sec.writeComments(w, sec.comments, eol); err != nil {
		return err
	}

	// section name
	if err := sec.writeSectionName(w, eol); err != nil {
		return err
	}

	// section entries
	if err := sec.writeSectionEntries(w, eol); err != nil {
		return err
	}

	// blank line
	if _, err := io.WriteString(w, eol); err != nil {
		return err
	}

	return nil
}

func (sec *Section) writeComments(w io.Writer, comments []string, eol string) (err error) {
	for _, s := range comments {
		if _, err = io.WriteString(w, s); err != nil {
			return
		}
		if _, err = io.WriteString(w, eol); err != nil {
			return
		}
	}
	return
}

func (sec *Section) writeSectionName(w io.Writer, eol string) (err error) {
	if sec.name != "" {
		if _, err = w.Write([]byte{'['}); err != nil {
			return
		}
		if _, err = io.WriteString(w, sec.name); err != nil {
			return
		}
		if _, err = w.Write([]byte{']'}); err != nil {
			return
		}
	}
	_, err = io.WriteString(w, eol)
	return
}

func (sec *Section) writeSectionEntries(w io.Writer, eol string) (err error) {
	for sei := sec.entries.Iterator(); sei.Next(); {
		es := sei.Value()
		for it := es.Iterator(); it.Next(); {
			if err = sec.writeSectionEntry(w, sei.Key(), it.Value(), eol); err != nil {
				return err
			}
		}
	}
	return err
}

func (sec *Section) writeSectionEntry(w io.Writer, key string, ve *Entry, eol string) (err error) {
	if len(ve.Comments) > 0 {
		if _, err = io.WriteString(w, eol); err != nil {
			return
		}
		if err = sec.writeComments(w, ve.Comments, eol); err != nil {
			return
		}
	}

	if _, err = io.WriteString(w, key); err != nil {
		return
	}
	if _, err = io.WriteString(w, " = "); err != nil {
		return
	}
	if _, err = io.WriteString(w, quote(ve.Value)); err != nil {
		return
	}
	_, err = io.WriteString(w, eol)
	return
}
