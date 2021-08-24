package col

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pandafw/pango/cmp"
	"github.com/pandafw/pango/iox"
)

// NewTreeSet creates a new TreeSet.
// Example: NewTreeSet(cmp.CompareString, "v1", "v2")
func NewTreeSet(compare cmp.Compare, vs ...interface{}) *TreeSet {
	ts := &TreeSet{compare: compare}
	ts.Add(vs...)
	return ts
}

// TreeSet implements an tree set that keeps the compare order of keys.
// The zero value for TreeSet is an empty set ready to use.
//
// https://en.wikipedia.org/wiki/Red%E2%80%93black_tree
//
// To iterate over a tree set (where ts is a *TreeSet):
//	for it := ts.Iterator(); it.Next(); {
//		// do something with it.Value()
//	}
//
type TreeSet struct {
	len     int
	root    *treeSetNode
	compare cmp.Compare
}

//-----------------------------------------------------------
// implements Collection interface

// Len returns the length of the tree set.
func (ts *TreeSet) Len() int {
	return ts.len
}

// IsEmpty returns true if the set has no items
func (ts *TreeSet) IsEmpty() bool {
	return ts.len == 0
}

// Clear clears the set
func (ts *TreeSet) Clear() {
	ts.len = 0
	ts.root = nil
}

// Add adds all items of vs and returns the last added item.
func (ts *TreeSet) Add(vs ...interface{}) {
	for _, v := range vs {
		ts.add(v)
	}
	return
}

// AddAll adds all items of another collection
func (ts *TreeSet) AddAll(ac Collection) {
	if ac.IsEmpty() || ts == ac {
		return
	}

	if ic, ok := ac.(Iterable); ok {
		it := ic.Iterator()
		for it.Next() {
			ts.add(it.Value())
		}
		return
	}

	ts.Add(ac.Values()...)
}

// Delete delete all items with associated value v of vs
func (ts *TreeSet) Delete(vs ...interface{}) {
	if ts.IsEmpty() {
		return
	}

	for _, v := range vs {
		if tn := ts.lookup(v); tn != nil {
			ts.deleteNode(tn)
		}
	}
}

// DeleteAll delete all of this collection's elements that are also contained in the specified collection
func (ts *TreeSet) DeleteAll(ac Collection) {
	if ts.IsEmpty() || ac.IsEmpty() {
		return
	}

	if ts == ac {
		ts.Clear()
		return
	}

	if ic, ok := ac.(Iterable); ok {
		it := ic.Iterator()
		for it.Next() {
			if tn := ts.lookup(it.Value()); tn != nil {
				ts.deleteNode(tn)
			}
		}
		return
	}

	ts.Delete(ac.Values()...)
}

// Contains Test to see if the collection contains all items of vs
func (ts *TreeSet) Contains(vs ...interface{}) bool {
	if len(vs) == 0 {
		return true
	}

	if ts.IsEmpty() {
		return false
	}

	for _, v := range vs {
		if tn := ts.lookup(v); tn == nil {
			return false
		}
	}
	return true
}

// ContainsAll Test to see if the collection contains all items of another collection
func (ts *TreeSet) ContainsAll(ac Collection) bool {
	if ts == ac || ac.IsEmpty() {
		return true
	}

	if ts.IsEmpty() {
		return false
	}

	if ic, ok := ac.(Iterable); ok {
		it := ic.Iterator()
		for it.Next() {
			if tn := ts.lookup(it.Value()); tn == nil {
				return false
			}
		}
		return true
	}

	return ts.Contains(ac.Values()...)
}

// Retain Retains only the elements in this collection that are contained in the argument array vs.
func (ts *TreeSet) Retain(vs ...interface{}) {
	if ts.IsEmpty() || len(vs) == 0 {
		return
	}

	ts.RetainAll(NewArrayList(vs...))
}

