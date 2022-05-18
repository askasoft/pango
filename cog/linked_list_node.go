package cog

import (
	"fmt"
)

// linkedListNode is a node of a doublly-linked list.
type linkedListNode[T any] struct {
	prev  *linkedListNode[T]
	next  *linkedListNode[T]
	value T
}

// String print the list item to string
func (ln *linkedListNode[T]) String() string {
	return fmt.Sprintf("%v", ln.value)
}
