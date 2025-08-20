package taglib

import (
	"html"
	"strings"
)

func SelectRender(args ...any) (any, error) {
	return TagRender(&SelectRenderer{}, args...)
}

type SelectRenderer struct {
	Name     string
	List     List
	Empty    string
	Values   Values
	Disabled bool
	Multiple bool
}

func (sr *SelectRenderer) TagName() string {
	return "Select"
}

func (sr *SelectRenderer) SetList(list any) {
	sr.List = AsList(list)
}

func (sr *SelectRenderer) SetValue(value any) {
	sr.Values = AsValues([]any{value})
}

func (sr *SelectRenderer) SetValues(values any) {
	sr.Values = AsValues(values)
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

	if sr.Empty != "" {
		sr.writeOption(sb, "", sr.Empty)
	}

	if sr.List != nil {
		sr.List.Each(func(k any, v string) bool {
			sr.writeOption(sb, k, v)
			return true
		})
	}

	TagClose(sb, "select")

	return nil
}

func (sr *SelectRenderer) writeOption(sb *strings.Builder, key any, text string) {
	a := Attrs{}

	a.Set("value", toString(key))
	if contains(sr.Values, key) {
		a.Set("selected", "")
	}

	TagStart(sb, "option", a)
	sb.WriteString(html.EscapeString(text))
	TagClose(sb, "option")
}