// RetainAll Retains only the elements in this collection that are contained in the specified collection.
func (ts *TreeSet) RetainAll(ac Collection) {
	if ts.IsEmpty() || ac.IsEmpty() || ts == ac {
		return
	}

	for tn := ts.front(); tn != nil; tn = tn.next() {
		if !ac.Contains(tn.value) {
			ts.deleteNode(tn)
		}
	}
}

// Values returns the value slice
func (ts *TreeSet) Values() []interface{} {
	vs := make([]interface{}, ts.len)
	for i, n := 0, ts.front(); n != nil; i, n = i+1, n.next() {
		vs[i] = n.value
	}
	return vs
}

// Each call f for each item in the set
func (ts *TreeSet) Each(f func(v interface{})) {
	for tn := ts.front(); tn != nil; tn = tn.next() {
		f(tn.value)
	}
}

// ReverseEach call f for each item in the set with reverse order
func (ts *TreeSet) ReverseEach(f func(v interface{})) {
	for tn := ts.back(); tn != nil; tn = tn.prev() {
		f(tn.value)
	}
}

// Iterator returns a iterator for the set
func (ts *TreeSet) Iterator() Iterator {
	return &treeSetIterator{tree: ts}
}

//----------------------------------------------------------------

// Front returns the first item of set ts or nil if the set is empty.
func (ts *TreeSet) Front() (v interface{}) {
	tn := ts.front()
	if tn != nil {
		v = tn.value
	}
	return
}

// Back returns the last item of set ts or nil if the set is empty.
func (ts *TreeSet) Back() (v interface{}) {
	tn := ts.back()
	if tn != nil {
		v = tn.value
	}
	return
}

// PopFront remove the first item of set.
func (ts *TreeSet) PopFront() (v interface{}) {
	tn := ts.front()
	if tn != nil {
		v = tn.value
		ts.deleteNode(tn)
	}
	return
}

// PopBack remove the last item of set.
func (ts *TreeSet) PopBack() (v interface{}) {
	tn := ts.back()
	if tn != nil {
		v = tn.value
		ts.deleteNode(tn)
	}
	return
}

// Floor Finds floor node of the input key, return the floor node's value or nil if no floor is found.
//
// Floor node is defined as the largest node that is smaller than or equal to the given node.
// A floor node may not be found, either because the tree is empty, or because
// all nodes in the tree are larger than the given node.
//
// key should adhere to the comparator's type assertion, otherwise method panics.
func (ts *TreeSet) Floor(v interface{}) interface{} {
	tn := ts.floor(v)
	if tn != nil {
		return tn.value
	}
	return nil
}

// Ceiling finds ceiling node of the input key, return the ceiling node's value or nil if no ceiling is found.
//
// Ceiling node is defined as the smallest node that is larger than or equal to the given node.
// A ceiling node may not be found, either because the tree is empty, or because
// all nodes in the tree are smaller than the given node.
//
// key should adhere to the comparator's type assertion, otherwise method panics.
func (ts *TreeSet) Ceiling(v interface{}) interface{} {
	tn := ts.floor(v)
	if tn != nil {
		return tn.value
	}
	return nil
}

// String print set to string
func (ts *TreeSet) String() string {
	bs, _ := json.Marshal(ts)
	return string(bs)
}

// Graph return the set's graph
func (ts *TreeSet) Graph() string {
	return ts.root.graph(0)
}

//-----------------------------------------------------
func (ts *TreeSet) checkItemIndex(index int) int {
	len := ts.len
	if index >= len || index < -len {
		panic(fmt.Sprintf("TreeSet out of bounds: index=%d, len=%d", index, len))
	}

	if index < 0 {
		index += len
	}
	return index
}

func (ts *TreeSet) checkSizeIndex(index int) int {
	len := ts.len
	if index > len || index < -len {
		panic(fmt.Sprintf("TreeSet out of bounds: index=%d, len=%d", index, len))
	}

	if index < 0 {
		index += len
	}
	return index
}

