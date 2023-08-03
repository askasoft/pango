package xvw

import (
	"github.com/askasoft/pango/tbs"
	"github.com/askasoft/pango/tpl"
	"github.com/askasoft/pango/xvw/tags"
)

// Functions default utility functions for template
func Functions() tpl.FuncMap {
	return tpl.FuncMap{
		// Text Bundle
		"T":  tbs.GetText,
		"TF": tbs.Format,
		"TR": tbs.Replace,

		// Tags
		"Pager":  tags.PagerRender,
		"Select": tags.SelectRender,
	}
}
