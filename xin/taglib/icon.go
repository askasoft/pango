package taglib

import (
	"html/template"

	"github.com/askasoft/pango/str"
	"github.com/askasoft/pango/tbs"
)

func Icon(icon string, clss ...string) template.HTML {
	clss = append(clss, tbs.GetText("", "icon."+icon))
	return template.HTML(`<i class="` + str.Join(clss, " ") + `"></i>`) //nolint: gosec
}
