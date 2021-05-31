package str

import (
	"strings"
)

// PadLeftRune left pad a string with a specified character.
// str.PadLeftRune("", 3, 'z')     = "zzz"
// str.PadLeftRune("bat", 3, 'z')  = "bat"
// str.PadLeftRune("bat", 5, 'z')  = "zzbat"
// str.PadLeftRune("bat", 1, 'z')  = "bat"
// str.PadLeftRune("bat", -1, 'z') = "bat"
func PadLeftRune(str string, size int, ch rune) string {
	psz := size - RuneCount(str)
	if psz <= 0 {
		return str
	}
	return Repeat(string(ch), psz) + str
}

// PadLeft Left pad a string with a specified String.
// Strings.leftPad("", 3, "z")      = "zzz"
// Strings.leftPad("bat", 3, "yz")  = "bat"
// Strings.leftPad("bat", 5, "yz")  = "yzbat"
// Strings.leftPad("bat", 8, "yz")  = "yzyzybat"
// Strings.leftPad("bat", 1, "yz")  = "bat"
// Strings.leftPad("bat", -1, "yz") = "bat"
// Strings.leftPad("bat", 5, "")    = "bat"
func PadLeft(str string, size int, ps string) string {
	if ps == "" {
		return str
	}

	padLen := RuneCount(ps)
	strLen := RuneCount(str)
	psz := size - strLen
	if psz <= 0 {
		return str
	}
	if padLen == 1 {
		return Repeat(ps, psz) + str
	}

	if psz == padLen {
		return ps + str
	}

	sb := &strings.Builder{}
	sb.Grow(len(ps)*((psz+psz%padLen)/padLen) + len(str))
	pad(sb, ps, padLen, psz)
	sb.WriteString(str)

	return sb.String()
}

func pad(sb *strings.Builder, ps string, padLen int, psz int) {
	for psz >= padLen {
		sb.WriteString(ps)
		psz -= padLen
	}
	if psz > 0 {
		for _, c := range ps {
			sb.WriteRune(c)
			psz--
			if psz == 0 {
				break
			}
		}
	}
}

// PadRightRune Right pad a string with a specified character.
// str.PadRightRune("", 3, 'z')     = "zzz"
// str.PadRightRune("bat", 3, 'z')  = "bat"
// str.PadRightRune("bat", 5, 'z')  = "batzz"
// str.PadRightRune("bat", 1, 'z')  = "bat"
// str.PadRightRune("bat", -1, 'z') = "bat"
func PadRightRune(str string, size int, ch rune) string {
	psz := size - RuneCount(str)
	if psz <= 0 {
		return str
	}
	return str + Repeat(string(ch), psz)
}

// PadRight Right pad a string with a specified String.
// str.PadRight("", 3, "z")      = "zzz"
// str.PadRight("bat", 3, "yz")  = "bat"
// str.PadRight("bat", 5, "yz")  = "batyz"
// str.PadRight("bat", 8, "yz")  = "batyzyzy"
// str.PadRight("bat", 1, "yz")  = "bat"
// str.PadRight("bat", -1, "yz") = "bat"
// str.PadRight("bat", 5, "")    = "bat"
func PadRight(str string, size int, ps string) string {
	if ps == "" {
		return str
	}

	padLen := RuneCount(ps)
	strLen := RuneCount(str)
	psz := size - strLen
	if psz <= 0 {
		return str
	}
	if padLen == 1 {
		return str + Repeat(ps, psz)
	}

	if psz == padLen {
		return str + ps
	}

	sb := &strings.Builder{}
	sb.Grow(len(str) + padLen*((psz+psz%padLen)/padLen))
	sb.WriteString(str)
	pad(sb, ps, padLen, psz)

	return sb.String()
}
