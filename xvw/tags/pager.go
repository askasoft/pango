package tags

import (
	"fmt"
	"strings"

	"github.com/askasoft/pango/cog"
	"github.com/askasoft/pango/num"
	"github.com/askasoft/pango/str"
	"github.com/askasoft/pango/tbs"
	"github.com/askasoft/pango/xvw/args"
)

type Pager = args.Pager

func PagerRender(args ...any) (any, error) {
	return TagRender(&PageRenderer{}, args...)
}

type PageRenderer struct {
	Pager

	Locale string

	LinkSize int    // page link size
	LinkHref string // link href url

	// Style:
	//  <: <ul>
	//  p: previous page (hidden)
	//  P: previous page (disabled)
	//  n: next page (hidden)
	//  N: next page (disabled)
	//  f: first page (hidden)
	//  F: first page (disabled)
	//  l: last page (hidden)
	//  L: last page (disabled)
	//  1: #1 first page (depends on '#')
	//  #: page number links
	//  x: #x last page (depends on '#')
	//  i: pager info label
	//  s: limit size select
	//  >: </ul>
	//  I: pager info text   (float left)
	//  S: limit size select (float right)
	Style string
}

func (pr *PageRenderer) TagName() string {
	return "Pager"
}

func (pr *PageRenderer) Render(sb *strings.Builder, args ...any) error {
	if len(args) > 0 {
		if p, ok := args[0].(*Pager); ok {
			pr.Pager = *p
			args = args[1:]
		} else if p, ok := args[0].(Pager); ok {
			pr.Pager = p
			args = args[1:]
		}
	}

	a := Attrs{}

	err := TagSetAttrs(pr, a, args)
	if err != nil {
		return err
	}

	if pr.Style == "" {
		pr.Style = tbs.GetText(pr.Locale, "pager.style", "IS<FP#NL>")
	}
	if pr.LinkSize == 0 {
		pr.LinkSize = num.Atoi(tbs.GetText(pr.Locale, "pager.link-size", "5"))
	}

	a.Class("ui-pager")
	a.Data("page", num.Itoa(pr.Page))
	a.Data("limit", num.Itoa(pr.Limit))
	a.Data("total", num.Itoa(pr.Total))
	a.Data("style", pr.Style)
	a.Data("spy", "pager")

	TagStart(sb, "div", a)

	if pr.Style != "" {
		for _, r := range pr.Style {
			switch r {
			case '<':
				sb.WriteString("<ul class=\"pagination\">")
			case 'f':
				pr.writePagerLinkFirst(sb, true)
			case 'F':
				pr.writePagerLinkFirst(sb, false)
			case 'p':
				pr.writePagerLinkPrev(sb, true)
			case 'P':
				pr.writePagerLinkPrev(sb, false)
			case '#':
				if pr.Total > 0 {
					pr.writePagerLinkNums(sb)
				}
			case 'n':
				pr.writePagerLinkNext(sb, true)
			case 'N':
				pr.writePagerLinkNext(sb, false)
			case 'l':
				pr.writePagerLinkLast(sb, true)
			case 'L':
				pr.writePagerLinkLast(sb, false)
			case 'i':
				if pr.Total > 0 {
					pr.writePagerInfoLabel(sb)
				} else {
					pr.writePagerInfoEmpty(sb)
				}
			case 's':
				if pr.Total > 0 {
					err := pr.writePagerLimits(sb)
					if err != nil {
						return err
					}
				}
			case '>':
				sb.WriteString("</ul>")
			case 'I':
				if pr.Total > 0 {
					pr.writeOuterInfosLabel(sb)
				} else {
					pr.writeOuterInfosEmpty(sb)
				}
			case 'S':
				if pr.Total > 0 {
					err := pr.writeOuterLimits(sb)
					if err != nil {
						return err
					}
				}
			}
		}
	}

	TagClose(sb, "div")

	return nil
}

func (pr *PageRenderer) getLinkHref(pn int) string {
	if pr.LinkHref == "" {
		return "#"
	}

	sr := strings.NewReplacer(
		"{page}", num.Itoa(pn),
		"{limit}", num.Itoa(pr.Limit),
	)

	return sr.Replace(pr.LinkHref)
}

