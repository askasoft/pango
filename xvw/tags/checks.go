package tags

import (
	"html"
	"strings"

	"github.com/askasoft/pango/asg"
)

func ChecksRender(args ...any) (any, error) {
	return TagRender(&ChecksRenderer{}, args...)
}

type ChecksRenderer struct {
	Name     string
	List     List[string, string]
	Values   []string
	Disabled bool
	Ordered  bool
}

func (cr *ChecksRenderer) TagName() string {
	return "Checks"
}

func (cr *ChecksRenderer) SetValue(val string) {
	cr.Values = []string{val}
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
			for _, k := range cr.Values {
				if v, ok := cr.List.Get(k); ok {
					cr.writeCheckbox(sb, k, v, true)
				}
			}
			sb.WriteString("<hr>")

			for it := cr.List.Iterator(); it.Next(); {
				if !asg.Contains(cr.Values, it.Key()) {
					cr.writeCheckbox(sb, it.Key(), it.Value(), false)
				}
			}
		} else {
			for it := cr.List.Iterator(); it.Next(); {
				cr.writeCheckbox(sb, it.Key(), it.Value(), asg.Contains(cr.Values, it.Key()))
			}
		}
	}

	TagClose(sb, "div")

	return nil
}

func (cr *ChecksRenderer) writeCheckbox(sb *strings.Builder, key, text string, checked bool) {
	TagStart(sb, "label")

	a := Attrs{
		"type":  "checkbox",
		"name":  cr.Name,
		"value": key,
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
