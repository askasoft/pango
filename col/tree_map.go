package col

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pandafw/pango/cmp"
	"github.com/pandafw/pango/iox"
	"github.com/pandafw/pango/str"
)

// NewTreeMap creates a new TreeMap.
// Example: NewTreeMap(cmp.CompareString, []P{{"k1", "v1"}, {"k2", "v2"}}...)
func NewTreeMap(compare cmp.Compare, kvs ...P) *TreeMap {
	tm := &TreeMap{compare: compare}
	tm.SetPairs(kvs...)
	return tm
}

// TreeMap implements an tree map that keeps the compare order of keys.
// The zero value for TreeMap is an empty map ready to use.
//
// https://en.wikipedia.org/wiki/Red%E2%80%93black_tree
//
// To iterate over a tree map (where tm is a *TreeMap):
//	for it := tm.Iterator(); it.Next(); {
//		// do something with it.Key(), it.Value()
//	}
//
type TreeMap struct {
	len     int
	root    *TreeMapNode
	compare cmp.Compare
}

//-----------------------------------------------------------
// implements Collection interface

// Len returns the length of the tree map.
func (tm *TreeMap) Len() int {
	return tm.len
}

// IsEmpty returns true if the map has no items
func (tm *TreeMap) IsEmpty() bool {
	return tm.len == 0
}

// Clear clears the map
func (tm *TreeMap) Clear() {
	tm.len = 0
	tm.root = nil
}

//-----------------------------------------------------------
// implements Map interface

// Keys returns the key slice
func (tm *TreeMap) Keys() []K {
	ks := make([]K, tm.len)
	for i, n := 0, tm.head(); n != nil; i, n = i+1, n.next() {
		ks[i] = n.key
	}
	return ks
}

// Values returns the value slice
func (tm *TreeMap) Values() []V {
	vs := make([]V, tm.len)
	for i, n := 0, tm.head(); n != nil; i, n = i+1, n.next() {
		vs[i] = n.value
	}
	return vs
}

// Contains looks for the given key, and returns true if the key exists in the map.
func (tm *TreeMap) Contains(ks ...K) bool {
	if len(ks) == 0 {
		return true
	}

	for _, k := range ks {
		if _, ok := tm.Get(k); !ok {
			return false
		}
	}
	return true
}

// Get looks for the given key, and returns the value associated with it,
// or nil if not found. The boolean it returns says whether the key is ok in the map.
func (tm *TreeMap) Get(key K) (V, bool) {
	node := tm.lookup(key)
	if node != nil {
		return node.value, true
	}
	return nil, false
}

// Set sets the paired key-value item, and returns what `Get` would have returned
// on that key prior to the call to `Set`.
// key should adhere to the comparator's type assertion, otherwise method panics.
func (tm *TreeMap) Set(key K, value V) (ov V, ok bool) {
	tn := tm.root
	if tn == nil {
		// Assert key is of comparator's type for initial tree
		tm.compare(key, key)

		tm.root = &TreeMapNode{key: key, value: value, color: black}
		tm.len = 1
		return
	}

	cmp := 0
	parent := tn
	for tn != nil {
		parent = tn
		cmp = tm.compare(key, tn.key)
		switch {
		case cmp < 0:
			tn = tn.left
		case cmp > 0:
			tn = tn.right
		default:
			ov, ok = tn.value, true
			tn.value = value
			return
		}
	}

	tn = &TreeMapNode{key: key, value: value, parent: parent}
	if cmp < 0 {
		parent.left = tn
	} else {
		parent.right = tn
	}

	tm.fixAfterInsertion(tn)

	tm.len++
	return
}

// SetPairs set items from key-value items array, override the existing items
func (tm *TreeMap) SetPairs(pairs ...P) {
	setMapPairs(tm, pairs...)
}

// SetAll add items from another map am, override the existing items
func (tm *TreeMap) SetAll(am Map) {
	setMapAll(tm, am)
}

// SetIfAbsent sets the key-value item if the key does not exists in the map,
// and returns true if the tree is changed.
func (tm *TreeMap) SetIfAbsent(key K, value V) (ov V, ok bool) {
	if node := tm.lookup(key); node != nil {
		return node.value, true
	}

	return tm.Set(key, value)
}

// Delete delete all items with key of ks,
// and returns what `Get` would have returned
// on that key prior to the call to `Delete`.
func (tm *TreeMap) Delete(ks ...K) (ov V, ok bool) {
	if tm.IsEmpty() {
		return
	}

	for _, k := range ks {
		ov, ok = tm.delete(k)
	}
	return
}

