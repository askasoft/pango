package col

// MapEntry key/value entry
type MapEntry struct {
	key   interface{}
	Value interface{}
}

// Key returns the entry's key
func (e *MapEntry) Key() interface{} {
	return e.key
}
