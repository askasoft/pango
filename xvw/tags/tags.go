package tags

import (
	"errors"
	"fmt"
	"html"
	"html/template"
	"strings"
	"unicode"

	"github.com/askasoft/pango/ref"
)

type TagRenderer interface {
	Name() string
	Render(sb *strings.Builder, args ...any) error
}

func TagRender(tr TagRenderer, args ...any) (any, error) {
	sb := &strings.Builder{}
	err := tr.Render(sb, args...)
	return template.HTML(sb.String()), err //nolint: gosec
}

func TagSetAttrs(tr TagRenderer, attrs Attrs, args []any) error {
	if len(args)&1 != 0 {
		return errors.New(tr.Name() + ": invalid arguments")
	}

	for i := 0; i < len(args); i += 2 {
		key, ok := args[i].(string)
		if !ok || key == "" {
			return errors.New(tr.Name() + ": key must be non-empty string")
		}

		if unicode.IsUpper(rune(key[0])) {
			err := ref.SetProperty(tr, key, args[i+1])
			if err != nil {
				return err
			}
			continue
		}

		val := args[i+1]
		if sv, ok := val.(string); ok {
			attrs.Set(key, sv)
		} else {
			attrs.Set(key, fmt.Sprint(val))
		}
	}
	return nil
}

func TagStart(t *strings.Builder, name string, ass ...Attrs) {
	t.WriteByte('<')
	t.WriteString(name)
	for _, as := range ass {
		for k, v := range as {
			if k != "" {
				t.WriteByte(' ')
				t.WriteString(k)
				if v != "" {
					t.WriteByte('=')
					t.WriteByte('"')
					t.WriteString(html.EscapeString(v))
					t.WriteByte('"')
				}
			}
		}
	}
	t.WriteByte('>')
}

func TagStartClose(t *strings.Builder, name string, ass ...Attrs) {
	TagStart(t, name, ass...)
	t.WriteString("/>")
}

func TagClose(t *strings.Builder, name string) {
	t.WriteString("</")
	t.WriteString(name)
	t.WriteByte('>')
}
