package col

// MapEntry key/value entry
type MapEntry struct {
	key   interface{}
	Value interface{}
	entry *ListEntry
}

// Key returns the entry's key
func (e *MapEntry) Key() interface{} {
	return e.key
}

// Next returns a pointer to the next entry.
func (e *MapEntry) Next() *MapEntry {
	return toMapEntry(e.entry.Next())
}

// Prev returns a pointer to the previous entry.
func (e *MapEntry) Prev() *MapEntry {
	return toMapEntry(e.entry.Prev())
}

func toMapEntry(le *ListEntry) *MapEntry {
	if le == nil {
		return nil
	}
	return le.Value.(*MapEntry)
}
