package taglib

import (
	"github.com/askasoft/pango/tbs"
	"github.com/askasoft/pango/tpl"
)

// Functions default utility functions for template
func Functions() tpl.FuncMap {
	return tpl.FuncMap{
		// Text Bundle
		"T":  tbs.GetText,
		"TF": tbs.Format,
		"TR": tbs.Replace,

		// Tags
		"I":      Icon,
		"Checks": ChecksRender,
		"Radios": RadiosRender,
		"Select": SelectRender,
		"Pager":  PagerRender,
	}
}
