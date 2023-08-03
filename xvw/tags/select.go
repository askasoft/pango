package tags

import (
	"html"
	"strings"
)

func SelectRender(args ...any) (any, error) {
	return TagRender(&SelectRenderer{}, args...)
}

type Iterator interface {
	Key() string
	Value() string
	Next() bool
}

type SelectRenderer struct {
	Locale string
	Value  string
	List   Iterator
}

func (sr *SelectRenderer) Name() string {
	return "Select"
}

func (sr *SelectRenderer) Render(sb *strings.Builder, args ...any) error {
	a := Attrs{}

	if err := TagSetAttrs(sr, a, args); err != nil {
		return err
	}

	TagStart(sb, "select", a)

	if sr.List != nil {
		for sr.List.Next() {
			sr.writeOption(sb, sr.List.Key(), sr.List.Value(), sr.List.Key() == sr.Value)
		}
	}

	TagClose(sb, "select")

	return nil
}

func (sr *SelectRenderer) writeOption(sb *strings.Builder, key, text string, selected bool) {
	a := Attrs{}

	a.Set("value", key)
	if selected {
		a.Set("selected", "")
	}

	TagStart(sb, "option", a)
	sb.WriteString(html.EscapeString(text))
	TagClose(sb, "option")
}