// Each call f for each item in the map
func (tm *TreeMap) Each(f func(k K, v V)) {
	for tn := tm.head(); tn != nil; tn = tn.next() {
		f(tn.key, tn.value)
	}
}

// ReverseEach call f for each item in the map with reverse order
func (tm *TreeMap) ReverseEach(f func(k K, v V)) {
	for tn := tm.tail(); tn != nil; tn = tn.prev() {
		f(tn.key, tn.value)
	}
}

// Iterator returns a iterator for the map
func (tm *TreeMap) Iterator() Iterator2 {
	return &treeMapIterator{tree: tm}
}

// Head returns a pointer to the minimum item.
func (tm *TreeMap) Head() *TreeMapNode {
	return tm.head()
}

// Tail returns a pointer to the maximum item.
func (tm *TreeMap) Tail() *TreeMapNode {
	return tm.tail()
}

// PollHead remove the first item of map.
func (tm *TreeMap) PollHead() *TreeMapNode {
	tn := tm.head()
	if tn != nil {
		tn = tm.deleteNode(tn)
	}
	return tn
}

// PollTail remove the last item of map.
func (tm *TreeMap) PollTail() *TreeMapNode {
	tn := tm.tail()
	if tn != nil {
		tn = tm.deleteNode(tn)
	}
	return tn
}

// Floor Finds floor node of the input key, return the floor node or nil if no floor is found.
//
// Floor node is defined as the largest node that is smaller than or equal to the given node.
// A floor node may not be found, either because the tree is empty, or because
// all nodes in the tree are larger than the given node.
//
// key should adhere to the comparator's type assertion, otherwise method panics.
func (tm *TreeMap) Floor(key K) *TreeMapNode {
	return tm.floor(key)
}

// Ceiling finds ceiling node of the input key, return the ceiling node or nil if no ceiling is found.
//
// Ceiling node is defined as the smallest node that is larger than or equal to the given node.
// A ceiling node may not be found, either because the tree is empty, or because
// all nodes in the tree are smaller than the given node.
//
// key should adhere to the comparator's type assertion, otherwise method panics.
func (tm *TreeMap) Ceiling(key K) *TreeMapNode {
	return tm.ceiling(key)
}

// Items returns the map item slice
func (tm *TreeMap) Items() []*TreeMapNode {
	ns := make([]*TreeMapNode, tm.Len())
	for i, n := 0, tm.Head(); n != nil; i, n = i+1, n.next() {
		ns[i] = n
	}
	return ns
}

// String print map to string
func (tm *TreeMap) String() string {
	bs, _ := json.Marshal(tm)
	return string(bs)
}

// Graph return the map's graph
func (tm *TreeMap) Graph(value bool) string {
	flag := 0
	if value {
		flag |= tmValue
	}
	return tm.root.graph(flag)
}

//-----------------------------------------------------------

// head returns a pointer to the minimum item.
func (tm *TreeMap) head() *TreeMapNode {
	tn := tm.root
	if tn != nil {
		for tn.left != nil {
			tn = tn.left
		}
	}
	return tn
}

// tail returns a pointer to the maximum item.
func (tm *TreeMap) tail() *TreeMapNode {
	tn := tm.root
	if tn != nil {
		for tn.right != nil {
			tn = tn.right
		}
	}
	return tn
}

