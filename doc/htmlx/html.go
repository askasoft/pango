package htmlx

import (
	"strings"

	"github.com/askasoft/pango/iox"
	"github.com/askasoft/pango/str"
	"github.com/askasoft/pango/wcu"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func ParseHTMLFile(name string, charsets ...string) (*html.Node, error) {
	wf, _, err := wcu.DetectAndOpenFile(name, charsets...)
	if err != nil {
		return nil, err
	}
	defer wf.Close()

	return html.Parse(wf)
}

func FindElementNode(n *html.Node, tag atom.Atom) *html.Node {
	if n.Type == html.ElementNode && n.DataAtom == tag {
		return n
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if t := FindElementNode(c, tag); t != nil {
			return t
		}
	}

	return nil
}

func FindElementNodes(n *html.Node, tag atom.Atom) (ns []*html.Node) {
	if n.Type == html.ElementNode && n.DataAtom == tag {
		ns = append(ns, n)
		return
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ts := FindElementNodes(c, tag)
		ns = append(ns, ts...)
	}
	return
}

func Stringify(n *html.Node) string {
	sb := &strings.Builder{}
	cw := iox.NewCompactWriter(sb, isSpace, ' ')
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode {
			cw.WriteString(c.Data) //nolint: errcheck
		}
	}
	return str.Strip(sb.String())
}

func GetTitle(n *html.Node) string {
	if h := FindElementNode(n, atom.Head); h != nil {
		if t := FindElementNode(h, atom.Title); n != nil {
			return Stringify(t)
		}
	}
	return ""
}

func GetMetas(n *html.Node) map[string]string {
	if h := FindElementNode(n, atom.Head); h != nil {
		ns := FindElementNodes(h, atom.Meta)
		if len(ns) > 0 {
			mm := make(map[string]string, len(ns))
			for _, m := range ns {
				k, v := "", ""
				for _, a := range m.Attr {
					switch a.Key {
					case "name", "property":
						k = a.Val
					case "content":
						v = a.Val
					}
				}
				if k != "" {
					mm[k] = v
				}
			}
			return mm
		}
	}
	return nil
}
