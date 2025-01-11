package str

import (
	"fmt"
	"io"
	"regexp"
)

func RegexpSubmatchReplace(re *regexp.Regexp, source, replace string) string {
	var sb Builder
	if _, err := RegexpSubmatchReplaceWrite(&sb, re, source, replace); err != nil {
		return source
	}
	return sb.String()
}

func RegexpSubmatchReplaceWrite(w io.Writer, re *regexp.Regexp, source, replace string) (int, error) {
	smis := re.FindAllStringSubmatchIndex(source, -1)
	if len(smis) == 0 {
		return io.WriteString(w, source)
	}

	written := 0
	last := 0
	for _, smi := range smis {
		if smi[0] > last {
			n, err := io.WriteString(w, source[last:smi[0]])
			written += n
			if err != nil {
				return written, err
			}
		}

		rps := make([]string, 0, len(smi)-2)
		for i := 2; i < len(smi); i += 2 {
			rps = append(rps, fmt.Sprintf("$%d", i>>1), source[smi[i]:smi[i+1]])
		}

		n, err := NewReplacer(rps...).WriteString(w, replace)
		written += n
		if err != nil {
			return written, err
		}

		last = smi[1]
	}

	if last < len(source) {
		n, err := io.WriteString(w, source[last:])
		written += n
		if err != nil {
			return written, err
		}
	}

	return written, nil
}

type RegexpSubmatchReplacer struct {
	Patterns []*regexp.Regexp
	Replaces []string
}

func NewRegexpSubmatchReplacer(patreps ...string) (*RegexpSubmatchReplacer, error) {
	if len(patreps)%2 == 1 {
		panic("str.NewRegexpSubmatchReplacer: odd argument count")
	}

	patterns := make([]*regexp.Regexp, 0, len(patreps)/2)
	replaces := make([]string, 0, len(patterns))

	for i := 0; i < len(patreps); i += 2 {
		p, err := regexp.Compile(patreps[i])
		if err != nil {
			return nil, err
		}
		patterns = append(patterns, p)
		replaces = append(replaces, patreps[i+1])
	}

	return &RegexpSubmatchReplacer{Patterns: patterns, Replaces: replaces}, nil
}

func (rsr *RegexpSubmatchReplacer) Replace(s string) string {
	if len(rsr.Patterns) == 0 {
		return s
	}

	var sb Builder
	for i, p := range rsr.Patterns {
		sb.Reset()
		if _, err := RegexpSubmatchReplaceWrite(&sb, p, s, rsr.Replaces[i]); err != nil {
			return s
		}
		s = sb.String()
	}
	return s
}

func (rsr *RegexpSubmatchReplacer) WriteString(w io.Writer, s string) (int, error) {
	n := len(rsr.Patterns)
	if n == 0 {
		return io.WriteString(w, s)
	}

	if n > 1 {
		var sb Builder
		for i, p := range rsr.Patterns[0 : n-1] {
			sb.Reset()
			if _, err := RegexpSubmatchReplaceWrite(&sb, p, s, rsr.Replaces[i]); err != nil {
				return 0, err
			}
			s = sb.String()
		}
	}

	return RegexpSubmatchReplaceWrite(w, rsr.Patterns[n-1], s, rsr.Replaces[n-1])
}
