//go:build go1.18
// +build go1.18

package ini

import (
	"bufio"
	"strconv"
	"strings"
	"time"

	"github.com/pandafw/pango/cog"
	"github.com/pandafw/pango/iox"
	"github.com/pandafw/pango/num"
)

// Entry ini entry
type Entry struct {
	Value    string
	Comments []string
}

type EntryList = cog.LinkedList[*Entry]

type EntriesMap = cog.LinkedHashMap[string, *EntryList]

// Section ini section
type Section struct {
	name     string     // Name for this section.
	comments []string   // Comment for this section.
	entries  EntriesMap // Entries for this section.
}

// NewSection create a INI section
func NewSection(name string, comments ...string) *Section {
	return &Section{
		name:     name,
		comments: comments,
	}
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
	es := cog.NewLinkedList(e)
	sec.entries.Set(key, es)
	return e
}

// Delete delete entries with key form the section
func (sec *Section) Delete(key string) {
	sec.entries.Delete(key)
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

// GetInt get a int value of the key from the section
// if not found or convert error, returns the default defs[0] int value
func (sec *Section) GetInt(key string, defs ...int) int {
	e := sec.GetEntry(key)
	if e != nil && e.Value != "" {
		if i, err := strconv.ParseInt(e.Value, 0, strconv.IntSize); err == nil {
			return int(i)
		}
	}
	if len(defs) > 0 {
		return defs[0]
	}
	return 0
}

// GetInt64 get a int64 value of the key from the section
// if not found or convert error, returns the default defs[0] int64 value
func (sec *Section) GetInt64(key string, defs ...int64) int64 {
	e := sec.GetEntry(key)
	if e != nil && e.Value != "" {
		if i, err := strconv.ParseInt(e.Value, 0, 64); err == nil {
			return i
		}
	}
	if len(defs) > 0 {
		return defs[0]
	}
	return 0
}

// GetFloat get a float value of the key from the section
// if not found or convert error, returns the default defs[0] float value
func (sec *Section) GetFloat(key string, defs ...float64) float64 {
	e := sec.GetEntry(key)
	if e != nil && e.Value != "" {
		if f, err := strconv.ParseFloat(e.Value, 64); err == nil {
			return f
		}
	}
	if len(defs) > 0 {
		return defs[0]
	}
	return 0
}

// GetBool get a bool value of the key from the section
// if not found or convert error, returns the default defs[0] bool value
func (sec *Section) GetBool(key string, defs ...bool) bool {
	e := sec.GetEntry(key)
	if e != nil && e.Value != "" {
		if b, err := strconv.ParseBool(e.Value); err == nil {
			return b
		}
	}
	if len(defs) > 0 {
		return defs[0]
	}
	return false
}

// GetSize get a int64 size value of the key from the section
// if not found or convert error, returns the default defs[0] int value
func (sec *Section) GetSize(key string, defs ...int64) int64 {
	e := sec.GetEntry(key)
	if e != nil && e.Value != "" {
		if sz, err := num.ParseSize(e.Value); err == nil {
			return sz
		}
	}
	if len(defs) > 0 {
		return defs[0]
	}
	return 0
}

// GetDuration get a time.Duration value of the key from the section
// if not found or convert error, returns the default defs[0] Duration value
func (sec *Section) GetDuration(key string, defs ...time.Duration) time.Duration {
	e := sec.GetEntry(key)
	if e != nil && e.Value != "" {
		if d, err := time.ParseDuration(e.Value); err == nil {
			return d
		}
	}
	if len(defs) > 0 {
		return defs[0]
	}
	return 0
}

func (sec *Section) toStrings(es *cog.LinkedList[*Entry]) []string {
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

// Copy copy entries from src section, overrite existing entries
func (sec *Section) Copy(src *Section) {
	if len(src.comments) > 0 {
		sec.comments = src.comments
	}
	sec.entries.SetAll(&src.entries)
}

// Merge merge entries from src section
func (sec *Section) Merge(src *Section, multiple bool) {
	sec.comments = append(sec.comments, src.comments...)
	for it := src.entries.Iterator(); it.Next(); {
		if multiple {
			if es, ok := sec.entries.Get(it.Key()); ok {
				es.AddAll(it.Value())
				continue
			}
		}
		sec.entries.Set(it.Key(), it.Value())
	}
}

// String write section to string
func (sec *Section) String() string {
	sb := &strings.Builder{}
	bw := bufio.NewWriter(sb)
	sec.Write(bw, iox.EOL) //nolint: errcheck
	bw.Flush()
	return sb.String()
}

// Write output section to the writer
func (sec *Section) Write(bw *bufio.Writer, eol string) error {
	// comments
	if err := sec.writeComments(bw, sec.comments, eol); err != nil {
		return err
	}

	// section name
	if err := sec.writeSectionName(bw, eol); err != nil {
		return err
	}

	// section entries
	if err := sec.writeSectionEntries(bw, eol); err != nil {
		return err
	}

	// blank line
	if _, err := bw.WriteString(eol); err != nil {
		return err
	}

	return nil
}

func (sec *Section) writeComments(bw *bufio.Writer, comments []string, eol string) (err error) {
	for _, s := range comments {
		if _, err = bw.WriteString(s); err != nil {
			return
		}
		if _, err = bw.WriteString(eol); err != nil {
			return
		}
	}
	return
}

func (sec *Section) writeSectionName(bw *bufio.Writer, eol string) (err error) {
	if sec.name != "" {
		if err = bw.WriteByte('['); err != nil {
			return
		}
		if _, err = bw.WriteString(sec.name); err != nil {
			return
		}
		if err = bw.WriteByte(']'); err != nil {
			return
		}
	}
	_, err = bw.WriteString(eol)
	return
}

func (sec *Section) writeSectionEntries(bw *bufio.Writer, eol string) (err error) {
	for sei := sec.entries.Iterator(); sei.Next(); {
		es := sei.Value()
		for it := es.Iterator(); it.Next(); {
			if err = sec.writeSectionEntry(bw, sei.Key(), it.Value(), eol); err != nil {
				return err
			}
		}
	}
	return err
}

func (sec *Section) writeSectionEntry(bw *bufio.Writer, key string, ve *Entry, eol string) (err error) {
	if len(ve.Comments) > 0 {
		if _, err = bw.WriteString(eol); err != nil {
			return
		}
		if err = sec.writeComments(bw, ve.Comments, eol); err != nil {
			return
		}
	}

	if _, err = bw.WriteString(key); err != nil {
		return
	}
	if _, err = bw.WriteString(" = "); err != nil {
		return
	}
	if _, err = bw.WriteString(quote(ve.Value)); err != nil {
		return
	}
	_, err = bw.WriteString(eol)
	return
}
