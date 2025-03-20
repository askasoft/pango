//go:build go1.21
// +build go1.21

package asg

import (
	"cmp"
	"fmt"
	"strings"
)

// Anys convert slice 'sa' to []any slice
func Anys[T any](sa []T) []any {
	sb := make([]any, len(sa))
	for i, a := range sa {
		sb[i] = a
	}
	return sb
}

// Clone returns a copy of the slice.
// The elements are copied using assignment, so this is a shallow clone.
func Clone[T any](a []T) []T {
	// Preserve nil in case it matters.
	if a == nil {
		return nil
	}
	return append([]T{}, a...)
}

// Contains reports whether the c is contained in the slice a.
func Contains[T comparable](a []T, c T) bool {
	return Index(a, c) >= 0
}

// ContainsAny reports whether the any item of cs is contained in the slice a.
func ContainsAny[T comparable](a []T, cs ...T) bool {
	for _, c := range cs {
		if Contains(a, c) {
			return true
		}
	}
	return false
}

// ContainsFunc reports whether at least one element e of a satisfies f(e).
func ContainsFunc[T any](a []T, f func(T) bool) bool {
	return IndexFunc(a, f) >= 0
}

// Equal reports whether a and b
// are the same length and contain the same element.
// A nil argument is equivalent to an empty slice.
func Equal[T comparable](a []T, b []T) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// EqualFunc reports whether two slices are equal using an equality
// function on each pair of elements. If the lengths are different,
// EqualFunc returns false. Otherwise, the elements are compared in
// increasing index order, and the comparison stops at the first index
// for which eq returns false.
func EqualFunc[A any, B any](a []A, b []B, eq func(A, B) bool) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if !eq(a[i], b[i]) {
			return false
		}
	}
	return true
}

// Get get element at the specified index i.
func Get[T any](a []T, i int) (v T, ok bool) {
	if i >= 0 && i < len(a) {
		v, ok = a[i], true
	}
	return
}

// Index returns the index of the first instance of v in a, or -1 if v is not present in a.
func Index[T comparable](a []T, v T) int {
	for i, e := range a {
		if e == v {
			return i
		}
	}
	return -1
}

// IndexFunc returns the first index i satisfying f(a[i]), or -1 if none do.
func IndexFunc[T any](a []T, f func(T) bool) int {
	for i, v := range a {
		if f(v) {
			return i
		}
	}
	return -1
}

// Delete removes the elements a[i:j] from a, returning the modified slice.
// Delete panics if a[i:j] is not a valid slice of a.
// Delete is O(len(a)-j), so if many items must be deleted, it is better to
// make a single call deleting them all together than to delete one at a time.
// Delete might not modify the elements a[len(a)-(j-i):len(a)]. If those
// elements contain pointers you might consider zeroing those elements so that
// objects they reference can be garbage collected.
func Delete[T any](a []T, i, j int) []T {
	_ = a[i:j] // bounds check
	copy(a[i:], a[j:])
	return a[:len(a)+i-j]
}

// DeleteFunc removes any elements from a for which del returns true,
// returning the modified slice.
// When DeleteFunc removes m elements, it might not modify the elements
// a[len(a)-m:len(a)]. If those elements contain pointers you might consider
// zeroing those elements so that objects they reference can be garbage
// collected.
func DeleteFunc[T any](a []T, del func(T) bool) []T {
	i := IndexFunc(a, del)
	if i < 0 {
		return a
	}

	// Don't start copying elements until we find one to delete.
	for j := i + 1; j < len(a); j++ {
		if v := a[j]; !del(v) {
			a[i] = v
			i++
		}
	}
	return a[:i]
}

// DeleteEqual removes any elements from a for which elemant == e, returning the modified slice.
// When DeleteFunc removes m elements, it might not modify the elements
// a[len(a)-m:len(a)]. If those elements contain pointers you might consider
// zeroing those elements so that objects they reference can be garbage
// collected.
func DeleteEqual[T comparable](a []T, e T) []T {
	i := Index(a, e)
	if i < 0 {
		return a
	}

	// Don't start copying elements until we find one to delete.
	for j := i + 1; j < len(a); j++ {
		if v := a[j]; v != e {
			a[i] = v
			i++
		}
	}
	return a[:i]
}