// floor Finds floor node of the input key, return the floor node or nil if no floor is found.
func (tm *TreeMap) floor(key K) (floor *TreeMapNode) {
	node := tm.root
	for node != nil {
		compare := tm.compare(key, node.key)
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

// Ceiling finds ceiling node of the input key, return the ceiling node or nil if no ceiling is found.
func (tm *TreeMap) ceiling(key K) (ceiling *TreeMapNode) {
	node := tm.root
	for node != nil {
		compare := tm.compare(key, node.key)
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
// or nil if not found. The Node struct can then be used to iterate over the tree map
// from that point, either forward or backward.
func (tm *TreeMap) lookup(key K) *TreeMapNode {
	node := tm.root
	for node != nil {
		compare := tm.compare(key, node.key)
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

// delete delete the node from the tree by key,
// and returns what `Get` would have returned
// on that key prior to the call to `Delete`.
// key should adhere to the comparator's type assertion, otherwise method panics.
func (tm *TreeMap) delete(key K) (ov V, ok bool) {
	tn := tm.lookup(key)
	if tn == nil {
		return
	}

	ov, ok = tn.value, true
	tm.deleteNode(tn)
	return
}

// deleteNode delete the node p, returns the deleted node
// NOTE: if p has both left/right, p.next() will be deleted and returned
func (tm *TreeMap) deleteNode(p *TreeMapNode) *TreeMapNode {
	tm.len--

	// If strictly internal, copy successor's element to p and then make p point to successor.
	if p.left != nil && p.right != nil {
		s := p.next()

		p.key, s.key = s.key, p.key
		p.value, s.value = s.value, p.value
		p = s
	} // p has 2 children

	// Start fixup at replacement node, if it exists.
	replacement := p.right
	if p.left != nil {
		replacement = p.left
	}

	if replacement != nil {
		// Link replacement to parent
		replacement.parent = p.parent
		if p.parent == nil {
			tm.root = replacement
		} else if p == p.parent.left {
			p.parent.left = replacement
		} else {
			p.parent.right = replacement
		}

		// Null out links so they are OK to use by fixAfterDeletion.
		p.left, p.right, p.parent = nil, nil, nil

		// Fix replacement
		if p.color == black {
			tm.fixAfterDeletion(replacement)
		}
	} else if p.parent == nil { // return if we are the only node.
		tm.root = nil
	} else { //  No children. Use self as phantom replacement and unlink.
		if p.color == black {
			tm.fixAfterDeletion(p)
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

func (tm *TreeMap) fixAfterInsertion(x *TreeMapNode) {
	x.color = red

	for x != nil && x != tm.root && x.parent.color == red {
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
					tm.rotateLeft(x)
				}
				x.getParent().setColor(black)
				x.getGrandParent().setColor(red)
				tm.rotateRight(x.getGrandParent())
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
					tm.rotateRight(x)
				}
				x.getParent().setColor(black)
				x.getGrandParent().setColor(red)
				tm.rotateLeft(x.getGrandParent())
			}
		}
	}
	tm.root.color = black
}

func (tm *TreeMap) fixAfterDeletion(x *TreeMapNode) {
	for x != tm.root && x.getColor() == black {
		if x == x.getParent().getLeft() {
			sib := x.getParent().getRight()

			if sib.getColor() == red {
				sib.setColor(black)
				x.getParent().setColor(red)
				tm.rotateLeft(x.getParent())
				sib = x.getParent().getRight()
			}

			if sib.getLeft().getColor() == black && sib.getRight().getColor() == black {
				sib.setColor(red)
				x = x.getParent()
			} else {
				if sib.getRight().getColor() == black {
					sib.getLeft().setColor(black)
					sib.setColor(red)
					tm.rotateRight(sib)
					sib = x.getParent().getRight()
				}
				sib.setColor(x.getParent().getColor())
				x.getParent().setColor(black)
				sib.getRight().setColor(black)
				tm.rotateLeft(x.getParent())
				x = tm.root
			}
		} else { // symmetric
			sib := x.getParent().getLeft()

			if sib.getColor() == red {
				sib.setColor(black)
				x.getParent().setColor(red)
				tm.rotateRight(x.getParent())
				sib = x.getParent().getLeft()
			}

			if sib.getRight().getColor() == black && sib.getLeft().getColor() == black {
				sib.setColor(red)
				x = x.getParent()
			} else {
				if sib.getLeft().getColor() == black {
					sib.getRight().setColor(black)
					sib.setColor(red)
					tm.rotateLeft(sib)
					sib = x.getParent().getLeft()
				}
				sib.setColor(x.getParent().getColor())
				x.getParent().setColor(black)
				sib.getLeft().setColor(black)
				tm.rotateRight(x.getParent())
				x = tm.root
			}
		}
	}

	x.setColor(black)
}

func (tm *TreeMap) rotateLeft(p *TreeMapNode) {
	if p != nil {
		r := p.right
		p.right = r.left
		if r.left != nil {
			r.left.parent = p
		}
		r.parent = p.parent
		if p.parent == nil {
			tm.root = r
		} else if p.parent.left == p {
			p.parent.left = r
		} else {
			p.parent.right = r
		}
		r.left = p
		p.parent = r
	}
}

func (tm *TreeMap) rotateRight(p *TreeMapNode) {
	if p != nil {
		l := p.left
		p.left = l.right
		if l.right != nil {
			l.right.parent = p
		}
		l.parent = p.parent
		if p.parent == nil {
			tm.root = l
		} else if p.parent.right == p {
			p.parent.right = l
		} else {
			p.parent.left = l
		}
		l.right = p
		p.parent = l
	}
}

// debug return the map's graph (debug)
func (tm *TreeMap) debug() string {
	return tm.root.graph(tmColor | tmValue | tmPoint)
}

//-----------------------------------------------------

// color node color
type color byte

// Red Black color
const (
	red   color = 'R'
	black color = 'B'
)

func (c color) String() string {
	return fmt.Sprint(string(c))
}

// TreeMapNode is a node of red-black tree
type TreeMapNode struct {
	color  color
	left   *TreeMapNode
	right  *TreeMapNode
	parent *TreeMapNode
	key    K
	value  V
}

// Key returns the key
func (tn *TreeMapNode) Key() K {
	return tn.key
}

// Value returns the key
func (tn *TreeMapNode) Value() V {
	return tn.value
}

// SetValue sets the value
func (tn *TreeMapNode) SetValue(v V) {
	tn.value = v
}

func (tn *TreeMapNode) getLeft() *TreeMapNode {
	if tn != nil {
		return tn.left
	}
	return nil
}

func (tn *TreeMapNode) getRight() *TreeMapNode {
	if tn != nil {
		return tn.right
	}
	return nil
}

func (tn *TreeMapNode) getParent() *TreeMapNode {
	if tn != nil {
		return tn.parent
	}
	return nil
}

func (tn *TreeMapNode) getGrandParent() *TreeMapNode {
	if tn != nil && tn.parent != nil {
		return tn.parent.parent
	}
	return nil
}

func (tn *TreeMapNode) getColor() color {
	if tn == nil {
		return black
	}
	return tn.color
}

func (tn *TreeMapNode) setColor(c color) {
	if tn != nil {
		tn.color = c
	}
}

// prev returns the previous node or nil.
func (tn *TreeMapNode) prev() *TreeMapNode {
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
func (tn *TreeMapNode) next() *TreeMapNode {
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

// String print the key/value node to string
func (tn *TreeMapNode) String() string {
	return fmt.Sprintf("%v => %v", tn.key, tn.value)
}

const (
	tmColor = 1 << iota
	tmPoint
	tmValue
)

func (tn *TreeMapNode) graph(flag int) string {
	if tn == nil {
		return "(empty)"
	}

	sb := &strings.Builder{}
	tn.output(sb, "", true, flag)
	return sb.String()
}

func (tn *TreeMapNode) output(sb *strings.Builder, prefix string, tail bool, flag int) {
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

	if flag&tmColor == tmColor {
		sb.WriteString(fmt.Sprintf("(%v) ", tn.color))
	}
	sb.WriteString(fmt.Sprint(tn.key))
	if flag&tmPoint == tmPoint {
		sb.WriteString(fmt.Sprintf(" (%p)", tn))
	}
	if flag&tmValue == tmValue {
		v := str.RemoveAny(fmt.Sprint(tn.value), "\r\n")
		sb.WriteString(" => ")
		sb.WriteString(v)
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

// treeMapIterator a iterator for TreeMap
type treeMapIterator struct {
	tree    *TreeMap
	node    *TreeMapNode
	removed bool
}

// Prev moves the iterator to the previous element and returns true if there was a previous element in the container.
// If Prev() returns true, then previous element's index and value can be retrieved by Index() and Value().
// Modifies the state of the iterator.
func (it *treeMapIterator) Prev() bool {
	if it.tree.IsEmpty() {
		return false
	}

	if it.node == nil {
		it.node = it.tree.tail()
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
func (it *treeMapIterator) Next() bool {
	if it.tree.IsEmpty() {
		return false
	}

	if it.node == nil {
		it.node = it.tree.head()
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

// Key returns the current element's key.
func (it *treeMapIterator) Key() K {
	if it.node == nil {
		return nil
	}
	return it.node.key
}

// Value returns the current element's value.
func (it *treeMapIterator) Value() V {
	if it.node == nil {
		return nil
	}
	return it.node.value
}

// SetValue set the value to the item
func (it *treeMapIterator) SetValue(v V) {
	if it.node != nil {
		it.node.value = v
	}
}

// Remove remove the current element
func (it *treeMapIterator) Remove() {
	if it.node == nil {
		return
	}

	if it.removed {
		panic("TreeMap can't remove a unlinked item")
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
func (it *treeMapIterator) Reset() {
	it.node = nil
	it.removed = false
}

//-----------------------------------------------------------
// implements JSON Marshaller/Unmarshaller interface

func (tm *TreeMap) addJSONObjectItem(k string, v V) jsonObject {
	tm.Set(k, v)
	return tm
}

// MarshalJSON implements type json.Marshaler interface, so can be called in json.Marshal(tm)
func (tm *TreeMap) MarshalJSON() (res []byte, err error) {
	return jsonMarshalMap(tm)
}

// UnmarshalJSON implements type json.Unmarshaler interface, so can be called in json.Unmarshal(data, tm)
func (tm *TreeMap) UnmarshalJSON(data []byte) error {
	tm.Clear()
	ju := &jsonUnmarshaler{
		newArray:  newJSONArray,
		newObject: newJSONObject,
	}
	return ju.unmarshalJSONObject(data, tm)
}
