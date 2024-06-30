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
	List     List
	Value    string
	Disabled bool
}

func (rr *RadiosRenderer) TagName() string {
	return "Radios"
}

func (rr *RadiosRenderer) SetList(list any) {
	rr.List = AsList(list)
}

func (rr *RadiosRenderer) Render(sb *strings.Builder, args ...any) error {
	a := Attrs{}

	if err := TagSetAttrs(rr, a, args); err != nil {
		return err
	}

	a.Class("ui-radios")

	TagStart(sb, "div", a)

	if rr.List != nil {
		rr.List.Each(func(k, v string) bool {
			rr.writeRadio(sb, k, v)
			return true
		})
	}

	TagClose(sb, "div")

	return nil
}

func (rr *RadiosRenderer) writeRadio(sb *strings.Builder, key, text string) {
	TagStart(sb, "label")

	a := Attrs{
		"type":  "radio",
		"name":  rr.Name,
		"value": key,
	}

	if key == rr.Value {
		a.Set("checked", "")
	}
	if rr.Disabled {
		a.Set("disabled", "")
	}

	TagStartClose(sb, "input", a)

	sb.WriteString(html.EscapeString(text))
	TagClose(sb, "label")
}
