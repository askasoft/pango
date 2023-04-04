package str

import (
	"strings"
	"unicode/utf8"
)

// RepeatByte returns a new string consisting of count copies of the byte b.
//
// It panics if count is negative or if
func RepeatByte(b byte, count int) string {
	if count <= 0 {
		return ""
	}

	sb := &strings.Builder{}
	sb.Grow(count)
	for i := 0; i < count; i++ {
		sb.WriteByte(b)
	}
	return sb.String()
}

// RepeatRune returns a new string consisting of count copies of the rune r.
//
// It panics if count is negative or if
// the result of (utf8.RuneLen(r) * count) overflows.
func RepeatRune(r rune, count int) string {
	if count <= 0 {
		return ""
	}

	rlen := utf8.RuneLen(r)

	// Since we cannot return an error on overflow,
	// we should panic if the repeat will generate
	// an overflow.
	// See Issue golang.org/issue/16237
	size := rlen * count
	if size/count != rlen {
		panic("str: RepeatRune count causes overflow")
	}

	sb := &strings.Builder{}
	sb.Grow(size)
	for i := 0; i < count; i++ {
		sb.WriteRune(r)
	}
	return sb.String()
}

// PadCenterByte Center pad a string with the byte b to a new string with len(s) = size.
// str.PadCenterByte("", 4, ' ')     = "    "
// str.PadCenterByte("ab", -1, ' ')  = "ab"
// str.PadCenterByte("ab", 4, ' ')   = " ab "
// str.PadCenterByte("abcd", 2, ' ') = "abcd"
// str.PadCenterByte("a", 4, ' ')    = " a  "
// str.PadCenterByte("a", 4, 'y')    = "yayy"
func PadCenterByte(s string, size int, b byte) string {
	if size <= 0 {
		return s
	}

	cnt := len(s)
	pad := size - cnt
	if pad <= 0 {
		return s
	}

	s = PadLeftByte(s, cnt+pad/2, b)
	s = PadRightByte(s, size, b)
	return s
}

// PadCenterRune Center pad a string with the rune r to a new string with RuneCount(s) = size.
// str.PadCenterRune("", 4, ' ')     = "    "
// str.PadCenterRune("ab", -1, ' ')  = "ab"
// str.PadCenterRune("ab", 4, ' ')   = " ab "
// str.PadCenterRune("abcd", 2, ' ') = "abcd"
// str.PadCenterRune("a", 4, ' ')    = " a  "
// str.PadCenterRune("a", 4, 'y')    = "yayy"
func PadCenterRune(s string, size int, c rune) string {
	if size <= 0 {
		return s
	}

	cnt := RuneCount(s)
	pad := size - cnt
	if pad <= 0 {
		return s
	}

	s = PadLeftRune(s, cnt+pad/2, c)
	s = PadRightRune(s, size, c)
	return s
}

// PadCenter Center pad a string with the string p to a new string with RuneCount(s) = size.
// str.PadCenter("", 4, " ")     = "    "
// str.PadCenter("ab", -1, " ")  = "ab"
// str.PadCenter("ab", 4, " ")   = " ab "
// str.PadCenter("abcd", 2, " ") = "abcd"
// str.PadCenter("a", 4, " ")    = " a  "
// str.PadCenter("a", 4, "yz")   = "yayz"
// str.PadCenter("abc", 7, "") = "abc"
func PadCenter(s string, size int, p string) string {
	if p == "" {
		return s
	}

	if size <= 0 {
		return s
	}

	cnt := RuneCount(s)
	pad := size - cnt
	if pad <= 0 {
		return s
	}

	s = PadLeft(s, cnt+pad/2, p)
	s = PadRight(s, size, p)
	return s
}

// PadLeftByte left pad the string s with the byte b to a new string with len(s) = size.
// str.PadLeftByte("", 3, 'z')     = "zzz"
// str.PadLeftByte("bat", 3, 'z')  = "bat"
// str.PadLeftByte("bat", 5, 'z')  = "zzbat"
// str.PadLeftByte("bat", 1, 'z')  = "bat"
// str.PadLeftByte("bat", -1, 'z') = "bat"
func PadLeftByte(s string, size int, b byte) string {
	size -= len(s)
	if size <= 0 {
		return s
	}
	return RepeatByte(b, size) + s
}

