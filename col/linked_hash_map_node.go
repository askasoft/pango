package col

import "fmt"

// LinkedMapNode is a node of a linked hash map.
type LinkedMapNode struct {
	prev  *LinkedMapNode
	next  *LinkedMapNode
	key   K
	value V
}

// Key returns the key
func (ln *LinkedMapNode) Key() K {
	return ln.key
}

// Value returns the key
func (ln *LinkedMapNode) Value() V {
	return ln.value
}

// SetValue sets the value
func (ln *LinkedMapNode) SetValue(v V) {
	ln.value = v
}

// String print the list item to string
func (ln *LinkedMapNode) String() string {
	return fmt.Sprintf("%v => %v", ln.key, ln.value)
}
