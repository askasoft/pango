package log

// Filter log filter
type Filter interface {
	Reject(le *Event) bool
}

// LevelFilter log level filter
type LevelFilter struct {
	Level int
}

// Reject filter event by logger name
func (lf *LevelFilter) Reject(le *Event) bool {
	return lf.Level < le.Level
}

// NewLevelFilter create a level filter
func NewLevelFilter(level int) *LevelFilter {
	return &LevelFilter{Level: level}
}

// NameFilter logger name filter
type NameFilter struct {
	Name string
}

// Reject filter event by logger name
func (lf *NameFilter) Reject(le *Event) bool {
	return lf.Name != le.Logger.GetName()
}

// NewNameFilter create a name filter
func NewNameFilter(name string) *NameFilter {
	return &NameFilter{Name: name}
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
