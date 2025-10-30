package bye

import (
	"bytes"
	"unicode"
)

// Clone returns a copy of b[:len(b)].
// The result may have additional unused capacity.
// Clone(nil) returns nil.
func Clone(b []byte) []byte {
	return bytes.Clone(b)
}

// Compare returns an integer comparing two byte slices lexicographically.
// The result will be 0 if a==b, -1 if a < b, and +1 if a > b.
// A nil argument is equivalent to an empty slice.
func Compare(a, b []byte) int {
	return bytes.Compare(a, b)
}

// Contains reports whether subslice is within b.
func Contains(b, subslice []byte) bool {
	return bytes.Contains(b, subslice)
}

// ContainsAny reports whether any of the UTF-8-encoded code points in chars are within b.
func ContainsAny(b []byte, chars string) bool {
	return bytes.ContainsAny(b, chars)
}

// ContainsFunc reports whether any of the UTF-8-encoded code points r within b satisfy f(r).
func ContainsFunc(b []byte, f func(rune) bool) bool {
	return bytes.ContainsFunc(b, f)
}

// ContainsRune reports whether the rune is contained in the UTF-8-encoded byte slice b.
func ContainsRune(b []byte, r rune) bool {
	return bytes.ContainsRune(b, r)
}

// Count counts the number of non-overlapping instances of sep in s.
// If sep is an empty slice, Count returns 1 + the number of UTF-8-encoded code points in s.
func Count(s, sep []byte) int {
	return bytes.Count(s, sep)
}

// Cut slices s around the first instance of sep,
// returning the text before and after sep.
// The found result reports whether sep appears in s.
// If sep does not appear in s, cut returns s, nil, false.
//
// Cut returns slices of the original slice s, not copies.
func Cut(s, sep []byte) (before, after []byte, found bool) {
	return bytes.Cut(s, sep)
}

// CutPrefix returns s without the provided leading prefix byte slice
// and reports whether it found the prefix.
// If s doesn't start with prefix, CutPrefix returns s, false.
// If prefix is the empty byte slice, CutPrefix returns s, true.
//
// CutPrefix returns slices of the original slice s, not copies.
func CutPrefix(s, prefix []byte) (after []byte, found bool) {
	return bytes.CutPrefix(s, prefix)
}

// CutSuffix returns s without the provided ending suffix byte slice
// and reports whether it found the suffix.
// If s doesn't end with suffix, CutSuffix returns s, false.
// If suffix is the empty byte slice, CutSuffix returns s, true.
//
// CutSuffix returns slices of the original slice s, not copies.
func CutSuffix(s, suffix []byte) (before []byte, found bool) {
	return bytes.CutSuffix(s, suffix)
}

// Equal reports whether a and b
// are the same length and contain the same bytes.
// A nil argument is equivalent to an empty slice.
func Equal(a, b []byte) bool {
	return bytes.Equal(a, b)
}

// EqualFold reports whether s and t, interpreted as UTF-8 strings,
// are equal under Unicode case-folding, which is a more general
// form of case-insensitivity.
func EqualFold(s, t []byte) bool {
	return bytes.EqualFold(s, t)
}

// Fields interprets s as a sequence of UTF-8-encoded code points.
// It splits the slice s around each instance of one or more consecutive white space
// characters, as defined by unicode.IsSpace, returning a slice of subslices of s or an
// empty slice if s contains only white space.
func Fields(s []byte) [][]byte {
	return bytes.Fields(s)
}

// FieldsFunc interprets s as a sequence of UTF-8-encoded code points.
// It splits the slice s at each run of code points c satisfying f(c) and
// returns a slice of subslices of s. If all code points in s satisfy f(c), or
// len(s) == 0, an empty slice is returned.
// FieldsFunc makes no guarantees about the order in which it calls f(c).
// If f does not return consistent results for a given c, FieldsFunc may crash.
func FieldsFunc(s []byte, f func(rune) bool) [][]byte {
	return bytes.FieldsFunc(s, f)
}

// HasPrefix tests whether the byte slice s begins with prefix.
func HasPrefix(s, prefix []byte) bool {
	return bytes.HasPrefix(s, prefix)
}

// HasSuffix tests whether the byte slice s ends with suffix.
func HasSuffix(s, suffix []byte) bool {
	return bytes.HasSuffix(s, suffix)
}

// Index returns the index of the first instance of sep in s, or -1 if sep is not present in s.
func Index(s, sep []byte) int {
	return bytes.Index(s, sep)
}

