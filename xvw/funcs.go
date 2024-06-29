package xvw

import (
	"html/template"

	"github.com/askasoft/pango/str"
	"github.com/askasoft/pango/tbs"
	"github.com/askasoft/pango/tpl"
	"github.com/askasoft/pango/xvw/tags"
)

// Functions default utility functions for template
func Functions() tpl.FuncMap {
	return tpl.FuncMap{
		// Icon
		"I": icon,

		// Text Bundle
		"T":  tbs.GetText,
		"TF": tbs.Format,
		"TR": tbs.Replace,

		// Tags
		"Checks": tags.ChecksRender,
		"Radios": tags.RadiosRender,
		"Select": tags.SelectRender,
		"Pager":  tags.PagerRender,
	}
}

func icon(icon string, clss ...string) template.HTML {
	clss = append(clss, tbs.GetText("", "icon."+icon))
	return template.HTML(`<i class="` + str.Join(clss, " ") + `"></i>`) //nolint: gosec
}
