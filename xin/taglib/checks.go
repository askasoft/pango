package taglib

import (
	"html"
	"strings"
)

func ChecksRender(args ...any) (any, error) {
	return TagRender(&ChecksRenderer{}, args...)
}

type ChecksRenderer struct {
	Name     string
	List     List
	Values   Values
	Disabled bool
	Ordered  bool
}

func (cr *ChecksRenderer) TagName() string {
	return "Checks"
}

func (cr *ChecksRenderer) SetList(list any) {
	cr.List = AsList(list)
}

func (cr *ChecksRenderer) SetValues(values any) {
	cr.Values = AsValues(values)
}

func (cr *ChecksRenderer) Render(sb *strings.Builder, args ...any) error {
	a := Attrs{}

	if err := TagSetAttrs(cr, a, args); err != nil {
		return err
	}

	a.Class("ui-checks")
	if cr.Ordered {
		a.Class("ordered")
	}

	TagStart(sb, "div", a)

	if cr.List != nil {
		if cr.Ordered {
			if cr.Values != nil {
				cr.Values.Each(func(k any) bool {
					if v, ok := cr.List.Get(k); ok {
						cr.writeCheckbox(sb, k, v, true)
					}
					return true
				})
			}
			sb.WriteString("<hr>")

			cr.List.Each(func(k any, v string) bool {
				if !contains(cr.Values, k) {
					cr.writeCheckbox(sb, k, v, false)
				}
				return true
			})
		} else {
			cr.List.Each(func(k any, v string) bool {
				cr.writeCheckbox(sb, k, v, contains(cr.Values, k))
				return true
			})
		}
	}

	TagClose(sb, "div")

	return nil
}

func (cr *ChecksRenderer) writeCheckbox(sb *strings.Builder, key any, text string, checked bool) {
	TagStart(sb, "label")

	a := Attrs{
		"type":  "checkbox",
		"name":  cr.Name,
		"value": toString(key),
	}

	if checked {
		a.Set("checked", "")
	}
	if cr.Disabled {
		a.Set("disabled", "")
	}

	TagStartClose(sb, "input", a)

	sb.WriteString(html.EscapeString(text))
	TagClose(sb, "label")
}
