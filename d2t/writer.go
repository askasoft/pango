package d2t

import (
	"strings"

	"github.com/askasoft/pango/str"
)

type LineWriter struct {
	sb strings.Builder
}

func (lw *LineWriter) String() string {
	return lw.sb.String()
}

func (lw *LineWriter) WriteString(s string) (n int, err error) {
	lw.sb.WriteString(s)
	lw.sb.WriteByte('\n')
	return len(s) + 1, nil
}

type StripWriter struct {
	sb strings.Builder
}

func (sw *StripWriter) String() string {
	return sw.sb.String()
}

func (sw *StripWriter) WriteString(s string) (n int, err error) {
	s = str.Strip(s)
	if s == "" {
		return 0, nil
	}
	sw.sb.WriteString(s)
	sw.sb.WriteByte('\n')
	return len(s) + 1, nil
}
