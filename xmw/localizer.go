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
			if loc, ok = ll.acceptable(loc); ok {
				return loc
			}
		}
		if loc, ok := c.GetQuery(k); ok {
			if loc, ok = ll.acceptable(loc); ok {
				return loc
			}
		}
	}
	return ""
}

func (ll *Localizer) getLocaleFromCookie(c *xin.Context, ks []string) string {
	for _, k := range ks {
		if loc, err := c.Cookie(k); err == nil {
			if al, ok := ll.acceptable(loc); ok {
				return al
			}
		}
	}
	return ""
}

func (ll *Localizer) getLocaleFromHeader(c *xin.Context, ks []string) string {
	for _, k := range ks {
		locs := str.FieldsAny(c.GetHeader(k), ",; ")
		for _, loc := range locs {
			if al, ok := ll.acceptable(loc); ok {
				return al
			}
		}
	}
	return ""
}

func (ll *Localizer) acceptable(loc string) (string, bool) {
	if loc != "" {
		locs := ll.Locales
		for _, al := range locs {
			if str.StartsWith(loc, al) {
				return al, true
			}
		}
	}
	return loc, false
}