// PadLeftRune left pad the string s with the rune r to a new string with RuneCount(s) = size.
// str.PadLeftRune("", 3, 'z')     = "zzz"
// str.PadLeftRune("bat", 3, 'z')  = "bat"
// str.PadLeftRune("bat", 5, 'z')  = "zzbat"
// str.PadLeftRune("bat", 1, 'z')  = "bat"
// str.PadLeftRune("bat", -1, 'z') = "bat"
func PadLeftRune(s string, size int, r rune) string {
	size -= RuneCount(s)
	if size <= 0 {
		return s
	}
	return RepeatRune(r, size) + s
}

// PadLeft left pad the string s with the string p to a new string with RuneCount(s) = size.
// str.PadLeft("", 3, "z")      = "zzz"
// str.PadLeft("bat", 3, "yz")  = "bat"
// str.PadLeft("bat", 5, "yz")  = "yzbat"
// str.PadLeft("bat", 8, "yz")  = "yzyzybat"
// str.PadLeft("bat", 1, "yz")  = "bat"
// str.PadLeft("bat", -1, "yz") = "bat"
// str.PadLeft("bat", 5, "")    = "bat"
func PadLeft(s string, size int, p string) string {
	if p == "" {
		return s
	}

	size -= RuneCount(s)
	if size <= 0 {
		return s
	}

	pad := RuneCount(p)
	if pad == 1 {
		r, _ := utf8.DecodeRuneInString(p)
		return RepeatRune(r, size) + s
	}

	if size == pad {
		return p + s
	}

	sb := &strings.Builder{}
	sb.Grow(len(p)*((size+size%pad)/pad) + len(s))
	padstr(sb, p, pad, size)
	sb.WriteString(s)

	return sb.String()
}

func padstr(sb *strings.Builder, p string, pad int, size int) {
	for size >= pad {
		sb.WriteString(p)
		size -= pad
	}
	if size > 0 {
		for _, c := range p {
			sb.WriteRune(c)
			size--
			if size == 0 {
				break
			}
		}
	}
}

// PadRightByte right pad the string s with the byte b to a new string with len(s) = size.
// str.PadRightByte("", 3, 'z')     = "zzz"
// str.PadRightByte("bat", 3, 'z')  = "bat"
// str.PadRightByte("bat", 5, 'z')  = "batzz"
// str.PadRightByte("bat", 1, 'z')  = "bat"
// str.PadRightByte("bat", -1, 'z') = "bat"
func PadRightByte(s string, size int, b byte) string {
	size -= len(s)
	if size <= 0 {
		return s
	}
	return s + RepeatByte(b, size)
}

// PadRightRune right pad the string s with the rune r to a new string with RuneCount(s) = size.
// str.PadRightRune("", 3, 'z')     = "zzz"
// str.PadRightRune("bat", 3, 'z')  = "bat"
// str.PadRightRune("bat", 5, 'z')  = "batzz"
// str.PadRightRune("bat", 1, 'z')  = "bat"
// str.PadRightRune("bat", -1, 'z') = "bat"
func PadRightRune(s string, size int, r rune) string {
	size -= RuneCount(s)
	if size <= 0 {
		return s
	}
	return s + RepeatRune(r, size)
}

// PadRight right pad the string s with the string p to a new string with RuneCount(s) = size.
// str.PadRight("", 3, "z")      = "zzz"
// str.PadRight("bat", 3, "yz")  = "bat"
// str.PadRight("bat", 5, "yz")  = "batyz"
// str.PadRight("bat", 8, "yz")  = "batyzyzy"
// str.PadRight("bat", 1, "yz")  = "bat"
// str.PadRight("bat", -1, "yz") = "bat"
// str.PadRight("bat", 5, "")    = "bat"
func PadRight(s string, size int, p string) string {
	if p == "" {
		return s
	}

	size -= RuneCount(s)
	if size <= 0 {
		return s
	}

	pad := RuneCount(p)
	if pad == 1 {
		return s + Repeat(p, size)
	}

	if size == pad {
		return s + p
	}

	sb := &strings.Builder{}
	sb.Grow(len(s) + pad*((size+size%pad)/pad))
	sb.WriteString(s)
	padstr(sb, p, pad, size)

	return sb.String()
}
