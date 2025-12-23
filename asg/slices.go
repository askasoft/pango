package asg

import (
	"cmp"
	"fmt"
	"slices"
	"strings"
)

// SliceOf returns []T{v[0], v[1], ...}
func SliceOf[T any](v ...T) []T {
	return v
}

// Anys convert slice 's' to []any slice
func Anys[T any](s []T) []any {
	b := make([]any, len(s))
	for i, v := range s {
		b[i] = v
	}
	return b
}

// Clone returns a copy of the slice with additional +n cap.
// The elements are copied using assignment, so this is a shallow clone.
func Clone[T any](s []T, n ...int) []T {
	// Preserve nil in case it matters.
	if s == nil {
		return nil
	}

	return append(make([]T, 0, len(s)+First(n)), s...)
}

// Contains reports whether the v is contained in the slice s.
func Contains[T comparable](s []T, v T) bool {
	return slices.Contains(s, v)
}

// ContainsAll reports whether all elements of cs are contained in the slice s.
func ContainsAll[T comparable](s []T, cs ...T) bool {
	if len(cs) == 0 {
		return true
	}

	if len(s) == 0 {
		return false
	}

	for _, c := range cs {
		if !Contains(s, c) {
			return false
		}
	}
	return true
}

// ContainsAny reports whether at least one element of cs is contained in the slice s.
func ContainsAny[T comparable](s []T, cs ...T) bool {
	if len(cs) == 0 {
		return true
	}

	if len(s) == 0 {
		return false
	}

	for _, c := range cs {
		if Contains(s, c) {
			return true
		}
	}
	return false
}

// ContainsFunc reports whether at least one element e of slice s which satisfies f(e).
func ContainsFunc[T any](s []T, f func(T) bool) bool {
	return slices.ContainsFunc(s, f)
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

// First get first element of a slice s.
// returns first value of `d...` if slice `s` is empty.
// returns zero value if slice `s` and `d` is empty.
func First[T any](s []T, d ...T) (v T) {
	switch {
	case len(s) > 0:
		v = s[0]
	case len(d) > 0:
		v = d[0]
	}
	return
}

// Last get last element of slice s.
// returns zero value if slice is empty.
func Last[T any](s []T) (v T) {
	if z := len(s); z > 0 {
		v = s[z-1]
	}
	return
}

// Get get element at the specified index i.
func Get[T any](s []T, i int) (v T, ok bool) {
	if i >= 0 && i < len(s) {
		v, ok = s[i], true
	}
	return
}

// Index returns the index of the first instance of v in s, or -1 if v is not present in s.
func Index[T comparable](s []T, v T) int {
	return slices.Index(s, v)
}

// IndexFunc returns the first index i satisfying f(s[i]), or -1 if none do.
func IndexFunc[T any](s []T, f func(T) bool) int {
	return slices.IndexFunc(s, f)
}

// FindFunc returns the first item satisfying f(s[i]), or (zero,false) if none do.
func FindFunc[T any](s []T, f func(T) bool) (v T, ok bool) {
	for _, e := range s {
		if f(e) {
			v, ok = e, true
			return
		}
	}
	return
}

// Delete removes the elements s[i:j] from s, returning the modified slice.
// Delete panics if s[i:j] is not a valid slice of s.
// Delete is O(len(s)-j), so if many items must be deleted, it is better to
// make a single call deleting them all together than to delete one at a time.
// Delete might not modify the elements s[len(s)-(j-i):len(s)]. If those
// elements contain pointers you might consider zeroing those elements so that
// objects they reference can be garbage collected.
func Delete[T any](s []T, i, j int) []T {
	return slices.Delete(s, i, j)
}

// DeleteFunc removes any elements from s for which del returns true,
// returning the modified slice.
// When DeleteFunc removes m elements, it might not modify the elements
// s[len(s)-m:len(s)]. If those elements contain pointers you might consider
// zeroing those elements so that objects they reference can be garbage
// collected.
func DeleteFunc[T any](s []T, del func(T) bool) []T {
	return slices.DeleteFunc(s, del)
}

// DeleteEqual removes any elements from s for which elemant == e, returning the modified slice.
// When DeleteFunc removes m elements, it might not modify the elements
// s[len(s)-m:len(s)]. If those elements contain pointers you might consider
// zeroing those elements so that objects they reference can be garbage
// collected.
func DeleteEqual[T comparable](s []T, e T) []T {
	i := Index(s, e)
	if i < 0 {
		return s
	}

	// Don't start copying elements until we find one to delete.
	for j := i + 1; j < len(s); j++ {
		if v := s[j]; v != e {
			s[i] = v
			i++
		}
	}
	clear(s[i:]) // zero/nil out the obsolete elements, for GC
	return s[:i]
}

// Insert inserts the values v... into s at index i,
// returning the modified slice.
// The elements at s[i:] are shifted up to make room.
// In the returned slice r, r[i] == v[0],
// and, if i < len(s), r[i+len(v)] == value originally at r[i].
// Insert panics if i > len(s).
// This function is O(len(s) + len(v)).
// If the result is empty, it has the same nilness as s.
func Insert[T any](s []T, i int, v ...T) []T {
	return slices.Insert(s, i, v...)
}

// Replace replaces the elements s[i:j] by the given v, and returns the
// modified slice.
// Replace panics if j > len(s) or s[i:j] is not a valid slice of s.
// When len(v) < (j-i), Replace zeroes the elements between the new length and the original length.
// If the result is empty, it has the same nilness as s.
func Replace[T any](s []T, i, j int, v ...T) []T {
	return slices.Replace(s, i, j, v...)
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
// Compact modifies the contents of the slice s and returns the modified slice,
// which may have s smaller length.
// When Compact discards m elements in total, it might not modify the elements
// s[len(s)-m:len(s)]. If those elements contain pointers you might consider
// zeroing those elements so that objects they reference can be garbage collected.
func Compact[T comparable](s []T) []T {
	return slices.Compact(s)
}

// CompactFunc is like [Compact] but uses an equality function to compare elements.
// For runs of elements that compare equal, CompactFunc keeps the first one.
func CompactFunc[T any](s []T, eq func(T, T) bool) []T {
	return slices.CompactFunc(s, eq)
}

// Grow increases the slice's capacity, if necessary, to guarantee space for
// another n elements. After Grow(n), at least n elements can be appended
// to the slice without another allocation. If n is negative or too large to
// allocate the memory, Grow panics.
func Grow[T any](s []T, n int) []T {
	return slices.Grow(s, n)
}

// Clip removes unused capacity from the slice, returning s[:len(s):len(s)].
func Clip[T any](s []T) []T {
	return s[:len(s):len(s)]
}

// Reverse reverses the elements of the slice in place.
func Reverse[T any](s []T) {
	slices.Reverse(s)
}

// Concat returns s new slice concatenating the passed in slices.
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

func sprint[T any](v T) string {
	return fmt.Sprint(v)
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
