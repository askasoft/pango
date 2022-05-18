package col

import "fmt"

// linkedListNode is a node of a doublly-linked list.
type linkedListNode struct {
	prev  *linkedListNode
	next  *linkedListNode
	value T
}

// String print the list item to string
func (ln *linkedListNode) String() string {
	return fmt.Sprintf("%v", ln.value)
}