// IndexAny interprets s as a sequence of UTF-8-encoded Unicode code points.
// It returns the byte index of the first occurrence in s of any of the Unicode
// code points in chars. It returns -1 if chars is empty or if there is no code
// point in common.
func IndexAny(s []byte, chars string) int {
	return bytes.IndexAny(s, chars)
}

// IndexByte returns the index of the first instance of c in b, or -1 if c is not present in b.
func IndexByte(b []byte, c byte) int {
	return bytes.IndexByte(b, c)
}

// IndexRune interprets s as a sequence of UTF-8-encoded code points.
// It returns the byte index of the first occurrence in s of the given rune.
// It returns -1 if rune is not present in s.
// If r is utf8.RuneError, it returns the first instance of any
// invalid UTF-8 byte sequence.
func IndexRune(s []byte, r rune) int {
	return bytes.IndexRune(s, r)
}

// IndexFunc interprets s as a sequence of UTF-8-encoded code points.
// It returns the byte index in s of the first Unicode
// code point satisfying f(c), or -1 if none do.
func IndexFunc(s []byte, f func(r rune) bool) int {
	return bytes.IndexFunc(s, f)
}

// Join concatenates the elements of s to create a new byte slice. The separator
// sep is placed between elements in the resulting slice.
func Join(s [][]byte, sep []byte) []byte {
	return bytes.Join(s, sep)
}

// LastIndex returns the index of the last instance of sep in s, or -1 if sep is not present in s.
func LastIndex(s, sep []byte) int {
	return bytes.LastIndex(s, sep)
}

// LastIndexAny interprets s as a sequence of UTF-8-encoded Unicode code
// points. It returns the byte index of the last occurrence in s of any of
// the Unicode code points in chars. It returns -1 if chars is empty or if
// there is no code point in common.
func LastIndexAny(s []byte, chars string) int {
	return bytes.LastIndexAny(s, chars)
}

// LastIndexByte returns the index of the last instance of c in s, or -1 if c is not present in s.
func LastIndexByte(s []byte, c byte) int {
	return bytes.LastIndexByte(s, c)
}

// LastIndexFunc interprets s as a sequence of UTF-8-encoded code points.
// It returns the byte index in s of the last Unicode
// code point satisfying f(c), or -1 if none do.
func LastIndexFunc(s []byte, f func(r rune) bool) int {
	return bytes.LastIndexFunc(s, f)
}

// Map returns a copy of the byte slice s with all its characters modified
// according to the mapping function. If mapping returns a negative value, the character is
// dropped from the byte slice with no replacement. The characters in s and the
// output are interpreted as UTF-8-encoded code points.
func Map(mapping func(r rune) rune, s []byte) []byte {
	return bytes.Map(mapping, s)
}

// Repeat returns a new byte slice consisting of count copies of b.
//
// It panics if count is negative or if
// the result of (len(b) * count) overflows.
func Repeat(b []byte, count int) []byte {
	return bytes.Repeat(b, count)
}

// Replace returns a copy of the slice s with the first n
// non-overlapping instances of old replaced by new.
// If old is empty, it matches at the beginning of the slice
// and after each UTF-8 sequence, yielding up to k+1 replacements
// for a k-rune slice.
// If n < 0, there is no limit on the number of replacements.
func Replace(s, old, new []byte, n int) []byte {
	return bytes.Replace(s, old, new, n)
}

// ReplaceAll returns a copy of the slice s with all
// non-overlapping instances of old replaced by new.
// If old is empty, it matches at the beginning of the slice
// and after each UTF-8 sequence, yielding up to k+1 replacements
// for a k-rune slice.
func ReplaceAll(s, old, new []byte) []byte {
	return bytes.ReplaceAll(s, old, new)
}

// Runes interprets s as a sequence of UTF-8-encoded code points.
// It returns a slice of runes (Unicode code points) equivalent to s.
func Runes(s []byte) []rune {
	return bytes.Runes(s)
}

// Split slices s into all subslices separated by sep and returns a slice of
// the subslices between those separators.
// If sep is empty, Split splits after each UTF-8 sequence.
// It is equivalent to SplitN with a count of -1.
func Split(s, sep []byte) [][]byte {
	return bytes.Split(s, sep)
}

// SplitAfter slices s into all subslices after each instance of sep and
// returns a slice of those subslices.
// If sep is empty, SplitAfter splits after each UTF-8 sequence.
// It is equivalent to SplitAfterN with a count of -1.
func SplitAfter(s, sep []byte) [][]byte {
	return bytes.Split(s, sep)
}

