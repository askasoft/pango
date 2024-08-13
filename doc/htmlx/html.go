package htmlx

import (
	"strings"

	"github.com/askasoft/pango/iox"
	"github.com/askasoft/pango/str"
	"github.com/askasoft/pango/wcu"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func ParseHTMLFile(name string, detect int, charsets ...string) (*html.Node, error) {
	wf, _, err := wcu.DetectAndOpenFile(name, detect, charsets...)
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

func FindElementNodes(root *html.Node, tag atom.Atom) (ns []*html.Node) {
	_ = IterElementNodes(root, func(n *html.Node) error {
		if n.DataAtom == tag {
			ns = append(ns, n)
		}
		return nil
	})
	return
}

func IterElementNodes(n *html.Node, iter func(n *html.Node) error) error {
	if n.Type == html.ElementNode {
		if err := iter(n); err != nil {
			return err
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if err := IterElementNodes(c, iter); err != nil {
			return err
		}
	}

	return nil
}

func IterNodes(n *html.Node, iter func(n *html.Node) error) error {
	if err := iter(n); err != nil {
		return err
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if err := IterNodes(c, iter); err != nil {
			return err
		}
	}

	return nil
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

func HTML(n *html.Node) (string, error) {
	sb := &strings.Builder{}
	err := html.Render(sb, n)
	return sb.String(), err
}

func FindNodeAttr(n *html.Node, k string) *html.Attribute {
	for i := 0; i < len(n.Attr); i++ {
		a := &n.Attr[i]
		if a.Key == k {
			return a
		}
	}
	return nil
}

func FindAndGetHtmlLang(doc *html.Node) string {
	if h := FindElementNode(doc, atom.Html); h != nil {
		if a := FindNodeAttr(h, "lang"); a != nil {
			return str.ToLower(a.Val)
		}
	}
	return ""
}

func FindAndGetHeadTitle(doc *html.Node) string {
	if h := FindElementNode(doc, atom.Head); h != nil {
		if t := FindElementNode(h, atom.Title); t != nil {
			return Stringify(t)
		}
	}
	return ""
}

func FindAndGetHeadMetas(doc *html.Node) map[string]string {
	if h := FindElementNode(doc, atom.Head); h != nil {
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