func (ts *TreeSet) setValue(tn *treeSetNode, v interface{}) *treeSetNode {
	if tn.value == v {
		return tn
	}

	// compare equals, just set the node's value
	if ts.compare(v, tn.value) == 0 {
		tn.value = v
		return tn
	}

	// delete and insert again
	ts.deleteNode(tn)
	return ts.add(v)
}

// front returns a pointer to the minimum node.
func (ts *TreeSet) front() *treeSetNode {
	tn := ts.root
	if tn != nil {
		for tn.left != nil {
			tn = tn.left
		}
	}
	return tn
}

// back returns a pointer to the maximum node.
func (ts *TreeSet) back() *treeSetNode {
	tn := ts.root
	if tn != nil {
		for tn.right != nil {
			tn = tn.right
		}
	}
	return tn
}

// floor Finds floor node of the input key, return the floor node or nil if no floor is found.
func (ts *TreeSet) floor(key interface{}) (floor *treeSetNode) {
	node := ts.root
	for node != nil {
		compare := ts.compare(key, node.value)
		switch {
		case compare == 0:
			return node
		case compare < 0:
			node = node.left
		case compare > 0:
			floor = node
			node = node.right
		}
	}
	return
}

// ceiling finds ceiling node of the input key, return the ceiling node or nil if no ceiling is found.
func (ts *TreeSet) ceiling(key interface{}) (ceiling *treeSetNode) {
	node := ts.root
	for node != nil {
		compare := ts.compare(key, node.value)
		switch {
		case compare == 0:
			return node
		case compare < 0:
			ceiling = node
			node = node.left
		case compare > 0:
			node = node.right
		}
	}
	return
}

// lookup looks for the given key, and returns the item associated with it,
// or nil if not found. The Node struct can then be used to iterate over the tree set
// from that point, either forward or backward.
func (ts *TreeSet) lookup(key interface{}) *treeSetNode {
	node := ts.root
	for node != nil {
		compare := ts.compare(key, node.value)
		switch {
		case compare == 0:
			return node
		case compare < 0:
			node = node.left
		case compare > 0:
			node = node.right
		}
	}
	return nil
}

// add adds the item, returns the item's node
// item should adhere to the comparator's type assertion, otherwise method panics.
func (ts *TreeSet) add(v interface{}) *treeSetNode {
	tn := ts.root
	if tn == nil {
		// Assert key is of comparator's type for initial tree
		ts.compare(v, v)

		ts.root = &treeSetNode{value: v, color: black}
		ts.len = 1
		return ts.root
	}

	cmp := 0
	parent := tn
	for tn != nil {
		parent = tn
		cmp = ts.compare(v, tn.value)
		switch {
		case cmp < 0:
			tn = tn.left
		case cmp > 0:
			tn = tn.right
		default:
			return tn
		}
	}

	tn = &treeSetNode{value: v, parent: parent}
	if cmp < 0 {
		parent.left = tn
	} else {
		parent.right = tn
	}

	ts.fixAfterInsertion(tn)

	ts.len++
	return tn
}

// delete delete the node from the tree by key,
// and returns what `Get` would have returned
// on that key prior to the call to `Delete`.
// key should adhere to the comparator's type assertion, otherwise method panics.
func (ts *TreeSet) delete(key interface{}) (ov interface{}, ok bool) {
	tn := ts.lookup(key)
	if tn == nil {
		return
	}

	ov, ok = tn.value, true
	ts.deleteNode(tn)
	return
}

