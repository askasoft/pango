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
	//  .: ellipsis
	//  >: </ul>
	//  i: pager info text
	//  s: limit size select
	Style string
}

func (pr *PageRenderer) Name() string {
	return "Pager"
}

func (pr *PageRenderer) Render(sb *strings.Builder, args ...any) error {
	attrs := Attrs{}

	if len(args) > 0 {
		if p, ok := args[0].(*Pager); ok {
			pr.Pager = *p
			args = args[1:]
		} else if p, ok := args[0].(Pager); ok {
			pr.Pager = p
			args = args[1:]
		}
	}

	err := TagSetAttrs(pr, attrs, args)
	if err != nil {
		return err
	}

	if pr.Style == "" {
		pr.Style = tbs.GetText(pr.Locale, "pager.style", "is<FP#NL>")
	}
	if pr.LinkSize == 0 {
		pr.LinkSize = num.Atoi(tbs.GetText(pr.Locale, "pager.link-size", "5"))
	}

	attrs.Class("ui-pager clearfix")
	attrs.Data("page", num.Itoa(pr.Page))
	attrs.Data("limit", num.Itoa(pr.Limit))
	attrs.Data("total", num.Itoa(pr.Total))
	attrs.Data("style", pr.Style)
	attrs.Data("spy", "pager")

	TagStart(sb, "div", attrs)

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
			case '>':
				sb.WriteString("</ul>")
			case 'i':
				if pr.Total > 0 {
					pr.writePagerTextInfo(sb)
				} else {
					pr.writePagerEmptyInfo(sb)
				}
			case 's':
				if pr.Total > 0 {
					err := pr.writePagerLimit(sb)
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

func (pr *PageRenderer) writePagerEmptyInfo(sb *strings.Builder) {
	pr.writePagerInfo(sb, tbs.GetText(pr.Locale, "pager.label-empty"))
}

func (pr *PageRenderer) writePagerTextInfo(sb *strings.Builder) {
	info := tbs.Replace(
		pr.Locale, "pager.infos.label", "{page} / {pages}",
		"{page}", num.Comma(pr.Page),
		"{pages}", num.Comma(pr.Pages()),
		"{begin}", num.Comma(pr.Begin()),
		"{end}", num.Comma(pr.End()),
		"{total}", num.Comma(pr.Total),
	)
	pr.writePagerInfo(sb, info)
}

func (pr *PageRenderer) writePagerInfo(sb *strings.Builder, info string) {
	sb.WriteString("<div class=\"infos\">")
	sb.WriteString(info)
	sb.WriteString("</div>")
}

func (pr *PageRenderer) writePagerLinkFirst(sb *strings.Builder, hidden bool) {
	sb.WriteString("<li class=\"first")
	if pr.Page <= 1 {
		sb.WriteString(str.If(hidden, " hidden", " disabled"))
	}
	sb.WriteString("\"><a href=\"")
	sb.WriteString(pr.getLinkHref(1))
	sb.WriteByte('"')
	sb.WriteString(" pageno=\"1")
	sb.WriteString("\" title=\"")
	sb.WriteString(tbs.GetText(pr.Locale, "pager.tooltip.first"))
	sb.WriteByte('"')
	sb.WriteByte('>')
	sb.WriteString(tbs.GetText(pr.Locale, "pager.label.first", "&lt;&lt;"))
	sb.WriteString("</a></li>")
}

func (pr *PageRenderer) writePagerLinkPrev(sb *strings.Builder, hidden bool) {
	p := pr.Page - 1
	if p < 1 {
		p = 1
	}

	sb.WriteString("<li class=\"prev")
	if pr.Page <= 1 {
		sb.WriteString(str.If(hidden, " hidden", " disabled"))
	}
	sb.WriteString("\"><a href=\"")
	sb.WriteString(pr.getLinkHref(p))
	sb.WriteString("\" pageno=\"")
	sb.WriteString(num.Itoa(p))
	sb.WriteString("\" title=\"")
	sb.WriteString(tbs.GetText(pr.Locale, "pager.tooltip.prev"))
	sb.WriteByte('"')
	sb.WriteByte('>')
	sb.WriteString(tbs.GetText(pr.Locale, "pager.label.prev", "&lt;"))
	sb.WriteString("</a></li>")
}

func (pr *PageRenderer) writePagerLinkNext(sb *strings.Builder, hidden bool) {
	p := pr.Page + 1

	sb.WriteString("<li class=\"next")
	if pr.Page >= pr.Pages() {
		sb.WriteString(str.If(hidden, " hidden", " disabled"))
	}
	sb.WriteString("\"><a href=\"")
	sb.WriteString(pr.getLinkHref(p))
	sb.WriteString("\" pageno=\"")
	sb.WriteString(num.Itoa(p))
	sb.WriteString("\" title=\"")
	sb.WriteString(tbs.GetText(pr.Locale, "pager.tooltip.next"))
	sb.WriteByte('"')
	sb.WriteByte('>')
	sb.WriteString(tbs.GetText(pr.Locale, "pager.label.next", "&gt;"))
	sb.WriteString("</a></li>")
}

func (pr *PageRenderer) writePagerLinkLast(sb *strings.Builder, hidden bool) {
	sb.WriteString("<li class=\"last")
	if pr.Page >= pr.Pages() {
		sb.WriteString(str.If(hidden, " hidden", " disabled"))
	}
	sb.WriteString("\"><a href=\"")
	sb.WriteString(pr.getLinkHref(pr.Pages()))
	sb.WriteString("\" pageno=\"")
	sb.WriteString(num.Itoa(pr.Pages()))
	sb.WriteString("\" title=\"")
	sb.WriteString(tbs.GetText(pr.Locale, "pager.tooltip.last"))
	sb.WriteByte('"')
	sb.WriteByte('>')
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
	sb.WriteString("<li class=\"")
	sb.WriteString(str.If(left, "eleft", "eright"))
	if !show {
		sb.WriteString(" hidden")
	}
	sb.WriteString("\"><span>&hellip;</span></li>")
}

func (pr *PageRenderer) linkp(sb *strings.Builder, p int) {
	sb.WriteString("<li class=\"page")
	if pr.Page == p {
		sb.WriteString(" active")
	}
	sb.WriteString("\"><a href=\"")
	sb.WriteString(pr.getLinkHref(p))
	sb.WriteString("\" pageno=\"")
	sb.WriteString(num.Itoa(p))
	sb.WriteString("\">")
	sb.WriteString(num.Itoa(p))
	sb.WriteString("</a></li>")
}

func (pr *PageRenderer) writePagerLimit(sb *strings.Builder) error {
	sb.WriteString("<div class=\"limits\">")

	sb.WriteString(tbs.GetText(pr.Locale, "pager.limits.label"))

	tlist := tbs.GetText(pr.Locale, "pager.limits.text", `%s Items`)
	slist := str.Fields(tbs.GetText(pr.Locale, "pager.limits.list", `20 50 100`))

	olist := &cog.LinkedHashMap[string, string]{}
	for _, s := range slist {
		olist.Set(s, fmt.Sprintf(tlist, s))
	}

	sr := &SelectRenderer{
		Value: num.Itoa(pr.Limit),
		List:  olist.Iterator(),
	}

	err := sr.Render(sb,
		"class", "select form-control",
		"title", tbs.GetText(pr.Locale, "pager.tooltip-limits"),
	)
	if err != nil {
		return err
	}

	sb.WriteString("</div>")
	return nil
}