// Compact replaces consecutive runs of equal elements with a single copy.
// This is like the uniq command found on Unix.
// Compact modifies the contents of the slice a and returns the modified slice,
// which may have a smaller length.
// When Compact discards m elements in total, it might not modify the elements
// a[len(a)-m:len(a)]. If those elements contain pointers you might consider
// zeroing those elements so that objects they reference can be garbage collected.
func Compact[T comparable](a []T) []T {
	if len(a) < 2 {
		return a
	}
	i := 1
	for k := 1; k < len(a); k++ {
		if a[k] != a[k-1] {
			if i != k {
				a[i] = a[k]
			}
			i++
		}
	}
	return a[:i]
}

// CompactFunc is like [Compact] but uses an equality function to compare elements.
// For runs of elements that compare equal, CompactFunc keeps the first one.
func CompactFunc[T any](a []T, eq func(T, T) bool) []T {
	if len(a) < 2 {
		return a
	}
	i := 1
	for k := 1; k < len(a); k++ {
		if !eq(a[k], a[k-1]) {
			if i != k {
				a[i] = a[k]
			}
			i++
		}
	}
	return a[:i]
}

// Grow increases the slice's capacity, if necessary, to guarantee space for
// another n elements. After Grow(n), at least n elements can be appended
// to the slice without another allocation. If n is negative or too large to
// allocate the memory, Grow panics.
func Grow[T any](a []T, n int) []T {
	if n < 0 {
		panic("cannot be negative")
	}
	if n -= cap(a) - len(a); n > 0 {
		a = append(a[:cap(a)], make([]T, n)...)[:len(a)]
	}
	return a
}

// Clip removes unused capacity from the slice, returning a[:len(a):len(a)].
func Clip[T any](a []T) []T {
	return a[:len(a):len(a)]
}

// Reverse reverses the elements of the slice in place.
func Reverse[T any](a []T) {
	for i, j := 0, len(a)-1; i < j; i, j = i+1, j-1 {
		a[i], a[j] = a[j], a[i]
	}
}

// Min returns the minimal value in x. It panics if x is empty.
// For floating-point numbers, Min propagates NaNs (any NaN value in x
// forces the output to be NaN).
func Min[T cmp.Ordered](x []T) T {
	if len(x) < 1 {
		panic("asg.Min: empty list")
	}

	m := x[0]
	for i := 1; i < len(x); i++ {
		m = min(m, x[i])
	}
	return m
}

// MinFunc returns the minimal value in x, using cmp to compare elements.
// It panics if x is empty. If there is more than one minimal element
// according to the cmp function, MinFunc returns the first one.
func MinFunc[T any](x []T, cmp func(a, b T) int) T {
	if len(x) < 1 {
		panic("asg.MinFunc: empty list")
	}

	m := x[0]
	for i := 1; i < len(x); i++ {
		if cmp(x[i], m) < 0 {
			m = x[i]
		}
	}
	return m
}

// Max returns the maximal value in x. It panics if x is empty.
// For floating-point E, Max propagates NaNs (any NaN value in x
// forces the output to be NaN).
func Max[T cmp.Ordered](x []T) T {
	if len(x) < 1 {
		panic("asg.Max: empty list")
	}

	m := x[0]
	for i := 1; i < len(x); i++ {
		m = max(m, x[i])
	}
	return m
}

// MaxFunc returns the maximal value in x, using cmp to compare elements.
// It panics if x is empty. If there is more than one maximal element
// according to the cmp function, MaxFunc returns the first one.
func MaxFunc[T any](x []T, cmp func(a, b T) int) T {
	if len(x) < 1 {
		panic("asg.MaxFunc: empty list")
	}

	m := x[0]
	for i := 1; i < len(x); i++ {
		if cmp(x[i], m) > 0 {
			m = x[i]
		}
	}
	return m
}

func sprint[T any](a T) string {
	return fmt.Sprint(a)
}

// Join concatenates the elements of its first argument to create a single string. The separator
// string sep is placed between elements in the resulting string.
func Join[T any](elems []T, sep string, fmt ...func(T) string) string {
	sp := sprint[T]
	if len(fmt) > 0 {
		sp = fmt[0]
	}

	switch len(elems) {
	case 0:
		return ""
	case 1:
		return sp(elems[0])
	}

	var b strings.Builder
	b.WriteString(sp(elems[0]))
	for _, n := range elems[1:] {
		b.WriteString(sep)
		b.WriteString(sp(n))
	}
	return b.String()
}
