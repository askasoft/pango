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

func (rr *RadiosRenderer) SetValue(value any) {
	rr.Value = toString(value)
}

func (rr *RadiosRenderer) Render(sb *strings.Builder, args ...any) error {
	a := Attrs{}

	if err := TagSetAttrs(rr, a, args); err != nil {
		return err
	}

	a.Class("ui-radios")

	TagStart(sb, "div", a)

	if rr.List != nil {
		rr.List.Each(func(k any, v string) bool {
			rr.writeRadio(sb, k, v)
			return true
		})
	}

	TagClose(sb, "div")

	return nil
}

func (rr *RadiosRenderer) writeRadio(sb *strings.Builder, key any, text string) {
	TagStart(sb, "label")

	sk := toString(key)

	a := Attrs{
		"type":  "radio",
		"name":  rr.Name,
		"value": sk,
	}

	if sk == rr.Value {
		a.Set("checked", "")
	}
	if rr.Disabled {
		a.Set("disabled", "")
	}

	TagStartClose(sb, "input", a)

	sb.WriteString(html.EscapeString(text))
	TagClose(sb, "label")
}
