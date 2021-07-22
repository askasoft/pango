package col

import "fmt"

// OrderedMapItem key/value item
type OrderedMapItem struct {
	// Next and previous pointers in the doubly-linked list of items.
	// To simplify the implementation, internally a list l is implemented
	// as a ring, such that &l.root is both the next item of the last
	// list item (l.Back()) and the previous item of the first list
	// item (l.Front()).
	next, prev *OrderedMapItem

	// The ordered map to which this item belongs.
	omap *OrderedMap

	key   interface{}
	Value interface{}
}

// Key returns the item's key
func (mi *OrderedMapItem) Key() interface{} {
	return mi.key
}

// Next returns the next list item or nil.
func (mi *OrderedMapItem) Next() *OrderedMapItem {
	if ni := mi.next; mi.omap != nil && ni != &mi.omap.root {
		return ni
	}
	return nil
}

// Prev returns the previous list item or nil.
func (mi *OrderedMapItem) Prev() *OrderedMapItem {
	if pi := mi.prev; mi.omap != nil && pi != &mi.omap.root {
		return pi
	}
	return nil
}

// Remove remove this item from the map
func (mi *OrderedMapItem) Remove() {
	if mi.omap == nil {
		return
	}

	delete(mi.omap.hash, mi.key)

	mi.prev.next = mi.next
	mi.next.prev = mi.prev

	// avoid memory leaks
	mi.next = nil
	mi.prev = nil
	mi.omap = nil
}

// String print the item to string
func (mi *OrderedMapItem) String() string {
	return fmt.Sprintf("%v => %v", mi.key, mi.Value)
}