// deleteNode delete the node p, returns the deleted node
// NOTE: if p has both left/right, p.next() will be deleted and returned
func (ts *TreeSet) deleteNode(p *treeSetNode) *treeSetNode {
	ts.len--

	// If strictly internal, copy successor's element to p and then make p point to successor.
	if p.left != nil && p.right != nil {
		s := p.next()

		p.value, s.value = s.value, p.value
		p = s
	} // p has 2 children

	// Start fixup at replacement node, if it exists.
	replacement := p.left
	if replacement == nil {
		replacement = p.right
	}

	if replacement != nil {
		// Link replacement to parent
		replacement.parent = p.parent
		if p.parent == nil {
			ts.root = replacement
		} else if p == p.parent.left {
			p.parent.left = replacement
		} else {
			p.parent.right = replacement
		}

		// Null out links so they are OK to use by fixAfterDeletion.
		p.left, p.right, p.parent = nil, nil, nil

		// Fix replacement
		if p.color == black {
			ts.fixAfterDeletion(replacement)
		}
	} else if p.parent == nil { // return if we are the only node.
		ts.root = nil
	} else { //  No children. Use self as phantom replacement and unlink.
		if p.color == black {
			ts.fixAfterDeletion(p)
		}

		if p.parent != nil {
			if p == p.parent.left {
				p.parent.left = nil
			} else if p == p.parent.right {
				p.parent.right = nil
			}
			p.parent = nil
		}
	}

	return p
}

func (ts *TreeSet) fixAfterInsertion(x *treeSetNode) {
	x.color = red

	for x != nil && x != ts.root && x.parent.color == red {
		if x.getParent() == x.getGrandParent().getLeft() {
			y := x.getGrandParent().getRight()
			if y.getColor() == red {
				x.getParent().setColor(black)
				y.setColor(black)
				x.getGrandParent().setColor(red)
				x = x.getGrandParent()
			} else {
				if x == x.getParent().getRight() {
					x = x.getParent()
					ts.rotateLeft(x)
				}
				x.getParent().setColor(black)
				x.getGrandParent().setColor(red)
				ts.rotateRight(x.getGrandParent())
			}
		} else {
			y := x.getGrandParent().getLeft()
			if y.getColor() == red {
				x.getParent().setColor(black)
				y.setColor(black)
				x.getGrandParent().setColor(red)
				x = x.getGrandParent()
			} else {
				if x == x.getParent().getLeft() {
					x = x.getParent()
					ts.rotateRight(x)
				}
				x.getParent().setColor(black)
				x.getGrandParent().setColor(red)
				ts.rotateLeft(x.getGrandParent())
			}
		}
	}
	ts.root.color = black
}

func (ts *TreeSet) fixAfterDeletion(x *treeSetNode) {
	for x != ts.root && x.getColor() == black {
		if x == x.getParent().getLeft() {
			sib := x.getParent().getRight()

			if sib.getColor() == red {
				sib.setColor(black)
				x.getParent().setColor(red)
				ts.rotateLeft(x.getParent())
				sib = x.getParent().getRight()
			}

			if sib.getLeft().getColor() == black && sib.getRight().getColor() == black {
				sib.setColor(red)
				x = x.getParent()
			} else {
				if sib.getRight().getColor() == black {
					sib.getLeft().setColor(black)
					sib.setColor(red)
					ts.rotateRight(sib)
					sib = x.getParent().getRight()
				}
				sib.setColor(x.getParent().getColor())
				x.getParent().setColor(black)
				sib.getRight().setColor(black)
				ts.rotateLeft(x.getParent())
				x = ts.root
			}
		} else { // symmetric
			sib := x.getParent().getLeft()

			if sib.getColor() == red {
				sib.setColor(black)
				x.getParent().setColor(red)
				ts.rotateRight(x.getParent())
				sib = x.getParent().getLeft()
			}

			if sib.getRight().getColor() == black && sib.getLeft().getColor() == black {
				sib.setColor(red)
				x = x.getParent()
			} else {
				if sib.getLeft().getColor() == black {
					sib.getRight().setColor(black)
					sib.setColor(red)
					ts.rotateLeft(sib)
					sib = x.getParent().getLeft()
				}
				sib.setColor(x.getParent().getColor())
				x.getParent().setColor(black)
				sib.getLeft().setColor(black)
				ts.rotateRight(x.getParent())
				x = ts.root
			}
		}
	}

	x.setColor(black)
}

