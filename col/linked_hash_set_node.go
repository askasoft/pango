package col

import "fmt"

// linkedSetNode is a node of a doublly-linked list.
type linkedSetNode struct {
	prev  *linkedSetNode
	next  *linkedSetNode
	value T
}

// String print the list item to string
func (ln *linkedSetNode) String() string {
	return fmt.Sprintf("%v", ln.value)
}
