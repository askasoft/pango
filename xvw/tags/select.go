package tags

import (
	"html"
	"strings"
)

func SelectRender(args ...any) (any, error) {
	return TagRender(&SelectRenderer{}, args...)
}

type SelectRenderer struct {
	Name  string
	List  Iterator
	Value string
}

func (sr *SelectRenderer) TagName() string {
	return "Select"
}

func (sr *SelectRenderer) Render(sb *strings.Builder, args ...any) error {
	a := Attrs{}

	if err := TagSetAttrs(sr, a, args); err != nil {
		return err
	}
	a.Name(sr.Name)

	TagStart(sb, "select", a)

	if sr.List != nil {
		it := sr.List
		for it.Next() {
			sr.writeOption(sb, it.Key(), it.Value(), it.Key() == sr.Value)
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