func (pr *PageRenderer) buildPagerTextInfo(name, defv string) string {
	return tbs.Replace(
		pr.Locale, name, defv,
		"{page}", num.Comma(pr.Page),
		"{pages}", num.Comma(pr.Pages()),
		"{begin}", num.Comma(pr.Begin()),
		"{end}", num.Comma(pr.End()),
		"{total}", num.Comma(pr.Total),
	)
}

func (pr *PageRenderer) writePagerInfo(sb *strings.Builder, info string) {
	sb.WriteString("<li class=\"info\">")
	sb.WriteString(info)
	sb.WriteString("</li>")
}

func (pr *PageRenderer) writePagerInfoEmpty(sb *strings.Builder) {
	pr.writePagerInfo(sb, tbs.GetText(pr.Locale, "pager.label.empty"))
}

func (pr *PageRenderer) writePagerInfoLabel(sb *strings.Builder) {
	info := pr.buildPagerTextInfo("pager.label.info", "{page} / {pages}")
	pr.writePagerInfo(sb, info)
}

func (pr *PageRenderer) writeOuterInfosEmpty(sb *strings.Builder) {
	pr.writeOuterInfos(sb, tbs.GetText(pr.Locale, "pager.infos.empty"))
}

func (pr *PageRenderer) writeOuterInfosLabel(sb *strings.Builder) {
	info := pr.buildPagerTextInfo("pager.infos.label", "{page} / {pages}")
	pr.writeOuterInfos(sb, info)
}

func (pr *PageRenderer) writeOuterInfos(sb *strings.Builder, info string) {
	sb.WriteString("<div class=\"infos\">")
	sb.WriteString(info)
	sb.WriteString("</div>")
}

func (pr *PageRenderer) writePagerLinkFirst(sb *strings.Builder, hidden bool) {
	sb.WriteString("<li class=\"page-item first")
	if pr.Page <= 1 {
		sb.WriteString(str.If(hidden, " hidden", " disabled"))
	}
	sb.WriteString("\"><a class=\"page-link\" href=\"")
	sb.WriteString(pr.getLinkHref(1))
	sb.WriteString("\" pageno=\"1\" title=\"")
	sb.WriteString(tbs.GetText(pr.Locale, "pager.tooltip.first"))
	sb.WriteString("\">")
	sb.WriteString(tbs.GetText(pr.Locale, "pager.label.first", "&lt;&lt;"))
	sb.WriteString("</a></li>")
}

func (pr *PageRenderer) writePagerLinkPrev(sb *strings.Builder, hidden bool) {
	p := pr.Page - 1
	if p < 1 {
		p = 1
	}

	sb.WriteString("<li class=\"page-item prev")
	if pr.Page <= 1 {
		sb.WriteString(str.If(hidden, " hidden", " disabled"))
	}
	sb.WriteString("\"><a class=\"page-link\" href=\"")
	sb.WriteString(pr.getLinkHref(p))
	sb.WriteString("\" pageno=\"")
	sb.WriteString(num.Itoa(p))
	sb.WriteString("\" title=\"")
	sb.WriteString(tbs.GetText(pr.Locale, "pager.tooltip.prev"))
	sb.WriteString("\">")
	sb.WriteString(tbs.GetText(pr.Locale, "pager.label.prev", "&lt;"))
	sb.WriteString("</a></li>")
}

func (pr *PageRenderer) writePagerLinkNext(sb *strings.Builder, hidden bool) {
	p := pr.Page + 1

	sb.WriteString("<li class=\"page-item next")
	if pr.Page >= pr.Pages() {
		sb.WriteString(str.If(hidden, " hidden", " disabled"))
	}
	sb.WriteString("\"><a class=\"page-link\" href=\"")
	sb.WriteString(pr.getLinkHref(p))
	sb.WriteString("\" pageno=\"")
	sb.WriteString(num.Itoa(p))
	sb.WriteString("\" title=\"")
	sb.WriteString(tbs.GetText(pr.Locale, "pager.tooltip.next"))
	sb.WriteString("\">")
	sb.WriteString(tbs.GetText(pr.Locale, "pager.label.next", "&gt;"))
	sb.WriteString("</a></li>")
}

