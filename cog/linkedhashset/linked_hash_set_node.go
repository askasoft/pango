package linkedhashset

import "fmt"

// linkedSetNode is a node of a doublly-linked list.
type linkedSetNode[T any] struct {
	prev  *linkedSetNode[T]
	next  *linkedSetNode[T]
	value T
}

// String print the list item to string
func (ln *linkedSetNode[T]) String() string {
	return fmt.Sprintf("%v", ln.value)
}
