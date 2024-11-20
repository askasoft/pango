package log

import (
	"strings"

	"github.com/askasoft/pango/str"
)

// Filter log filter
type Filter interface {
	Reject(le *Event) bool
}

// LevelFilter log level filter
type LevelFilter struct {
	Level Level
}

// Reject filter event by logger name
func (lf *LevelFilter) Reject(le *Event) bool {
	return lf.Level < le.Level
}

// NewLevelFilter create a level filter
func NewLevelFilter(lvl Level) *LevelFilter {
	return &LevelFilter{Level: lvl}
}

// NameFilter logger name filter
type NameFilter struct {
	Name string
}

// Reject filter event by logger name
func (nf *NameFilter) Reject(le *Event) bool {
	return nf.Name != le.Name
}

// NewNameFilter create a name filter
func NewNameFilter(name string) *NameFilter {
	return &NameFilter{Name: name}
}

// NameNotFilter logger name filter
type NameNotFilter struct {
	Name string
}

// Reject filter event by logger name
func (nnf *NameNotFilter) Reject(le *Event) bool {
	return nnf.Name == le.Name
}

// NewNameNotFilter create a name filter (not equal)
func NewNameNotFilter(name string) Filter {
	return &NameNotFilter{Name: name}
}

// MultiFilter a multiple filter
type MultiFilter struct {
	Filters []Filter
}

// Reject filter event by multiple filters
func (mf *MultiFilter) Reject(le *Event) bool {
	for _, f := range mf.Filters {
		if f.Reject(le) {
			return true
		}
	}
	return false
}

// NewMultiFilter create a multiple filter
func NewMultiFilter(fs ...Filter) *MultiFilter {
	return &MultiFilter{Filters: fs}
}

// FilterCreator filter create function
type FilterCreator func(s string) Filter

var filterCreators = make(map[string]FilterCreator)

// RegisterFilter register log filter type
func RegisterFilter(name string, fc FilterCreator) {
	filterCreators[name] = fc
}

// CreateFilter create a log filter by name and config
func CreateFilter(name string, conf string) Filter {
	if f, ok := filterCreators[name]; ok {
		return f(conf)
	}
	return nil
}

// NewLogFilter create a log filter by the configuration string 'c'
func NewLogFilter(c string) Filter {
	fs := []Filter{}
	ss := strings.Split(c, " ")
	for _, s := range ss {
		cs := strings.Split(s, ":")
		if len(cs) == 2 {
			f := CreateFilter(cs[0], cs[1])
			if f != nil {
				fs = append(fs, f)
			}
		}
	}

	if len(fs) < 1 {
		return nil
	}
	if len(fs) == 1 {
		return fs[0]
	}
	return &MultiFilter{Filters: fs}
}

type FilterSupport struct {
	Filter Filter // log filter
}

// SetFilter set the log filter
func (fs *FilterSupport) SetFilter(filter string) {
	fs.Filter = NewLogFilter(filter)
}

func (fs *FilterSupport) Reject(le *Event) bool {
	return fs.Filter != nil && fs.Filter.Reject(le)
}

func init() {
	RegisterFilter("name", func(s string) Filter {
		if str.StartsWithByte(s, '!') {
			return NewNameNotFilter(s[1:])
		}
		return NewNameFilter(s)
	})
	RegisterFilter("level", func(s string) Filter {
		return NewLevelFilter(ParseLevel(s))
	})
}
