package rex

import (
	"io"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Converter a convert function used by ExpandConvertWrite()
// `n` is the number of $0, $1, $2, ...
// `name` is the name of ${x}, $y, ...
// (n = -1 and name = "") means other unmatched text
type Converter func(n int, name string, value string) string

func NoConvert(n int, name string, value string) string {
	return value
}

func ReplaceAllString(source string, re *regexp.Regexp, template string) string {
	return ReplaceAllConvertString(source, re, template, NoConvert)
}

func ReplaceAllConvertString(source string, re *regexp.Regexp, template string, convert Converter) string {
	var sb strings.Builder

	if _, err := ReplaceAllConvertWrite(source, &sb, re, template, convert); err != nil {
		panic(err)
	}

	return sb.String()
}

func ReplaceAllWrite(source string, w io.Writer, re *regexp.Regexp, template string) (int, error) {
	return ReplaceAllConvertWrite(source, w, re, template, NoConvert)
}

func ReplaceAllConvertWrite(source string, w io.Writer, re *regexp.Regexp, template string, convert Converter) (int, error) {
	smis := re.FindAllStringSubmatchIndex(source, -1)
	if len(smis) == 0 {
		return io.WriteString(w, source)
	}

	last := 0
	written := 0

	for _, smi := range smis {
		if smi[0] > last {
			sub := convert(-1, "", source[last:smi[0]])
			n, err := io.WriteString(w, sub)
			written += n
			if err != nil {
				return written, err
			}
		}
		last = smi[1]

		n, err := ExpandConvertWrite(source, w, re, template, smi, convert)
		written += n
		if err != nil {
			return written, err
		}
	}

	if last < len(source) {
		sub := convert(-1, "", source[last:])
		n, err := io.WriteString(w, sub)
		written += n
		if err != nil {
			return written, err
		}
	}

	return written, nil
}

func ExpandWrite(source string, w io.Writer, re *regexp.Regexp, template string, match []int) (int, error) {
	return ExpandConvertWrite(source, w, re, template, match, NoConvert)
}

func ExpandConvertWrite(source string, w io.Writer, re *regexp.Regexp, template string, match []int, convert Converter) (int, error) {
	written := 0

	for len(template) > 0 {
		before, after, ok := strings.Cut(template, "$")
		if !ok {
			break
		}

		n, err := io.WriteString(w, before)
		written += n
		if err != nil {
			return written, err
		}

		template = after
		if template != "" && template[0] == '$' {
			// Treat $$ as $.
			n, err := io.WriteString(w, "$")
			written += n
			if err != nil {
				return written, err
			}
			template = template[1:]
			continue
		}

		name, num, rest, ok := extract(template)
		if !ok {
			// Malformed; treat $ as raw text.
			n, err := io.WriteString(w, "$")
			written += n
			if err != nil {
				return written, err
			}
			continue
		}

		template = rest
		if num >= 0 {
			if 2*num+1 < len(match) && match[2*num] >= 0 {
				sub := convert(num, "", source[match[2*num]:match[2*num+1]])
				n, err := io.WriteString(w, sub)
				written += n
				if err != nil {
					return written, err
				}
			}
		} else {
			for i, namei := range re.SubexpNames() {
				if name == namei && 2*i+1 < len(match) && match[2*i] >= 0 {
					sub := convert(-1, name, source[match[2*i]:match[2*i+1]])
					n, err := io.WriteString(w, sub)
					written += n
					if err != nil {
						return written, err
					}
					break
				}
			}
		}
	}

	n, err := io.WriteString(w, template)
	written += n
	return written, err
}

// extract returns the name from a leading "name" or "{name}" in str.
// (The $ has already been removed by the caller.)
// If it is a number, extract returns num set to that number; otherwise num = -1.
func extract(str string) (name string, num int, rest string, ok bool) {
	if str == "" {
		return
	}
	brace := false
	if str[0] == '{' {
		brace = true
		str = str[1:]
	}
	i := 0
	for i < len(str) {
		rune, size := utf8.DecodeRuneInString(str[i:])
		if !unicode.IsLetter(rune) && !unicode.IsDigit(rune) && rune != '_' {
			break
		}
		i += size
	}
	if i == 0 {
		// empty name is not okay
		return
	}
	name = str[:i]
	if brace {
		if i >= len(str) || str[i] != '}' {
			// missing closing brace
			return
		}
		i++
	}

	// Parse number.
	num = 0
	for i := 0; i < len(name); i++ {
		if name[i] < '0' || '9' < name[i] || num >= 1e8 {
			num = -1
			break
		}
		num = num*10 + int(name[i]) - '0'
	}
	// Disallow leading zeros.
	if name[0] == '0' && len(name) > 1 {
		num = -1
	}

	rest = str[i:]
	ok = true
	return
}