// SplitAfterN slices s into subslices after each instance of sep and
// returns a slice of those subslices.
// If sep is empty, SplitAfterN splits after each UTF-8 sequence.
// The count determines the number of subslices to return:
//
//	n > 0: at most n subslices; the last subslice will be the unsplit remainder.
//	n == 0: the result is nil (zero subslices)
//	n < 0: all subslices
func SplitAfterN(s, sep []byte, n int) [][]byte {
	return bytes.SplitAfterN(s, sep, n)
}

// SplitN slices s into subslices separated by sep and returns a slice of
// the subslices between those separators.
// If sep is empty, SplitN splits after each UTF-8 sequence.
// The count determines the number of subslices to return:
//
//	n > 0: at most n subslices; the last subslice will be the unsplit remainder.
//	n == 0: the result is nil (zero subslices)
//	n < 0: all subslices
func SplitN(s, sep []byte, n int) [][]byte {
	return bytes.SplitN(s, sep, n)
}

// ToLower returns a copy of the byte slice s with all Unicode letters mapped to
// their lower case.
func ToLower(s []byte) []byte {
	return bytes.ToLower(s)
}

// ToLowerSpecial treats s as UTF-8-encoded bytes and returns a copy with all the Unicode letters mapped to their
// lower case, giving priority to the special casing rules.
func ToLowerSpecial(c unicode.SpecialCase, s []byte) []byte {
	return bytes.ToLowerSpecial(c, s)
}

// ToTitle treats s as UTF-8-encoded bytes and returns a copy with all the Unicode letters mapped to their title case.
func ToTitle(s []byte) []byte {
	return bytes.ToTitle(s)
}

// ToTitleSpecial treats s as UTF-8-encoded bytes and returns a copy with all the Unicode letters mapped to their
// title case, giving priority to the special casing rules.
func ToTitleSpecial(c unicode.SpecialCase, s []byte) []byte {
	return bytes.ToTitleSpecial(c, s)
}

// ToUpper returns a copy of the byte slice s with all Unicode letters mapped to
// their upper case.
func ToUpper(s []byte) []byte {
	return bytes.ToUpper(s)
}

// ToUpperSpecial treats s as UTF-8-encoded bytes and returns a copy with all the Unicode letters mapped to their
// upper case, giving priority to the special casing rules.
func ToUpperSpecial(c unicode.SpecialCase, s []byte) []byte {
	return bytes.ToUpperSpecial(c, s)
}

// ToValidUTF8 treats s as UTF-8-encoded bytes and returns a copy with each run of bytes
// representing invalid UTF-8 replaced with the bytes in replacement, which may be empty.
func ToValidUTF8(s, replacement []byte) []byte {
	return bytes.ToValidUTF8(s, replacement)
}

// Trim returns a subslice of s by slicing off all leading and
// trailing UTF-8-encoded code points contained in cutset.
func Trim(s []byte, cutset string) []byte {
	return bytes.Trim(s, cutset)
}

// TrimFunc returns a subslice of s by slicing off all leading and trailing
// UTF-8-encoded code points c that satisfy f(c).
func TrimFunc(s []byte, f func(r rune) bool) []byte {
	return bytes.TrimFunc(s, f)
}

// TrimLeft returns a subslice of s by slicing off all leading
// UTF-8-encoded code points contained in cutset.
func TrimLeft(s []byte, cutset string) []byte {
	return bytes.TrimLeft(s, cutset)
}

// TrimLeftFunc treats s as UTF-8-encoded bytes and returns a subslice of s by slicing off
// all leading UTF-8-encoded code points c that satisfy f(c).
func TrimLeftFunc(s []byte, f func(r rune) bool) []byte {
	return bytes.TrimLeftFunc(s, f)
}

// TrimRight returns a subslice of s by slicing off all trailing
// UTF-8-encoded code points that are contained in cutset.
func TrimRight(s []byte, cutset string) []byte {
	return bytes.TrimRight(s, cutset)
}

// TrimRightFunc returns a subslice of s by slicing off all trailing
// UTF-8-encoded code points c that satisfy f(c).
func TrimRightFunc(s []byte, f func(r rune) bool) []byte {
	return bytes.TrimRightFunc(s, f)
}

// TrimPrefix returns s without the provided leading prefix string.
// If s doesn't start with prefix, s is returned unchanged.
func TrimPrefix(s, prefix []byte) []byte {
	return bytes.TrimPrefix(s, prefix)
}

// TrimSuffix returns s without the provided trailing suffix string.
// If s doesn't end with suffix, s is returned unchanged.
func TrimSuffix(s, suffix []byte) []byte {
	return bytes.TrimSuffix(s, suffix)
}

// TrimSpace returns a subslice of s by slicing off all leading and
// trailing white space, as defined by Unicode.
func TrimSpace(s []byte) []byte {
	return bytes.TrimSpace(s)
}