func (ts *TreeSet) rotateLeft(p *treeSetNode) {
	if p != nil {
		r := p.right
		p.right = r.left
		if r.left != nil {
			r.left.parent = p
		}
		r.parent = p.parent
		if p.parent == nil {
			ts.root = r
		} else if p.parent.left == p {
			p.parent.left = r
		} else {
			p.parent.right = r
		}
		r.left = p
		p.parent = r
	}
}

func (ts *TreeSet) rotateRight(p *treeSetNode) {
	if p != nil {
		l := p.left
		p.left = l.right
		if l.right != nil {
			l.right.parent = p
		}
		l.parent = p.parent
		if p.parent == nil {
			ts.root = l
		} else if p.parent.right == p {
			p.parent.right = l
		} else {
			p.parent.left = l
		}
		l.right = p
		p.parent = l
	}
}

// debug return the set's graph (debug)
func (ts *TreeSet) debug() string {
	return ts.root.graph(tsColor | tsPoint)
}

//-----------------------------------------------------

// treeSetNode is a node of red-black tree
type treeSetNode struct {
	color  color
	left   *treeSetNode
	right  *treeSetNode
	parent *treeSetNode
	value  interface{}
}

func (tn *treeSetNode) getLeft() *treeSetNode {
	if tn != nil {
		return tn.left
	}
	return nil
}

func (tn *treeSetNode) getRight() *treeSetNode {
	if tn != nil {
		return tn.right
	}
	return nil
}

func (tn *treeSetNode) getParent() *treeSetNode {
	if tn != nil {
		return tn.parent
	}
	return nil
}

func (tn *treeSetNode) getGrandParent() *treeSetNode {
	if tn != nil && tn.parent != nil {
		return tn.parent.parent
	}
	return nil
}

func (tn *treeSetNode) getColor() color {
	if tn == nil {
		return black
	}
	return tn.color
}

func (tn *treeSetNode) setColor(c color) {
	if tn != nil {
		tn.color = c
	}
}

// prev returns the previous node or nil.
func (tn *treeSetNode) prev() *treeSetNode {
	if tn == nil {
		return nil
	}

	if tn.left != nil {
		p := tn.left
		for p.right != nil {
			p = p.right
		}
		return p
	}

	c := tn
	p := tn.parent
	for p != nil && c == p.left {
		c = p
		p = p.parent
	}
	return p
}

// next returns the next node or nil.
func (tn *treeSetNode) next() *treeSetNode {
	if tn == nil {
		return nil
	}

	if tn.right != nil {
		n := tn.right
		for n.left != nil {
			n = n.left
		}
		return n
	}

	c := tn
	n := tn.parent
	for n != nil && c == n.right {
		c = n
		n = n.parent
	}
	return n
}

// String print the set item to string
func (tn *treeSetNode) String() string {
	return fmt.Sprint(tn.value)
}

const (
	tsColor = 1 << iota
	tsPoint
)

func (tn *treeSetNode) graph(flag int) string {
	if tn == nil {
		return "(empty)"
	}

	sb := &strings.Builder{}
	tn.output(sb, "", true, flag)
	return sb.String()
}

func (tn *treeSetNode) output(sb *strings.Builder, prefix string, tail bool, flag int) {
	if tn.right != nil {
		newPrefix := prefix
		if tail {
			newPrefix += "│   "
		} else {
			newPrefix += "    "
		}
		tn.right.output(sb, newPrefix, false, flag)
	}

	sb.WriteString(prefix)
	if tail {
		sb.WriteString("└── ")
	} else {
		sb.WriteString("┌── ")
	}

	if flag&tsColor == tsColor {
		sb.WriteString(fmt.Sprintf("(%v) ", tn.color))
	}
	sb.WriteString(fmt.Sprint(tn.value))
	if flag&tsPoint == tsPoint {
		sb.WriteString(fmt.Sprintf(" (%p)", tn))
	}
	sb.WriteString(iox.EOL)

	if tn.left != nil {
		newPrefix := prefix
		if tail {
			newPrefix += "    "
		} else {
			newPrefix += "│   "
		}
		tn.left.output(sb, newPrefix, true, flag)
	}
}

