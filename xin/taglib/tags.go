package taglib

import (
	"errors"
	"fmt"
	"html"
	"html/template"
	"strings"
	"unicode"

	"github.com/askasoft/pango/ref"
	"github.com/askasoft/pango/str"
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
	for i := 0; i < len(args); i++ {
		key, ok := args[i].(string)
		if !ok || key == "" {
			return errors.New(tr.TagName() + ": key must be non-empty string")
		}

		var val any
		if str.EndsWithByte(key, '=') {
			key = key[:len(key)-1]
			if key == "" {
				return errors.New(tr.TagName() + ": key must be non-empty string")
			}

			i++
			if i < len(args) {
				val = args[i]
			}
		}

		if unicode.IsUpper(rune(key[0])) {
			if val != nil {
				if err := ref.SetProperty(tr, key, val); err != nil {
					return err
				}
			}
			continue
		}

		if val == nil {
			attrs.Set(key, "")
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

func TagStart(sb *strings.Builder, name string, as ...Attrs) {
	sb.WriteByte('<')
	sb.WriteString(name)
	for _, a := range as {
		TagWriteAttrs(sb, a)
	}
	sb.WriteByte('>')
}

func TagWriteAttrs(sb *strings.Builder, a Attrs) {
	for k, v := range a {
		if k != "" {
			sb.WriteByte(' ')
			sb.WriteString(k)
			if v != "" {
				sb.WriteByte('=')
				sb.WriteByte('"')
				sb.WriteString(html.EscapeString(v))
				sb.WriteByte('"')
			}
		}
	}
}

func TagStartClose(sb *strings.Builder, name string, as ...Attrs) {
	sb.WriteByte('<')
	sb.WriteString(name)
	for _, a := range as {
		TagWriteAttrs(sb, a)
	}
	sb.WriteString("/>")
}

func TagClose(sb *strings.Builder, name string) {
	sb.WriteString("</")
	sb.WriteString(name)
	sb.WriteByte('>')
}
