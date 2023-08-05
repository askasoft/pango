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
	TagName() string
	Render(sb *strings.Builder, args ...any) error
}

func TagRender(tr TagRenderer, args ...any) (any, error) {
	sb := &strings.Builder{}
	err := tr.Render(sb, args...)
	return template.HTML(sb.String()), err //nolint: gosec
}

func TagSetAttrs(tr TagRenderer, attrs Attrs, args []any) error {
	if len(args)&1 != 0 {
		return errors.New(tr.TagName() + ": invalid arguments")
	}

	for i := 0; i < len(args); i += 2 {
		key, ok := args[i].(string)
		if !ok || key == "" {
			return errors.New(tr.TagName() + ": key must be non-empty string")
		}

		val := args[i+1]
		if val == nil {
			continue
		}

		if unicode.IsUpper(rune(key[0])) {
			err := ref.SetProperty(tr, key, val)
			if err != nil {
				return err
			}
			continue
		}

		if sv, ok := val.(string); ok {
			attrs.Set(key, sv)
		} else {
			attrs.Set(key, fmt.Sprint(val))
		}
	}
	return nil
}

func TagStart(t *strings.Builder, name string, as ...Attrs) {
	t.WriteByte('<')
	t.WriteString(name)
	for _, a := range as {
		TagWriteAttrs(t, a)
	}
	t.WriteByte('>')
}

func TagWriteAttrs(t *strings.Builder, a Attrs) {
	for k, v := range a {
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

func TagStartClose(t *strings.Builder, name string, as ...Attrs) {
	t.WriteByte('<')
	t.WriteString(name)
	for _, a := range as {
		TagWriteAttrs(t, a)
	}
	t.WriteString("/>")
}

func TagClose(t *strings.Builder, name string) {
	t.WriteString("</")
	t.WriteString(name)
	t.WriteByte('>')
}
