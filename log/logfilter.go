package log

import (
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

// AndFilter a AND multiple filter
type AndFilter struct {
	Filters []Filter
}

// Reject reject event if any filter reject it
func (af *AndFilter) Reject(le *Event) bool {
	for _, f := range af.Filters {
		if f.Reject(le) {
			return true
		}
	}
	return false
}

// NewAndFilter create a AND multiple filter
func NewAndFilter(fs ...Filter) *AndFilter {
	return &AndFilter{Filters: fs}
}

// OrFilter a OR multiple filter
type OrFilter struct {
	Filters []Filter
}

// Reject accept event if any filter accept it
func (of *OrFilter) Reject(le *Event) bool {
	for _, f := range of.Filters {
		if !f.Reject(le) {
			return false
		}
	}
	return true
}

// NewOrFilter create a OR multiple filter
func NewOrFilter(fs ...Filter) *OrFilter {
	return &OrFilter{Filters: fs}
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

// ParseFilter parse filter expression "name:condition" to create a log filter
func ParseFilter(expr string) Filter {
	name, conf, ok := str.CutByte(expr, ':')
	if ok {
		return CreateFilter(name, conf)
	}
	return nil
}

func parseAndFilter(expr string) Filter {
	fs := []Filter{}

	ss := str.Fields(expr)
	for _, s := range ss {
		f := ParseFilter(s)
		if f != nil {
			fs = append(fs, f)
		}
	}

	if len(fs) < 1 {
		return nil
	}
	if len(fs) == 1 {
		return fs[0]
	}
	return &AndFilter{Filters: fs}
}

// NewLogFilter parse filter expression to create a log filter
func NewLogFilter(expr string) Filter {
	fs := []Filter{}

	ss := str.Split(expr, "||")
	for _, s := range ss {
		f := parseAndFilter(s)
		if f != nil {
			fs = append(fs, f)
		}
	}

	if len(fs) < 1 {
		return nil
	}
	if len(fs) == 1 {
		return fs[0]
	}
	return &OrFilter{Filters: fs}
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
