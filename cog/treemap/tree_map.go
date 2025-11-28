package treemap

import (
	"fmt"
	"iter"

	"github.com/askasoft/pango/cog"
	"github.com/askasoft/pango/cog/internal/imap"
	"github.com/askasoft/pango/doc/jsonx"
)

// NewTreeMap creates a new TreeMap.
// Example: cog.NewTreeMap(cog.CompareString, cog.KV("k1", "v1"), cog.KV("k2", "v2"))
func NewTreeMap[K any, V any](compare cog.Compare[K], kvs ...cog.P[K, V]) *TreeMap[K, V] {
	tm := &TreeMap[K, V]{compare: compare}
	tm.SetEntries(kvs...)
	return tm
}

// TreeMap implements an tree map that keeps the compare order of keys.
// The zero value for TreeMap is an empty map ready to use.
//
// https://en.wikipedia.org/wiki/Red%E2%80%93black_tree
//
// To iterate over a tree map (where tm is a *TreeMap):
//
//	for it := tm.Iterator(); it.Next(); {
//		// do something with it.Key(), it.Value()
//	}
type TreeMap[K any, V any] struct {
	len     int
	root    *TreeMapNode[K, V]
	compare cog.Compare[K]
}

//-----------------------------------------------------------
// implements Collection interface

// Len returns the length of the tree map.
func (tm *TreeMap[K, V]) Len() int {
	return tm.len
}

// IsEmpty returns true if the map has no items
func (tm *TreeMap[K, V]) IsEmpty() bool {
	return tm.len == 0
}

// Clear clears the map
func (tm *TreeMap[K, V]) Clear() {
	tm.len = 0
	tm.root = nil
}

//-----------------------------------------------------------
// implements Map interface

// Keys returns the key slice
func (tm *TreeMap[K, V]) Keys() []K {
	ks := make([]K, tm.len)
	for i, n := 0, tm.head(); n != nil; i, n = i+1, n.next() {
		ks[i] = n.key
	}
	return ks
}

// Values returns the value slice
func (tm *TreeMap[K, V]) Values() []V {
	vs := make([]V, tm.len)
	for i, n := 0, tm.head(); n != nil; i, n = i+1, n.next() {
		vs[i] = n.value
	}
	return vs
}

// Entries returns the key-value pair slice
func (tm *TreeMap[K, V]) Entries() []cog.P[K, V] {
	ps := make([]cog.P[K, V], tm.Len())
	for i, n := 0, tm.head(); n != nil; i, n = i+1, n.next() {
		ps[i] = cog.KV(n.key, n.value)
	}
	return ps
}

// Contains Test to see if the map contains the given key
func (tm *TreeMap[K, V]) Contains(k K) bool {
	if tm.IsEmpty() {
		return false
	}
	if _, ok := tm.Get(k); ok {
		return true
	}
	return false
}

// ContainsAny Test to see if the map contains any key of ks
func (tm *TreeMap[K, V]) ContainsAny(ks ...K) bool {
	if len(ks) == 0 {
		return true
	}

	if tm.IsEmpty() {
		return false
	}

	for _, k := range ks {
		if _, ok := tm.Get(k); ok {
			return true
		}
	}
	return false
}

