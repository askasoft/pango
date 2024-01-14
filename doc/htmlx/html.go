package htmlx

import (
	"io"
	"strings"

	"github.com/askasoft/pango/wcu"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func ExtractTextFromHTMLFile(name string) (string, error) {
	sb := &strings.Builder{}
	err := ExtractStringFromHTMLFile(name, sb)
	return sb.String(), err
}

func ExtractStringFromHTMLFile(name string, w io.Writer) error {
	wf, err := wcu.DetectAndOpenFile(name)
	if err != nil {
		return err
	}
	defer wf.Close()

	return ExtractStringFromHTMLReader(wf, w)
}

func ParseHTMLFile(name string) (*html.Node, error) {
	wf, err := wcu.DetectAndOpenFile(name)
	if err != nil {
		return nil, err
	}
	defer wf.Close()

	return html.Parse(wf)
}

func ExtractStringFromHTMLReader(r io.Reader, w io.Writer) error {
	doc, err := html.Parse(r)
	if err != nil {
		return err
	}
	return ExtractStringFromHTMLNode(doc, w)
}

func ExtractTextFromHTMLNode(n *html.Node) string {
	sb := strings.Builder{}
	_ = ExtractStringFromHTMLNode(n, &sb)
	return sb.String()
}

func ExtractStringFromHTMLNode(n *html.Node, w io.Writer) error {
	if n.Type == html.CommentNode {
		return nil
	}

	if n.Type == html.ElementNode {
		switch n.DataAtom {
		case atom.Script:
			return nil
		case atom.H1, atom.H2, atom.H3, atom.H4, atom.H5, atom.H6:
			fallthrough
		case atom.Div, atom.P, atom.Pre, atom.Textarea, atom.Li, atom.Br:
			if _, err := io.WriteString(w, "\n"); err != nil {
				return err
			}
		}
	}

	if n.Type == html.TextNode {
		// Keep newlines and spaces, like jQuery
		if _, err := io.WriteString(w, n.Data); err != nil {
			return err
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if err := ExtractStringFromHTMLNode(c, w); err != nil {
			return err
		}
	}

	return nil
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

func ExtractTitleFromHTMLNode(n *html.Node) string {
	head := FindElementNode(n, atom.Head)
	if head != nil {
		title := FindElementNode(n, atom.Title)
		if title != nil {
			return ExtractTextFromHTMLNode(title)
		}
	}
	return ""
}

func ExtractMetasFromHTMLNode(n *html.Node) map[string]string {
	head := FindElementNode(n, atom.Head)
	if head != nil {
		ns := FindElementNodes(n, atom.Meta)
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
