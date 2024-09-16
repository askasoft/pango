package sqx

import (
	"bufio"
	"io"
	"strings"
	"unicode"

	"github.com/askasoft/pango/str"
)

type SqlReader struct {
	br  *bufio.Reader
	err error
}

func NewSqlReader(r io.Reader) *SqlReader {
	return &SqlReader{br: bufio.NewReader(r)}
}

func (sr *SqlReader) Error() error {
	return sr.err
}

func (sr *SqlReader) ReadSql() (string, error) {
	if sr.err != nil {
		return "", sr.err
	}

	var sb strings.Builder
	var c rune
	for sr.err == nil {
		c, _, sr.err = sr.br.ReadRune()
		if sr.err != nil {
			break
		}
		if c == '/' {
			c, _, sr.err = sr.br.ReadRune()
			if sr.err != nil {
				sb.WriteRune('/')
				break
			}
			if c == '*' {
				for {
					c, _, sr.err = sr.br.ReadRune()
					if sr.err != nil {
						break
					}
					if c == '*' {
						c, _, sr.err = sr.br.ReadRune()
						if sr.err != nil {
							break
						}
						if c == '/' {
							break
						}
					}
				}
			} else {
				sb.WriteRune('/')
			}
		} else if c == '-' {
			c, _, sr.err = sr.br.ReadRune()
			if sr.err != nil {
				sb.WriteRune('-')
				break
			}
			if c == '-' {
				for {
					c, _, sr.err = sr.br.ReadRune()
					if sr.err != nil {
						break
					}
					if c == '\n' {
						break
					}
				}
			} else {
				sb.WriteRune('-')
				sb.WriteRune(c)
			}
		} else if c == '\'' {
			sb.WriteRune(c)
			for {
				c, _, sr.err = sr.br.ReadRune()
				if sr.err != nil {
					break
				}
				sb.WriteRune(c)

				if c == '\'' {
					c, _, sr.err = sr.br.ReadRune()
					if sr.err != nil {
						break
					}
					if c != '\'' {
						sr.err = sr.br.UnreadRune()
						break
					}
					sb.WriteRune(c)
				}
			}
		} else if unicode.IsSpace(c) {
			// do not append leading space
			if sb.Len() > 0 {
				sb.WriteRune(c)
			}
		} else if c == ';' {
			break
		} else {
			sb.WriteRune(c)
		}
	}

	sql := str.Strip(sb.String())
	if sql != "" {
		return sql, nil
	}
	return "", sr.err
}
