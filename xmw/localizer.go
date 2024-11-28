package xmw

import (
	"github.com/askasoft/pango/str"
	"github.com/askasoft/pango/xin"
)

var (
	// LocaleParamNames default parameter key names
	LocaleParamNames = []string{"_locale_"}

	// LocaleCookieNames default cookie names
	LocaleCookieNames = []string{"X_LOCALE"}

	// LocaleHeaderNames default http header names
	LocaleHeaderNames = []string{"Accept-Language"}
)

// Localizer localizer middleware
type Localizer struct {
	Locales []string

	LocaleParamNames  []string
	LocaleCookieNames []string
	LocaleHeaderNames []string
}

// NewLocalizer create a default Localizer
func NewLocalizer(locales ...string) *Localizer {
	if len(locales) == 0 {
		locales = []string{"en"}
	}

	return &Localizer{
		Locales:           locales,
		LocaleHeaderNames: LocaleHeaderNames,
		LocaleParamNames:  LocaleParamNames,
		LocaleCookieNames: LocaleCookieNames,
	}
}

// Handler returns the xin.HandlerFunc
func (ll *Localizer) Handler() xin.HandlerFunc {
	return ll.Handle
}

// Handle process xin request
func (ll *Localizer) Handle(c *xin.Context) {
	loc := ll.getLocaleFromParameter(c, ll.LocaleParamNames)

	if loc == "" {
		loc = ll.getLocaleFromCookie(c, ll.LocaleCookieNames)
	}

	if loc == "" {
		loc = ll.getLocaleFromHeader(c, ll.LocaleHeaderNames)
	}

	if loc == "" {
		if locs := ll.Locales; len(locs) > 0 {
			loc = locs[0]
		}
	}

	c.Locale = loc

	c.Next()
}

func (ll *Localizer) getLocaleFromParameter(c *xin.Context, ks []string) string {
	for _, k := range ks {
		if loc, ok := c.GetPostForm(k); ok {
			if ll.acceptable(loc) {
				return loc
			}
		}
		if loc, ok := c.GetQuery(k); ok {
			if ll.acceptable(loc) {
				return loc
			}
		}
	}
	return ""
}

func (ll *Localizer) getLocaleFromCookie(c *xin.Context, ks []string) string {
	for _, k := range ks {
		if loc, err := c.Cookie(k); err == nil {
			if ll.acceptable(loc) {
				return loc
			}
		}
	}
	return ""
}

func (ll *Localizer) getLocaleFromHeader(c *xin.Context, ks []string) string {
	for _, k := range ks {
		loc := c.GetHeader(k)
		qls := str.FieldsAny(loc, ",; ")
		for _, ql := range qls {
			if ll.acceptable(ql) {
				return ql
			}
		}
	}
	return ""
}

func (ll *Localizer) acceptable(loc string) bool {
	if loc != "" {
		for _, al := range ll.Locales {
			if str.StartsWith(loc, al) {
				return true
			}
		}
	}
	return false
}