func (pr *PageRenderer) writePagerLinkLast(sb *strings.Builder, hidden bool) {
	sb.WriteString("<li class=\"page-item last")
	if pr.Page >= pr.Pages() {
		sb.WriteString(str.If(hidden, " hidden", " disabled"))
	}
	sb.WriteString("\"><a class=\"page-link\" href=\"")
	sb.WriteString(pr.getLinkHref(pr.Pages()))
	sb.WriteString("\" pageno=\"")
	sb.WriteString(num.Itoa(pr.Pages()))
	sb.WriteString("\" title=\"")
	sb.WriteString(tbs.GetText(pr.Locale, "pager.tooltip.last"))
	sb.WriteString("\">")
	sb.WriteString(tbs.GetText(pr.Locale, "pager.label.last", "&gt;&gt;"))
	sb.WriteString("</a></li>")
}

func (pr *PageRenderer) writePagerLinkNums(sb *strings.Builder) {
	pe := str.ContainsByte(pr.Style, '.')
	p1 := str.ContainsByte(pr.Style, '1')
	px := str.ContainsByte(pr.Style, 'x')

	pages := pr.Pages()
	linkSize := pr.LinkSize
	linkMax := linkSize
	if p1 {
		linkMax += 2
	} else if pe {
		linkMax++
	}
	if px {
		linkMax += 2
	} else if pe {
		linkMax++
	}

	if linkMax >= pages {
		for p := 1; p <= pages; p++ {
			pr.linkp(sb, p)
		}
		return
	}

	p := 1
	if pr.Page > linkSize/2 {
		p = pr.Page - (linkSize / 2)
	}
	if p+linkSize > pages {
		p = pages - linkSize + 1
	}
	if p < 1 {
		p = 1
	}

	if p1 {
		if p > 1 {
			pr.linkp(sb, 1)
		}

		if p == 3 {
			pr.linkp(sb, 2)
		} else if p > 3 {
			pr.ellipsis(sb, true, true)
		}
	} else {
		pr.ellipsis(sb, true, pe && p > 2)
	}

	for i := 0; i < linkSize && p <= pages; i++ {
		pr.linkp(sb, p)
		p++
	}

	if px {
		if p < pages-1 {
			pr.ellipsis(sb, false, pe)
			pr.linkp(sb, pages)
		} else if p == pages-1 {
			pr.linkp(sb, p)
			pr.linkp(sb, pages)
		} else if p == pages {
			pr.linkp(sb, pages)
		}
	} else {
		pr.ellipsis(sb, false, pe && p < pages)
	}
}

func (pr *PageRenderer) ellipsis(sb *strings.Builder, left, show bool) {
	sb.WriteString("<li class=\"page-item ")
	sb.WriteString(str.If(left, "eleft", "eright"))
	if !show {
		sb.WriteString(" hidden")
	}
	sb.WriteString("\"><span>&hellip;</span></li>")
}

func (pr *PageRenderer) linkp(sb *strings.Builder, p int) {
	sb.WriteString("<li class=\"page-item page")
	if pr.Page == p {
		sb.WriteString(" active")
	}
	sb.WriteString("\"><a class=\"page-link\" href=\"")
	sb.WriteString(pr.getLinkHref(p))
	sb.WriteString("\" pageno=\"")
	sb.WriteString(num.Itoa(p))
	sb.WriteString("\">")
	sb.WriteString(num.Itoa(p))
	sb.WriteString("</a></li>")
}

func (pr *PageRenderer) writeLimitsSelect(sb *strings.Builder) error {
	sb.WriteString(tbs.GetText(pr.Locale, "pager.limits.label"))

	tlist := tbs.GetText(pr.Locale, "pager.limits.text", `%s Items`)
	slist := str.Fields(tbs.GetText(pr.Locale, "pager.limits.list", `20 50 100`))

	olist := &cog.LinkedHashMap[string, string]{}
	for _, s := range slist {
		olist.Set(s, fmt.Sprintf(tlist, s))
	}

	sr := &SelectRenderer{
		List:  olist.Iterator(),
		Value: num.Itoa(pr.Limit),
	}

	args := []any{
		"title=", tbs.GetText(pr.Locale, "pager.limits.tooltip"),
		"class=", tbs.GetText(pr.Locale, "pager.limits.class"),
	}

	return sr.Render(sb, args...)
}

func (pr *PageRenderer) writePagerLimits(sb *strings.Builder) error {
	sb.WriteString("<li class=\"limits\">")
	err := pr.writeLimitsSelect(sb)
	sb.WriteString("</li>")
	return err
}

func (pr *PageRenderer) writeOuterLimits(sb *strings.Builder) error {
	sb.WriteString("<div class=\"limits\">")
	err := pr.writeLimitsSelect(sb)
	sb.WriteString("</div>")
	return err
}
