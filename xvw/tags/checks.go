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
	List     Iterator
	Values   []string
	Disabled bool
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

	TagStart(sb, "div", a)

	if cr.List != nil {
		it := cr.List
		for it.Next() {
			cr.writeCheckbox(sb, it.Key(), it.Value(), asg.Contains(cr.Values, it.Key()))
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
