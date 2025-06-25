package treeset

import (
	"encoding/json"
	"iter"

	"github.com/askasoft/pango/cog"
	"github.com/askasoft/pango/cog/internal/iarray"
	"github.com/askasoft/pango/cog/internal/jsoncol"
	"github.com/askasoft/pango/str"
)

// NewTreeSet creates a new TreeSet.
// Example: cog.NewTreeSet(cog.CompareString, "v1", "v2")
func NewTreeSet[T any](compare cog.Compare[T], vs ...T) *TreeSet[T] {
	ts := &TreeSet[T]{compare: compare}
	ts.AddAll(vs...)
	return ts
}

// TreeSet implements an tree set that keeps the compare order of keys.
// The zero value for TreeSet is an empty set ready to use.
//
// https://en.wikipedia.org/wiki/Red%E2%80%93black_tree
//
// To iterate over a tree set (where ts is a *TreeSet):
//
//	for it := ts.Iterator(); it.Next(); {
//		// do something with it.Value()
//	}
type TreeSet[T any] struct {
	len     int
	root    *treeSetNode[T]
	compare cog.Compare[T]
}

//-----------------------------------------------------------
// implements Collection interface

// Len returns the length of the tree set.
func (ts *TreeSet[T]) Len() int {
	return ts.len
}

// IsEmpty returns true if the set has no items
func (ts *TreeSet[T]) IsEmpty() bool {
	return ts.len == 0
}

// Clear clears the set
func (ts *TreeSet[T]) Clear() {
	ts.len = 0
	ts.root = nil
}

// Add add item v.
func (ts *TreeSet[T]) Add(v T) {
	ts.add(v)
}

// AddAll adds all items of vs.
func (ts *TreeSet[T]) AddAll(vs ...T) {
	for _, v := range vs {
		ts.add(v)
	}
}

// AddCol adds all items of another collection
func (ts *TreeSet[T]) AddCol(ac cog.Collection[T]) {
	if ac.IsEmpty() || ts == ac {
		return
	}

	if ic, ok := ac.(cog.Iterable[T]); ok {
		it := ic.Iterator()
		for it.Next() {
			ts.add(it.Value())
		}
		return
	}

	ts.AddAll(ac.Values()...)
}

// Remove remove all items with associated value v
func (ts *TreeSet[T]) Remove(v T) {
	if tn := ts.lookup(v); tn != nil {
		ts.deleteNode(tn)
	}
}

// RemoveAll remove all items in the array vs
func (ts *TreeSet[T]) RemoveAll(vs ...T) {
	if ts.IsEmpty() {
		return
	}

	for _, v := range vs {
		ts.Remove(v)
	}
}

// RemoveCol remove all of this collection's elements that are also contained in the specified collection
func (ts *TreeSet[T]) RemoveCol(ac cog.Collection[T]) {
	if ts.IsEmpty() || ac.IsEmpty() {
		return
	}

	if ts == ac {
		ts.Clear()
		return
	}

	if ic, ok := ac.(cog.Iterable[T]); ok {
		ts.RemoveIter(ic.Iterator())
		return
	}

	ts.RemoveAll(ac.Values()...)
}

// RemoveIter remove all items in the iterator it
func (ts *TreeSet[T]) RemoveIter(it cog.Iterator[T]) {
	for it.Next() {
		ts.Remove(it.Value())
	}
}

// RemoveFunc remove all items that function f returns true
func (ts *TreeSet[T]) RemoveFunc(f func(T) bool) {
	if ts.IsEmpty() {
		return
	}

	for tn := ts.head(); tn != nil; tn = tn.next() {
		if f(tn.value) {
			ts.deleteNode(tn)
		}
	}
}

// Contains Test to see if the list contains the value v
func (ts *TreeSet[T]) Contains(v T) bool {
	return ts.lookup(v) != nil
}

