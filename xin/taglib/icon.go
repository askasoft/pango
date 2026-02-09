package taglib

import (
	"strings"

	"github.com/askasoft/pango/tbs"
)

func IconRender(args ...any) (any, error) {
	return TagRender(&IconRenderer{}, args...)
}

type IconRenderer struct {
	Locale string
}

func (ir *IconRenderer) TagName() string {
	return "I"
}

func (ir *IconRenderer) Render(sb *strings.Builder, args ...any) error {
	var icon string

	if len(args) > 0 {
		icon = args[0].(string)
		args = args[1:]
	}

	a := Attrs{}

	if err := TagSetAttrs(ir, a, args); err != nil {
		return err
	}

	if icon != "" {
		a.Class(tbs.GetText(ir.Locale, "icon."+icon))
	}

	TagStart(sb, "i", a)
	TagClose(sb, "i")

	return nil
}
