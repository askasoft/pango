package tags

import (
	"html"
	"strings"

	"github.com/askasoft/pango/ars"
)

func ChecksRender(args ...any) (any, error) {
	return TagRender(&ChecksRenderer{}, args...)
}

type ChecksRenderer struct {
	Name     string
	List     Iterator
	Value    []string
	Disabled bool
}

func (sr *ChecksRenderer) TagName() string {
	return "Checks"
}

func (sr *ChecksRenderer) Render(sb *strings.Builder, args ...any) error {
	a := Attrs{}

	if err := TagSetAttrs(sr, a, args); err != nil {
		return err
	}

	a.Class("ui-checks")

	TagStart(sb, "div", a)

	if sr.List != nil {
		it := sr.List
		for it.Next() {
			sr.writeCheckbox(sb, it.Key(), it.Value(), ars.ContainsString(sr.Value, it.Key()))
		}
	}

	TagClose(sb, "div")

	return nil
}

func (sr *ChecksRenderer) writeCheckbox(sb *strings.Builder, key, text string, checked bool) {
	TagStart(sb, "label")

	a := Attrs{
		"type":  "checkbox",
		"name":  sr.Name,
		"value": key,
	}

	if checked {
		a.Set("checked", "")
	}
	if sr.Disabled {
		a.Set("disabled", "")
	}

	TagStartClose(sb, "input", a)

	sb.WriteString(html.EscapeString(text))
	TagClose(sb, "label")
}
