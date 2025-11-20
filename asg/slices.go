package asg

import (
	"cmp"
	"fmt"
	"slices"
	"strings"
)

// SliceOf returns a []T{args[0], args[1], ...}
func SliceOf[T any](args ...T) []T {
	return args
}

// Anys convert slice 'sa' to []any slice
func Anys[T any](a []T) []any {
	b := make([]any, len(a))
	for i, v := range a {
		b[i] = v
	}
	return b
}

// Clone returns a copy of the slice with additional +n cap.
// The elements are copied using assignment, so this is a shallow clone.
func Clone[T any](a []T, n ...int) []T {
	// Preserve nil in case it matters.
	if a == nil {
		return nil
	}

	return append(make([]T, 0, len(a)+First(n)), a...)
}

// Contains reports whether the v is contained in the slice a.
func Contains[T comparable](a []T, v T) bool {
	return slices.Contains(a, v)
}

// ContainsAll reports whether all elements of cs are contained in the slice a.
func ContainsAll[T comparable](a []T, cs ...T) bool {
	if len(cs) == 0 {
		return true
	}

	if len(a) == 0 {
		return false
	}

	for _, c := range cs {
		if !Contains(a, c) {
			return false
		}
	}
	return true
}

// ContainsAny reports whether at least one element of cs is contained in the slice a.
func ContainsAny[T comparable](a []T, cs ...T) bool {
	if len(cs) == 0 {
		return true
	}

	if len(a) == 0 {
		return false
	}

	for _, c := range cs {
		if Contains(a, c) {
			return true
		}
	}
	return false
}

// ContainsFunc reports whether at least one element e of a satisfies f(e).
func ContainsFunc[T any](a []T, f func(T) bool) bool {
	return slices.ContainsFunc(a, f)
}

// Equal reports whether a and b
// are the same length and contain the same element.
// A nil argument is equivalent to an empty slice.
func Equal[T comparable](a []T, b []T) bool {
	return slices.Equal(a, b)
}

// EqualFunc reports whether two slices are equal using an equality
// function on each pair of elements. If the lengths are different,
// EqualFunc returns false. Otherwise, the elements are compared in
// increasing index order, and the comparison stops at the first index
// for which eq returns false.
func EqualFunc[A any, B any](a []A, b []B, eq func(A, B) bool) bool {
	return slices.EqualFunc(a, b, eq)
}

// First get first element of a slice a.
// returns first value of `d...` if slice `a` is empty.
// returns zero value if slice `a` and `d` is empty.
func First[T any](a []T, d ...T) (v T) {
	switch {
	case len(a) > 0:
		v = a[0]
	case len(d) > 0:
		v = d[0]
	}
	return
}

// Last get last element of a slice.
// returns zero value if slice is empty.
func Last[T any](a []T) (v T) {
	if z := len(a); z > 0 {
		v = a[z-1]
	}
	return
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
	return slices.Index(a, v)
}

// IndexFunc returns the first index i satisfying f(a[i]), or -1 if none do.
func IndexFunc[T any](a []T, f func(T) bool) int {
	return slices.IndexFunc(a, f)
}

// FindFunc returns the first item satisfying f(a[i]), or (zero,false) if none do.
func FindFunc[T any](a []T, f func(T) bool) (v T, ok bool) {
	for _, e := range a {
		if f(e) {
			v, ok = e, true
			return
		}
	}
	return
}

// Delete removes the elements a[i:j] from a, returning the modified slice.
// Delete panics if a[i:j] is not a valid slice of a.
// Delete is O(len(a)-j), so if many items must be deleted, it is better to
// make a single call deleting them all together than to delete one at a time.
// Delete might not modify the elements a[len(a)-(j-i):len(a)]. If those
// elements contain pointers you might consider zeroing those elements so that
// objects they reference can be garbage collected.
func Delete[T any](a []T, i, j int) []T {
	return slices.Delete(a, i, j)
}