// ContainsAll Test to see if the map contains all keys of ks
func (tm *TreeMap[K, V]) ContainsAll(ks ...K) bool {
	if len(ks) == 0 {
		return true
	}

	if tm.IsEmpty() {
		return false
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
func (tm *TreeMap[K, V]) Get(key K) (V, bool) {
	node := tm.lookup(key)
	if node != nil {
		return node.value, true
	}

	var v V
	return v, false
}

// MustGet looks for the given key, and returns the value associated with it.
// Panic if not found.
func (tm *TreeMap[K, V]) MustGet(key K) V {
	if v, ok := tm.Get(key); ok {
		return v
	}
	panic(fmt.Errorf("TreeMap key '%v' does not exist", key))
}

// SafeGet looks for the given key, and returns the value associated with it.
// If not found, return defaults[0] or zero V.
func (tm *TreeMap[K, V]) SafeGet(key K, defaults ...V) V {
	if v, ok := tm.Get(key); ok {
		return v
	}
	if len(defaults) > 0 {
		return defaults[0]
	}

	var v V
	return v
}

// Set sets the paired key-value item, and returns what `Get` would have returned
// on that key prior to the call to `Set`.
// key should adhere to the comparator's type assertion, otherwise method panics.
func (tm *TreeMap[K, V]) Set(key K, value V) (ov V, ok bool) {
	tn := tm.root
	if tn == nil {
		// Assert key is of comparator's type for initial tree
		tm.compare(key, key)

		tm.root = &TreeMapNode[K, V]{key: key, value: value, color: black}
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

	tn = &TreeMapNode[K, V]{key: key, value: value, parent: parent}
	if cmp < 0 {
		parent.left = tn
	} else {
		parent.right = tn
	}

	tm.fixAfterInsertion(tn)

	tm.len++
	return
}

// SetIfAbsent sets the key-value item if the key does not exists in the map,
// and returns true if the tree is changed.
func (tm *TreeMap[K, V]) SetIfAbsent(key K, value V) (ov V, ok bool) {
	if node := tm.lookup(key); node != nil {
		return node.value, true
	}

	return tm.Set(key, value)
}

// SetEntries set items from key-value items array, override the existing items
func (tm *TreeMap[K, V]) SetEntries(pairs ...cog.P[K, V]) {
	imap.SetMapPairs(tm, pairs...)
}

// Copy copy items from another map am, override the existing items
func (tm *TreeMap[K, V]) Copy(am cog.Map[K, V]) {
	imap.CopyMap(tm, am)
}

// Remove remove the item with key k,
// and returns what `Get` would have returned
// on that key prior to the call to `Set`.
func (tm *TreeMap[K, V]) Remove(k K) (ov V, ok bool) {
	if tm.IsEmpty() {
		return
	}

	return tm.delete(k)
}

// RemoveAll remove all items with key of ks.
func (tm *TreeMap[K, V]) RemoveAll(ks ...K) {
	if tm.IsEmpty() {
		return
	}

	for _, k := range ks {
		tm.delete(k)
	}
}

// Each call f for each item in the map
func (tm *TreeMap[K, V]) Each(f func(K, V) bool) {
	for tn := tm.head(); tn != nil; tn = tn.next() {
		if !f(tn.key, tn.value) {
			return
		}
	}
}

// ReverseEach call f for each item in the map with reverse order
func (tm *TreeMap[K, V]) ReverseEach(f func(K, V) bool) {
	for tn := tm.tail(); tn != nil; tn = tn.prev() {
		if !f(tn.key, tn.value) {
			return
		}
	}
}

// Seq returns a iter.Seq[K, V] for range
func (tm *TreeMap[K, V]) Seq() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for tn := tm.head(); tn != nil; tn = tn.next() {
			if !yield(tn.key, tn.value) {
				return
			}
		}
	}
}

// Iterator returns a iterator for the map
func (tm *TreeMap[K, V]) Iterator() cog.Iterator2[K, V] {
	return &treeMapIterator[K, V]{tree: tm}
}

// Head returns a pointer to the minimum item.
func (tm *TreeMap[K, V]) Head() *TreeMapNode[K, V] {
	return tm.head()
}

// Tail returns a pointer to the maximum item.
func (tm *TreeMap[K, V]) Tail() *TreeMapNode[K, V] {
	return tm.tail()
}

// PollHead remove the first item of map.
func (tm *TreeMap[K, V]) PollHead() *TreeMapNode[K, V] {
	tn := tm.head()
	if tn != nil {
		tn = tm.deleteNode(tn)
	}
	return tn
}

// PollTail remove the last item of map.
func (tm *TreeMap[K, V]) PollTail() *TreeMapNode[K, V] {
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
func (tm *TreeMap[K, V]) Floor(key K) *TreeMapNode[K, V] {
	return tm.floor(key)
}

// Ceiling finds ceiling node of the input key, return the ceiling node or nil if no ceiling is found.
//
// Ceiling node is defined as the smallest node that is larger than or equal to the given node.
// A ceiling node may not be found, either because the tree is empty, or because
// all nodes in the tree are smaller than the given node.
//
// key should adhere to the comparator's type assertion, otherwise method panics.
func (tm *TreeMap[K, V]) Ceiling(key K) *TreeMapNode[K, V] {
	return tm.ceiling(key)
}

// Items returns the map item slice
func (tm *TreeMap[K, V]) Items() []*TreeMapNode[K, V] {
	ns := make([]*TreeMapNode[K, V], tm.Len())
	for i, n := 0, tm.Head(); n != nil; i, n = i+1, n.next() {
		ns[i] = n
	}
	return ns
}

// String print map to string
func (tm *TreeMap[K, V]) String() string {
	return jsonx.Stringify(tm)
}

// Graph return the map's graph
func (tm *TreeMap[K, V]) Graph(value bool) string {
	flag := 0
	if value {
		flag |= tmValue
	}
	return tm.root.graph(flag)
}

//-----------------------------------------------------------

// head returns a pointer to the minimum item.
func (tm *TreeMap[K, V]) head() *TreeMapNode[K, V] {
	tn := tm.root
	if tn != nil {
		for tn.left != nil {
			tn = tn.left
		}
	}
	return tn
}

// tail returns a pointer to the maximum item.
func (tm *TreeMap[K, V]) tail() *TreeMapNode[K, V] {
	tn := tm.root
	if tn != nil {
		for tn.right != nil {
			tn = tn.right
		}
	}
	return tn
}

// floor Finds floor node of the input key, return the floor node or nil if no floor is found.
func (tm *TreeMap[K, V]) floor(key K) (floor *TreeMapNode[K, V]) {
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
func (tm *TreeMap[K, V]) ceiling(key K) (ceiling *TreeMapNode[K, V]) {
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
func (tm *TreeMap[K, V]) lookup(key K) *TreeMapNode[K, V] {
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
func (tm *TreeMap[K, V]) delete(key K) (ov V, ok bool) {
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
func (tm *TreeMap[K, V]) deleteNode(p *TreeMapNode[K, V]) *TreeMapNode[K, V] {
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

func (tm *TreeMap[K, V]) fixAfterInsertion(x *TreeMapNode[K, V]) {
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

func (tm *TreeMap[K, V]) fixAfterDeletion(x *TreeMapNode[K, V]) {
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

func (tm *TreeMap[K, V]) rotateLeft(p *TreeMapNode[K, V]) {
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

func (tm *TreeMap[K, V]) rotateRight(p *TreeMapNode[K, V]) {
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
func (tm *TreeMap[K, V]) debug() string {
	return tm.root.graph(tmColor | tmValue | tmPoint)
}

//-----------------------------------------------------------
// implements JSON Marshaller/Unmarshaller interface

// MarshalJSON implements type json.Marshaler interface, so can be called in json.Marshal(tm)
func (tm *TreeMap[K, V]) MarshalJSON() ([]byte, error) {
	return imap.JsonMarshalMap(tm)
}

// UnmarshalJSON implements type json.Unmarshaler interface, so can be called in json.Unmarshal(data, tm)
func (tm *TreeMap[K, V]) UnmarshalJSON(data []byte) error {
	tm.Clear()
	return imap.JsonUnmarshalMap(data, tm)
}