// ContainsAll Test to see if the collection contains all items of vs
func (ts *TreeSet[T]) ContainsAll(vs ...T) bool {
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

// ContainsCol Test to see if the collection contains all items of another collection
func (ts *TreeSet[T]) ContainsCol(ac cog.Collection[T]) bool {
	if ac.IsEmpty() || ts == ac {
		return true
	}

	if ts.IsEmpty() {
		return false
	}

	if ic, ok := ac.(cog.Iterable[T]); ok {
		return ts.ContainsIter(ic.Iterator())
	}

	return ts.ContainsAll(ac.Values()...)
}

// ContainsIter Test to see if the collection contains all items of iterator 'it'
func (ts *TreeSet[T]) ContainsIter(it cog.Iterator[T]) bool {
	for it.Next() {
		if tn := ts.lookup(it.Value()); tn == nil {
			return false
		}
	}
	return true
}

// RetainAll Retains only the elements in this collection that are contained in the argument array vs.
func (ts *TreeSet[T]) RetainAll(vs ...T) {
	if ts.IsEmpty() {
		return
	}

	if len(vs) == 0 {
		ts.Clear()
		return
	}

	for tn := ts.head(); tn != nil; tn = tn.next() {
		if !iarray.Contains(vs, tn.value) {
			ts.deleteNode(tn)
		}
	}
}

// RetainCol Retains only the elements in this collection that are contained in the specified collection.
func (ts *TreeSet[T]) RetainCol(ac cog.Collection[T]) {
	if ts.IsEmpty() || ts == ac {
		return
	}

	if ac.IsEmpty() {
		ts.Clear()
		return
	}

	ts.RetainFunc(ac.Contains)
}

// RetainFunc Retains all items that function f returns true
func (ts *TreeSet[T]) RetainFunc(f func(T) bool) {
	if ts.IsEmpty() {
		return
	}

	for tn := ts.head(); tn != nil; tn = tn.next() {
		if !f(tn.value) {
			ts.deleteNode(tn)
		}
	}
}

// Values returns the value slice
func (ts *TreeSet[T]) Values() []T {
	vs := make([]T, ts.len)
	for i, n := 0, ts.head(); n != nil; i, n = i+1, n.next() {
		vs[i] = n.value
	}
	return vs
}

// Each call f for each item in the set
func (ts *TreeSet[T]) Each(f func(int, T) bool) {
	i := 0
	for tn := ts.head(); tn != nil; tn = tn.next() {
		if !f(i, tn.value) {
			return
		}
		i++
	}
}

// ReverseEach call f for each item in the set with reverse order
func (ts *TreeSet[T]) ReverseEach(f func(int, T) bool) {
	i := ts.len - 1
	for tn := ts.tail(); tn != nil; tn = tn.prev() {
		if !f(i, tn.value) {
			return
		}
		i--
	}
}

// Iterator returns a iterator for the set
func (ts *TreeSet[T]) Iterator() cog.Iterator[T] {
	return &treeSetIterator[T]{tree: ts}
}

// Seq returns a iter.Seq[T] for range
func (ts *TreeSet[T]) Seq() iter.Seq[T] {
	return func(yield func(T) bool) {
		for tn := ts.head(); tn != nil; tn = tn.next() {
			if !yield(tn.value) {
				return
			}
		}
	}
}

//----------------------------------------------------------------

// PeekHead get the first item of set.
func (ts *TreeSet[T]) PeekHead() (v T, ok bool) {
	tn := ts.head()
	if tn != nil {
		v, ok = tn.value, true
	}
	return
}

// PeekTail get the last item of set.
func (ts *TreeSet[T]) PeekTail() (v T, ok bool) {
	tn := ts.tail()
	if tn != nil {
		v, ok = tn.value, true
	}
	return
}

// PollHead remove the first item of set.
func (ts *TreeSet[T]) PollHead() (v T, ok bool) {
	tn := ts.head()
	if tn != nil {
		v, ok = tn.value, true
		ts.deleteNode(tn)
	}
	return
}

// PollTail remove the last item of set.
func (ts *TreeSet[T]) PollTail() (v T, ok bool) {
	tn := ts.tail()
	if tn != nil {
		v, ok = tn.value, true
		ts.deleteNode(tn)
	}
	return
}

//----------------------------------------------------------------

// Head returns the first item of set ts or nil if the set is empty.
func (ts *TreeSet[T]) Head() (v T) {
	v, _ = ts.PeekHead()
	return
}

// Tail returns the last item of set ts or nil if the set is empty.
func (ts *TreeSet[T]) Tail() (v T) {
	v, _ = ts.PeekTail()
	return
}

// Floor Finds floor node of the input key, return the floor node's value or nil if no floor is found.
//
// Floor node is defined as the largest node that is smaller than or equal to the given node.
// A floor node may not be found, either because the tree is empty, or because
// all nodes in the tree are larger than the given node.
//
// key should adhere to the comparator's type assertion, otherwise method panics.
func (ts *TreeSet[T]) Floor(v T) *T {
	tn := ts.floor(v)
	if tn != nil {
		return &tn.value
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
func (ts *TreeSet[T]) Ceiling(v T) *T {
	tn := ts.ceiling(v)
	if tn != nil {
		return &tn.value
	}
	return nil
}

// String print set to string
func (ts *TreeSet[T]) String() string {
	bs, _ := json.Marshal(ts)
	return str.UnsafeString(bs)
}

// Graph return the set's graph
func (ts *TreeSet[T]) Graph() string {
	return ts.root.graph(0)
}

// -----------------------------------------------------
func (ts *TreeSet[T]) setValue(tn *treeSetNode[T], v T) *treeSetNode[T] {
	if any(tn.value) == any(v) {
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

// head returns a pointer to the minimum node.
func (ts *TreeSet[T]) head() *treeSetNode[T] {
	tn := ts.root
	if tn != nil {
		for tn.left != nil {
			tn = tn.left
		}
	}
	return tn
}

// tail returns a pointer to the maximum node.
func (ts *TreeSet[T]) tail() *treeSetNode[T] {
	tn := ts.root
	if tn != nil {
		for tn.right != nil {
			tn = tn.right
		}
	}
	return tn
}

// floor Finds floor node of the input key, return the floor node or nil if no floor is found.
func (ts *TreeSet[T]) floor(key T) (floor *treeSetNode[T]) {
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
func (ts *TreeSet[T]) ceiling(key T) (ceiling *treeSetNode[T]) {
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
func (ts *TreeSet[T]) lookup(key T) *treeSetNode[T] {
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
func (ts *TreeSet[T]) add(v T) *treeSetNode[T] {
	tn := ts.root
	if tn == nil {
		// Assert key is of comparator's type for initial tree
		ts.compare(v, v)

		ts.root = &treeSetNode[T]{value: v, color: black}
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

	tn = &treeSetNode[T]{value: v, parent: parent}
	if cmp < 0 {
		parent.left = tn
	} else {
		parent.right = tn
	}

	ts.fixAfterInsertion(tn)

	ts.len++
	return tn
}

// deleteNode delete the node p, returns the deleted node
// NOTE: if p has both left/right, p.next() will be deleted and returned
func (ts *TreeSet[T]) deleteNode(p *treeSetNode[T]) *treeSetNode[T] {
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
			switch p {
			case p.parent.left:
				p.parent.left = nil
			case p.parent.right:
				p.parent.right = nil
			}
			p.parent = nil
		}
	}

	return p
}

func (ts *TreeSet[T]) fixAfterInsertion(x *treeSetNode[T]) {
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

func (ts *TreeSet[T]) fixAfterDeletion(x *treeSetNode[T]) {
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

func (ts *TreeSet[T]) rotateLeft(p *treeSetNode[T]) {
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

func (ts *TreeSet[T]) rotateRight(p *treeSetNode[T]) {
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
func (ts *TreeSet[T]) debug() string {
	return ts.root.graph(tsColor | tsPoint)
}

//-----------------------------------------------------------
// implements JSON Marshaller/Unmarshaller interface

// MarshalJSON implements type json.Marshaler interface, so can be called in json.Marshal(ts)
func (ts *TreeSet[T]) MarshalJSON() ([]byte, error) {
	return jsoncol.JsonMarshalCol(ts)
}

// UnmarshalJSON implements type json.Unmarshaler interface, so can be called in json.Unmarshal(data, ts)
func (ts *TreeSet[T]) UnmarshalJSON(data []byte) error {
	ts.Clear()
	return jsoncol.JsonUnmarshalCol(data, ts)
}
