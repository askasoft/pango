package xmw

import (
	"github.com/askasoft/pango/str"
	"github.com/askasoft/pango/xin"
)

const (
	// LocaleParamName default parameter key name
	LocaleParamName = "__locale"

	// LocaleHeaderName default http header name
	LocaleHeaderName = "X-Accept-Language"

	// LocaleCookieName default cookie name
	LocaleCookieName = "WW_LOCALE"
)

// Localizer localizer middleware
type Localizer struct {
	Locales []string

	LocaleParamName    string
	LocaleHeaderName   string
	LocaleCookieName   string
	FromAcceptLanguage bool
}

// NewLocalizer create a default Localizer
func NewLocalizer(locales ...string) *Localizer {
	if len(locales) == 0 {
		locales = []string{"en"}
	}

	return &Localizer{
		Locales:            locales,
		LocaleHeaderName:   LocaleHeaderName,
		LocaleParamName:    LocaleParamName,
		LocaleCookieName:   LocaleCookieName,
		FromAcceptLanguage: true,
	}
}

// Handler returns the xin.HandlerFunc
func (ll *Localizer) Handler() xin.HandlerFunc {
	return func(c *xin.Context) {
		ll.handle(c)
	}
}

// handle process xin request
func (ll *Localizer) handle(c *xin.Context) {
	loc := ""

	if ll.LocaleParamName != "" {
		loc = ll.getLocaleFromParameter(c, ll.LocaleParamName)
	}

	if loc == "" && ll.LocaleHeaderName != "" {
		loc = ll.getLocaleFromHeader(c, ll.LocaleHeaderName)
	}

	if loc == "" && ll.LocaleCookieName != "" {
		loc = ll.getLocaleFromCookie(c, ll.LocaleCookieName)
	}

	if loc == "" && ll.FromAcceptLanguage {
		loc = ll.getLocaleFromHeader(c, "Accept-Language")
	}

	if loc == "" {
		loc = ll.Locales[0]
	}

	c.Locale = loc

	c.Next()
}

func (ll *Localizer) getLocaleFromHeader(c *xin.Context, k string) string {
	loc := c.GetHeader(k)
	qls := str.FieldsAny(loc, ",; ")
	for _, ql := range qls {
		if ll.acceptable(ql) {
			return ql
		}
	}
	return ""
}

func (ll *Localizer) getLocaleFromParameter(c *xin.Context, k string) string {
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
	return ""
}

func (ll *Localizer) getLocaleFromCookie(c *xin.Context, k string) string {
	if loc, err := c.Cookie(k); err == nil {
		if ll.acceptable(loc) {
			return loc
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
