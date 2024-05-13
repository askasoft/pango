package tags

import (
	"html"
	"strings"
)

func RadiosRender(args ...any) (any, error) {
	return TagRender(&RadiosRenderer{}, args...)
}

type RadiosRenderer struct {
	Name     string
	List     List[string, string]
	Value    string
	Disabled bool
}

func (rr *RadiosRenderer) TagName() string {
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
		for it := rr.List.Iterator(); it.Next(); {
			rr.writeRadio(sb, it.Key(), it.Value(), it.Key() == rr.Value)
		}
	}

	TagClose(sb, "div")

	return nil
}

func (rr *RadiosRenderer) writeRadio(sb *strings.Builder, key, text string, checked bool) {
	TagStart(sb, "label")

	a := Attrs{
		"type":  "radio",
		"name":  rr.Name,
		"value": key,
	}

	if checked {
		a.Set("checked", "")
	}
	if rr.Disabled {
		a.Set("disabled", "")
	}

	TagStartClose(sb, "input", a)

	sb.WriteString(html.EscapeString(text))
	TagClose(sb, "label")
}