//-----------------------------------------------------

// treeSetIterator a iterator for TreeSet
type treeSetIterator struct {
	tree    *TreeSet
	node    *treeSetNode
	removed bool
}

// Prev moves the iterator to the previous element and returns true if there was a previous element in the container.
// If Prev() returns true, then previous element's index and value can be retrieved by Index() and Value().
// Modifies the state of the iterator.
func (it *treeSetIterator) Prev() bool {
	if it.tree.IsEmpty() {
		return false
	}

	if it.node == nil {
		it.node = it.tree.back()
		it.removed = false
		return true
	}

	if it.removed {
		if it.node.left == nil {
			return false
		}
		it.node = it.node.left
		it.removed = false
		return true
	}

	if pi := it.node.prev(); pi != nil {
		it.node = pi
		it.removed = false
		return true
	}
	return false
}

// Next moves the iterator to the next element and returns true if there was a next element in the collection.
// If Next() returns true, then next element's value can be retrieved by Value().
// If Next() was called for the first time, then it will point the iterator to the first element if it exists.
// Modifies the state of the iterator.
func (it *treeSetIterator) Next() bool {
	if it.tree.IsEmpty() {
		return false
	}

	if it.node == nil {
		it.node = it.tree.front()
		it.removed = false
		return true
	}

	if it.removed {
		if it.node.right == nil {
			return false
		}
		it.node = it.node.right
		it.removed = false
		return true
	}

	if ni := it.node.next(); ni != nil {
		it.node = ni
		it.removed = false
		return true
	}
	return false
}

// Value returns the current element's value.
func (it *treeSetIterator) Value() interface{} {
	if it.node == nil {
		return nil
	}
	return it.node.value
}

// SetValue set the value to the item
func (it *treeSetIterator) SetValue(v interface{}) {
	if it.node == nil {
		return
	}

	if it.removed {
		// unlinked item
		it.node.value = v
		return
	}

	it.node = it.tree.setValue(it.node, v)
}

// Remove remove the current element
func (it *treeSetIterator) Remove() {
	if it.node == nil {
		return
	}

	if it.removed {
		panic("TreeSet can't remove a unlinked item")
	}

	p, n := it.node.prev(), it.node.next()
	d := it.tree.deleteNode(it.node)
	if d != it.node {
		n, it.node = it.node, d
	}

	// save prev/next for iterator
	it.node.left, it.node.right = p, n
	it.removed = true
}

// Reset resets the iterator to its initial state (one-before-first/one-after-last)
// Call Next()/Prev() to fetch the first/last element if any.
func (it *treeSetIterator) Reset() {
	it.node = nil
	it.removed = false
}

//-----------------------------------------------------------
// implements JSON Marshaller/Unmarshaller interface

func (ts *TreeSet) addJSONArrayItem(v interface{}) jsonArray {
	ts.Add(v)
	return ts
}

// MarshalJSON implements type json.Marshaler interface, so can be called in json.Marshal(ts)
func (ts *TreeSet) MarshalJSON() (res []byte, err error) {
	return jsonMarshalSet(ts)
}

// UnmarshalJSON implements type json.Unmarshaler interface, so can be called in json.Unmarshal(data, ts)
func (ts *TreeSet) UnmarshalJSON(data []byte) error {
	ts.Clear()
	ju := &jsonUnmarshaler{
		newArray:  newJSONArray,
		newObject: newJSONObject,
	}
	return ju.unmarshalJSONArray(data, ts)
}
