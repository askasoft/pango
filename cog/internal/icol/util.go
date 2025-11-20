package icol

import "github.com/askasoft/pango/cog"

// ContainsAny Test to see if the collection contains any item of vs
func ContainsAny[T any](c cog.Collection[T], vs ...T) bool {
	if len(vs) == 0 {
		return true
	}

	if c.IsEmpty() {
		return false
	}

	for _, v := range vs {
		if c.Contains(v) {
			return true
		}
	}
	return false
}

// ContainsAll Test to see if the collection contains all items of vs
func ContainsAll[T any](c cog.Collection[T], vs ...T) bool {
	if len(vs) == 0 {
		return true
	}

	if c.IsEmpty() {
		return false
	}

	for _, v := range vs {
		if !c.Contains(v) {
			return false
		}
	}
	return true
}

// ContainsCol Test to see if the collection contains all items of another collection
func ContainsCol[T any](a cog.Collection[T], b cog.Collection[T]) bool {
	if b.IsEmpty() || a == b {
		return true
	}

	if a.IsEmpty() {
		return false
	}

	if i, ok := b.(cog.Iterable[T]); ok {
		return a.ContainsIter(i.Iterator())
	}

	return ContainsAll(a, b.Values()...)
}
