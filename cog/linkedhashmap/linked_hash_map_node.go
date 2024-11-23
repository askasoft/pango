package linkedhashmap

import "fmt"

// LinkedMapNode is a node of a linked hash map.
type LinkedMapNode[K any, V any] struct {
	prev  *LinkedMapNode[K, V]
	next  *LinkedMapNode[K, V]
	key   K
	value V
}

// Key returns the key
func (ln *LinkedMapNode[K, V]) Key() K {
	return ln.key
}

// Value returns the key
func (ln *LinkedMapNode[K, V]) Value() V {
	return ln.value
}

// SetValue sets the value
func (ln *LinkedMapNode[K, V]) SetValue(v V) {
	ln.value = v
}

// String print the list item to string
func (ln *LinkedMapNode[K, V]) String() string {
	return fmt.Sprintf("%v => %v", ln.key, ln.value)
}
