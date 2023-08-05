package tags

import (
	"html"
	"strings"
)

func RadiosRender(args ...any) (any, error) {
	return TagRender(&RadiosRenderer{}, args...)
}

type RadiosRenderer struct {
	Locale string
	Value  string
	List   Iterator
}

func (rr *RadiosRenderer) Name() string {
	return "Radios"
}

func (rr *RadiosRenderer) Render(sb *strings.Builder, args ...any) error {
	a := Attrs{}

	if err := TagSetAttrs(rr, a, args); err != nil {
		return err
	}

	a.Class("ui-radios")

	TagStart(sb, "div", a)

	if rr.List != nil {
		it := rr.List
		for it.Next() {
			rr.writeRadio(sb, it.Key(), it.Value(), it.Key() == rr.Value)
		}
	}

	TagClose(sb, "div")

	return nil
}

func (rr *RadiosRenderer) writeRadio(sb *strings.Builder, key, text string, checked bool) {
	TagStart(sb, "label")

	a := Attrs{}
	a.Set("type", "radio")
	a.Set("value", key)
	if checked {
		a.Set("checked", "")
	}
	TagStartClose(sb, "radio", a)

	sb.WriteString(html.EscapeString(text))
	TagClose(sb, "label")
}
