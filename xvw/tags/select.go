package tags

import (
	"html"
	"strings"

	"github.com/askasoft/pango/asg"
)

func SelectRender(args ...any) (any, error) {
	return TagRender(&SelectRenderer{}, args...)
}

type SelectRenderer struct {
	Name     string
	List     Iterator
	Values   []string
	Disabled bool
	Multiple bool
}

func (sr *SelectRenderer) TagName() string {
	return "Select"
}

func (sr *SelectRenderer) SetValue(val string) {
	sr.Values = []string{val}
}

func (sr *SelectRenderer) Render(sb *strings.Builder, args ...any) error {
	a := Attrs{}

	if err := TagSetAttrs(sr, a, args); err != nil {
		return err
	}

	a.Class("ui-select")
	a.Name(sr.Name)
	if sr.Disabled {
		a.Set("disabled", "")
	}
	if sr.Multiple {
		a.Set("multiple", "")
	}

	TagStart(sb, "select", a)

	if sr.List != nil {
		it := sr.List
		for it.Next() {
			sr.writeOption(sb, it.Key(), it.Value())
		}
	}

	TagClose(sb, "select")

	return nil
}

func (sr *SelectRenderer) writeOption(sb *strings.Builder, key, text string) {
	a := Attrs{}

	a.Set("value", key)
	if asg.Contains(sr.Values, key) {
		a.Set("selected", "")
	}

	TagStart(sb, "option", a)
	sb.WriteString(html.EscapeString(text))
	TagClose(sb, "option")
}