// DeleteFunc removes any elements from a for which del returns true,
// returning the modified slice.
// When DeleteFunc removes m elements, it might not modify the elements
// a[len(a)-m:len(a)]. If those elements contain pointers you might consider
// zeroing those elements so that objects they reference can be garbage
// collected.
func DeleteFunc[T any](a []T, del func(T) bool) []T {
	return slices.DeleteFunc(a, del)
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

// Compare compares the elements of s1 and s2, using [cmp.Compare] on each pair
// of elements. The elements are compared sequentially, starting at index 0,
// until one element is not equal to the other.
// The result of comparing the first non-matching elements is returned.
// If both slices are equal until one of them ends, the shorter slice is
// considered less than the longer one.
// The result is 0 if s1 == s2, -1 if s1 < s2, and +1 if s1 > s2.
func Compare[T cmp.Ordered](s1, s2 []T) int {
	return slices.Compare(s1, s2)
}

// CompareFunc is like [Compare] but uses a custom comparison function on each
// pair of elements.
// The result is the first non-zero result of cmp; if cmp always
// returns 0 the result is 0 if len(s1) == len(s2), -1 if len(s1) < len(s2),
// and +1 if len(s1) > len(s2).
func CompareFunc[S1 ~[]E1, S2 ~[]E2, E1, E2 any](s1 S1, s2 S2, cmp func(E1, E2) int) int {
	return slices.CompareFunc(s1, s2, cmp)
}

// Compact replaces consecutive runs of equal elements with a single copy.
// This is like the uniq command found on Unix.
// Compact modifies the contents of the slice a and returns the modified slice,
// which may have a smaller length.
// When Compact discards m elements in total, it might not modify the elements
// a[len(a)-m:len(a)]. If those elements contain pointers you might consider
// zeroing those elements so that objects they reference can be garbage collected.
func Compact[T comparable](a []T) []T {
	return slices.Compact(a)
}

// CompactFunc is like [Compact] but uses an equality function to compare elements.
// For runs of elements that compare equal, CompactFunc keeps the first one.
func CompactFunc[T any](a []T, eq func(T, T) bool) []T {
	return slices.CompactFunc(a, eq)
}

// Grow increases the slice's capacity, if necessary, to guarantee space for
// another n elements. After Grow(n), at least n elements can be appended
// to the slice without another allocation. If n is negative or too large to
// allocate the memory, Grow panics.
func Grow[T any](a []T, n int) []T {
	return slices.Grow(a, n)
}

// Clip removes unused capacity from the slice, returning a[:len(a):len(a)].
func Clip[T any](a []T) []T {
	return a[:len(a):len(a)]
}

// Reverse reverses the elements of the slice in place.
func Reverse[T any](a []T) {
	slices.Reverse(a)
}

// Concat returns a new slice concatenating the passed in slices.
func Concat[T any](ss ...[]T) []T {
	return slices.Concat(ss...)
}

// Min returns the minimal value in x. It panics if x is empty.
// For floating-point numbers, Min propagates NaNs (any NaN value in x
// forces the output to be NaN).
func Min[T cmp.Ordered](x []T) T {
	return slices.Min(x)
}

// MinFunc returns the minimal value in x, using cmp to compare elements.
// It panics if x is empty. If there is more than one minimal element
// according to the cmp function, MinFunc returns the first one.
func MinFunc[T any](x []T, cmp func(a, b T) int) T {
	return slices.MinFunc(x, cmp)
}

// Max returns the maximal value in x. It panics if x is empty.
// For floating-point E, Max propagates NaNs (any NaN value in x
// forces the output to be NaN).
func Max[T cmp.Ordered](x []T) T {
	return slices.Max(x)
}

// MaxFunc returns the maximal value in x, using cmp to compare elements.
// It panics if x is empty. If there is more than one maximal element
// according to the cmp function, MaxFunc returns the first one.
func MaxFunc[T any](x []T, cmp func(a, b T) int) T {
	return slices.MaxFunc(x, cmp)
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
