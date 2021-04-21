package ini

import (
	"strconv"

	"github.com/pandafw/pango/col"
)

// Entry ini entry
type Entry struct {
	Value    string
	Comments []string
}

// Section ini section
type Section struct {
	name     string          // Name for tihs section.
	comments []string        // Comment for this section.
	entries  *col.OrderedMap // Entries for this section.
}

// NewSection create a INI section
func NewSection(name string, comments ...string) *Section {
	return &Section{
		name:     name,
		comments: comments,
		entries:  col.NewOrderedMap(),
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

// Kvmap return the section's entries key/value map
func (sec *Section) Kvmap() map[string]interface{} {
	m := make(map[string]interface{}, sec.entries.Len())
	for e := sec.entries.Front(); e != nil; e = e.Next() {
		var v interface{}
		switch se := e.Value.(type) {
		case *col.List:
			v = sec.toStrings(se)
		case *Entry:
			v = se.Value
		}
		m[e.Key().(string)] = v
	}
	return m
}

// Add add a key/value entry to the section
func (sec *Section) Add(key string, value string, comments ...string) *Entry {
	e := &Entry{Value: value, Comments: comments}

	if v, ok := sec.entries.Get(key); ok {
		if l, ok := v.(*col.List); ok {
			l.PushBack(e)
			return e
		}
		l := col.NewList()
		l.PushBack(v)
		l.PushBack(e)
		sec.entries.Set(key, l)
		return e
	}

	return sec.Set(key, value, comments...)
}

// Set set a key/value entry to the section
func (sec *Section) Set(key string, value string, comments ...string) *Entry {
	e := &Entry{Value: value, Comments: comments}
	sec.entries.Set(key, e)
	return e
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
	if e != nil {
		return e.Value
	}
	if len(defs) > 0 {
		return defs[0]
	}
	return ""
}

// GetInt get a int value of the key from the section
// if not found, returns the default defs[0] int value
func (sec *Section) GetInt(key string, defs ...int) int {
	e := sec.GetEntry(key)
	if e != nil {
		if i, err := strconv.Atoi(e.Value); err != nil {
			return i
		}
	}
	if len(defs) > 0 {
		return defs[0]
	}
	return 0
}

// GetBool get a bool value of the key from the section
// if not found, returns the default defs[0] int value
func (sec *Section) GetBool(key string, defs ...bool) bool {
	e := sec.GetEntry(key)
	if e != nil {
		if b, err := strconv.ParseBool(e.Value); err != nil {
			return b
		}
	}
	if len(defs) > 0 {
		return defs[0]
	}
	return false
}

func (sec *Section) toStrings(l *col.List) []string {
	ss := make([]string, 0, l.Len())
	for e := l.Front(); e != nil; e = e.Next() {
		ss = append(ss, e.Value.(*Entry).Value)
	}
	return ss
}

// GetValues get the key's values from the section
func (sec *Section) GetValues(key string) []string {
	if v, ok := sec.entries.Get(key); ok {
		switch se := v.(type) {
		case *col.List:
			return sec.toStrings(se)
		case *Entry:
			return []string{se.Value}
		}
	}
	return nil
}

// GetEntry get the key's entry from the section
func (sec *Section) GetEntry(key string) *Entry {
	if v, ok := sec.entries.Get(key); ok {
		switch se := v.(type) {
		case *col.List:
			return se.Front().Value.(*Entry)
		case *Entry:
			return se
		}
	}
	return nil
}

// Clear clear the entries and comments
func (sec *Section) Clear() {
	sec.comments = nil
	sec.entries.Clear()
}
