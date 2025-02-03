package rex

import (
	"io"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

func SubmatchReplace(re *regexp.Regexp, source, replace []byte) []byte {
	smis := re.FindAllSubmatchIndex(source, -1)
	if len(smis) == 0 {
		return source
	}

	last := 0

	var bs []byte
	for _, smi := range smis {
		if smi[0] > last {
			bs = append(bs, source[last:smi[0]]...)
		}
		last = smi[1]

		bs = re.Expand(bs, replace, source, smi)
	}

	if last < len(source) {
		bs = append(bs, source[last:]...)
	}

	return bs
}

func SubmatchReplaceString(re *regexp.Regexp, source, replace string) string {
	sr := SubmatchReplacer{Pattern: re, Template: replace}
	return sr.Replace(source)
}

type Converter func(n int, name string, value string) string

func noconvert(n int, name string, value string) string {
	return value
}

type SubmatchReplacer struct {
	Pattern   *regexp.Regexp
	Template  string
	Converter Converter
}

// Replace returns a copy of s with all replacements performed.
func (sr *SubmatchReplacer) Replace(source string) string {
	smis := sr.Pattern.FindAllStringSubmatchIndex(source, -1)
	if len(smis) == 0 {
		return source
	}

	var sb strings.Builder
	if _, err := sr.write(&sb, source, smis); err != nil {
		panic(err)
	}

	return sb.String()
}

// WriteString writes source to w with all replacements performed.
func (sr *SubmatchReplacer) WriteString(w io.Writer, source string) (int, error) {
	smis := sr.Pattern.FindAllStringSubmatchIndex(source, -1)
	if len(smis) == 0 {
		return io.WriteString(w, source)
	}

	return sr.write(w, source, smis)
}

func (sr *SubmatchReplacer) write(w io.Writer, source string, smis [][]int) (int, error) {
	convert := sr.Converter
	if convert == nil {
		convert = noconvert
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

		n, err := sr.expand(w, source, smi, convert)
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

func (sr *SubmatchReplacer) expand(dst io.Writer, src string, match []int, convert Converter) (int, error) {
	written := 0
	template := sr.Template

	for len(template) > 0 {
		before, after, ok := strings.Cut(template, "$")
		if !ok {
			break
		}

		n, err := io.WriteString(dst, before)
		written += n
		if err != nil {
			return written, err
		}

		template = after
		if template != "" && template[0] == '$' {
			// Treat $$ as $.
			n, err := io.WriteString(dst, "$")
			written += n
			if err != nil {
				return written, err
			}
			template = template[1:]
			continue
		}

		name, num, rest, ok := sr.extract(template)
		if !ok {
			// Malformed; treat $ as raw text.
			n, err := io.WriteString(dst, "$")
			written += n
			if err != nil {
				return written, err
			}
			continue
		}

		template = rest
		if num >= 0 {
			if 2*num+1 < len(match) && match[2*num] >= 0 {
				sub := convert(num, "", src[match[2*num]:match[2*num+1]])
				n, err := io.WriteString(dst, sub)
				written += n
				if err != nil {
					return written, err
				}
			}
		} else {
			for i, namei := range sr.Pattern.SubexpNames() {
				if name == namei && 2*i+1 < len(match) && match[2*i] >= 0 {
					sub := convert(-1, name, src[match[2*i]:match[2*i+1]])
					n, err := io.WriteString(dst, sub)
					written += n
					if err != nil {
						return written, err
					}
					break
				}
			}
		}
	}

	n, err := io.WriteString(dst, template)
	written += n
	return written, err
}

// extract returns the name from a leading "name" or "{name}" in str.
// (The $ has already been removed by the caller.)
// If it is a number, extract returns num set to that number; otherwise num = -1.
func (sr *SubmatchReplacer) extract(str string) (name string, num int, rest string, ok bool) {
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
