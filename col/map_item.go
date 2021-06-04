package col

import "fmt"

// MapItem key-value item
type MapItem struct {
	key   interface{}
	Value interface{}
}

// Key returns the item's key
func (mi *MapItem) Key() interface{} {
	return mi.key
}

// String print the item to string
func (mi *MapItem) String() string {
	return fmt.Sprintf("%v => %v", mi.key, mi.Value)
}
