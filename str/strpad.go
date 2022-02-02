package str

import (
	"strings"
	"unicode/utf8"
)

// RepeatRune returns a new string consisting of count copies of the rune r.
//
// It panics if count is negative or if
// the result of (utf8.RuneLen(r) * count) overflows.
func RepeatRune(r rune, count int) string {
	if count <= 0 {
		return ""
	}

	len := utf8.RuneLen(r)

	// Since we cannot return an error on overflow,
	// we should panic if the repeat will generate
	// an overflow.
	// See Issue golang.org/issue/16237
	size := len * count
	if size/count != len {
		panic("str: RepeatRune count causes overflow")
	}

	sb := &strings.Builder{}
	sb.Grow(size)
	for i := 0; i < count; i++ {
		sb.WriteRune(r)
	}
	return sb.String()
}

// PadCenterRune Center pad a string with the rune r to a new string with len(s) = size.
// str.CenterRune("", 4, ' ')     = "    "
// str.CenterRune("ab", -1, ' ')  = "ab"
// str.CenterRune("ab", 4, ' ')   = " ab "
// str.CenterRune("abcd", 2, ' ') = "abcd"
// str.CenterRune("a", 4, ' ')    = " a  "
// str.CenterRune("a", 4, 'y')    = "yayy"
func PadCenterRune(str string, size int, c rune) string {
	if size <= 0 {
		return str
	}

	strCnt := RuneCount(str)
	pads := size - strCnt
	if pads <= 0 {
		return str
	}

	str = PadLeftRune(str, strCnt+pads/2, c)
	str = PadRightRune(str, size, c)
	return str
}

// PadCenter Center pad a string with the string pad to a new string with len(s) = size.
// str.PadCenter("", 4, " ")     = "    "
// str.PadCenter("ab", -1, " ")  = "ab"
// str.PadCenter("ab", 4, " ")   = " ab "
// str.PadCenter("abcd", 2, " ") = "abcd"
// str.PadCenter("a", 4, " ")    = " a  "
// str.PadCenter("a", 4, "yz")   = "yayz"
// str.PadCenter("abc", 7, "") = "  abc  "
func PadCenter(str string, size int, pad string) string {
	if size <= 0 {
		return str
	}

	if pad == "" {
		return PadCenterRune(str, size, ' ')
	}

	strCnt := RuneCount(str)
	pads := size - strCnt
	if pads <= 0 {
		return str
	}

	str = PadLeft(str, strCnt+pads/2, pad)
	str = PadRight(str, size, pad)
	return str
}

// PadLeftRune left pad the string str with the rune r to a new string with len(s) = size.
// str.PadLeftRune("", 3, 'z')     = "zzz"
// str.PadLeftRune("bat", 3, 'z')  = "bat"
// str.PadLeftRune("bat", 5, 'z')  = "zzbat"
// str.PadLeftRune("bat", 1, 'z')  = "bat"
// str.PadLeftRune("bat", -1, 'z') = "bat"
func PadLeftRune(str string, size int, r rune) string {
	size -= RuneCount(str)
	if size <= 0 {
		return str
	}
	return RepeatRune(r, size) + str
}

// PadLeft left pad the string str with the string pad to a new string with len(s) = size.
// str.PadLeft("", 3, "z")      = "zzz"
// str.PadLeft("bat", 3, "yz")  = "bat"
// str.PadLeft("bat", 5, "yz")  = "yzbat"
// str.PadLeft("bat", 8, "yz")  = "yzyzybat"
// str.PadLeft("bat", 1, "yz")  = "bat"
// str.PadLeft("bat", -1, "yz") = "bat"
// str.PadLeft("bat", 5, "")    = "bat"
func PadLeft(str string, size int, pad string) string {
	if pad == "" {
		return str
	}

	size -= RuneCount(str)
	if size <= 0 {
		return str
	}

	padCnt := RuneCount(pad)
	if padCnt == 1 {
		r, _ := utf8.DecodeRuneInString(pad)
		return RepeatRune(r, size) + str
	}

	if size == padCnt {
		return pad + str
	}

	sb := &strings.Builder{}
	sb.Grow(len(pad)*((size+size%padCnt)/padCnt) + len(str))
	padstr(sb, pad, padCnt, size)
	sb.WriteString(str)

	return sb.String()
}

func padstr(sb *strings.Builder, pad string, padCnt int, count int) {
	for count >= padCnt {
		sb.WriteString(pad)
		count -= padCnt
	}
	if count > 0 {
		for _, c := range pad {
			sb.WriteRune(c)
			count--
			if count == 0 {
				break
			}
		}
	}
}

// PadRightRune right pad the string str with the rune r to a new string with len(s) = size.
// str.PadRightRune("", 3, 'z')     = "zzz"
// str.PadRightRune("bat", 3, 'z')  = "bat"
// str.PadRightRune("bat", 5, 'z')  = "batzz"
// str.PadRightRune("bat", 1, 'z')  = "bat"
// str.PadRightRune("bat", -1, 'z') = "bat"
func PadRightRune(str string, size int, r rune) string {
	size -= RuneCount(str)
	if size <= 0 {
		return str
	}
	return str + RepeatRune(r, size)
}

// PadRight right pad the string str with the string pad to a new string with len(s) = size.
// str.PadRight("", 3, "z")      = "zzz"
// str.PadRight("bat", 3, "yz")  = "bat"
// str.PadRight("bat", 5, "yz")  = "batyz"
// str.PadRight("bat", 8, "yz")  = "batyzyzy"
// str.PadRight("bat", 1, "yz")  = "bat"
// str.PadRight("bat", -1, "yz") = "bat"
// str.PadRight("bat", 5, "")    = "bat"
func PadRight(str string, size int, pad string) string {
	if pad == "" {
		return str
	}

	size -= RuneCount(str)
	if size <= 0 {
		return str
	}

	padCnt := RuneCount(pad)
	if padCnt == 1 {
		return str + Repeat(pad, size)
	}

	if size == padCnt {
		return str + pad
	}

	sb := &strings.Builder{}
	sb.Grow(len(str) + padCnt*((size+size%padCnt)/padCnt))
	sb.WriteString(str)
	padstr(sb, pad, padCnt, size)

	return sb.String()
}
