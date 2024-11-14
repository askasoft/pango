package htmlx

import (
	"fmt"
	"io"
	"strings"

	"github.com/askasoft/pango/iox"
	"github.com/askasoft/pango/str"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func ExtractTextFromHTMLFile(name string, detect int, charsets ...string) (string, error) {
	sb := &strings.Builder{}
	err := HTMLFileTextify(sb, name, detect, charsets...)
	return sb.String(), err
}

func ExtractTextFromHTMLString(html string) (string, error) {
	sb := &strings.Builder{}
	err := HTMLStringTextify(sb, html)
	return sb.String(), err
}

func ExtractTextFromHTMLNode(node *html.Node) string {
	sb := &strings.Builder{}
	_ = Textify(sb, node)
	return sb.String()
}

func HTMLFileTextify(w io.Writer, name string, detect int, charsets ...string) error {
	doc, err := ParseHTMLFile(name, detect, charsets...)
	if err != nil {
		return err
	}

	return Textify(w, doc)
}

func HTMLStringTextify(w io.Writer, html string) error {
	r := strings.NewReader(html)

	return HTMLReaderTextify(w, r)
}

func HTMLReaderTextify(w io.Writer, r io.Reader) error {
	doc, err := html.Parse(r)
	if err != nil {
		return err
	}
	return Textify(w, doc)
}

func Textify(w io.Writer, n *html.Node) error {
	tf := NewTextifier(w)
	return tf.Textify(n)
}

func isSpace(r rune) bool {
	return r <= 0x20
}

type Textifier struct {
	Custom func(tf *Textifier, n *html.Node) (bool, error)

	ow io.Writer          // current writer
	pw *iox.ProxyWriter   // proxy writer
	cw *iox.CompactWriter // compact writer
	ln int                // ol>li number
	lv int                // ol/ul/dir/menu level
	th bool               // thead
	tc int                // td count
}

func NewTextifier(w io.Writer) *Textifier {
	pw := &iox.ProxyWriter{W: w}
	cw := iox.NewCompactWriter(pw, isSpace, ' ')
	tf := &Textifier{
		ow: cw,
		pw: pw,
		cw: cw,
	}
	return tf
}

func (tf *Textifier) Textify(n *html.Node) error {
	if cf := tf.Custom; cf != nil {
		if ok, err := cf(tf, n); ok || err != nil {
			return err
		}
	}
	switch n.Type {
	case html.CommentNode:
		return nil
	case html.ElementNode:
		switch n.DataAtom {
		case atom.Noscript, atom.Script, atom.Style, atom.Select, atom.Object, atom.Applet, atom.Iframe, atom.Frameset, atom.Frame:
			return nil
		case atom.Br:
			return tf.LbrDeep(n)
		case atom.H1:
			return tf.HbrDeep(n, 1)
		case atom.H2:
			return tf.HbrDeep(n, 2)
		case atom.H3:
			return tf.HbrDeep(n, 3)
		case atom.H4:
			return tf.HbrDeep(n, 4)
		case atom.H5:
			return tf.HbrDeep(n, 5)
		case atom.H6:
			return tf.HbrDeep(n, 6)
		case atom.Title:
			return tf.Title(n)
		case atom.Body, atom.Table:
			return tf.WbrDeep(n)
		case atom.Thead:
			return tf.Thead(n)
		case atom.Tr:
			return tf.Tr(n)
		case atom.Th, atom.Td:
			return tf.Td(n)
		case atom.Ol:
			return tf.Ol(n)
		case atom.Ul, atom.Dir, atom.Menu:
			return tf.Ul(n)
		case atom.Li:
			return tf.Li(n)
		case atom.Dl:
			return tf.RbrDeep(n)
		case atom.Dt:
			return tf.WbrDeep(n)
		case atom.Dd:
			return tf.Dd(n)
		case atom.Div, atom.P:
			return tf.RbrDeep(n)
		case atom.Blockquote:
			return tf.Blockquote(n)
		case atom.Code, atom.Pre, atom.Textarea, atom.Xmp:
			return tf.RawDeep(n)
		default:
			return tf.Deep(n)
		}
	case html.TextNode:
		return tf.Write(n.Data)
	default:
		return tf.Deep(n)
	}
}

func (tf *Textifier) Write(s string) error {
	_, err := iox.WriteString(tf.ow, s)
	return err
}

func (tf *Textifier) Eol() error {
	tf.cw.Reset(true)
	_, err := tf.pw.Write([]byte{'\n'})
	return err
}

func (tf *Textifier) LbrDeep(n *html.Node) error {
	if err := tf.Eol(); err != nil {
		return err
	}
	return tf.Deep(n)
}

func (tf *Textifier) RbrDeep(n *html.Node) error {
	if err := tf.Deep(n); err != nil {
		return err
	}
	return tf.Eol()
}

func (tf *Textifier) HbrDeep(n *html.Node, x int) error {
	if err := tf.Eol(); err != nil {
		return err
	}
	s := str.RepeatRune('#', x) + " "
	if _, err := tf.pw.Write(str.UnsafeBytes(s)); err != nil {
		return err
	}
	tf.cw.Reset(true)
	if err := tf.Deep(n); err != nil {
		return err
	}
	return tf.Eol()
}

func (tf *Textifier) WbrDeep(n *html.Node) error {
	if err := tf.Eol(); err != nil {
		return err
	}
	if err := tf.Deep(n); err != nil {
		return err
	}
	return tf.Eol()
}

func (tf *Textifier) Title(n *html.Node) error {
	s := Stringify(n)
	if _, err := iox.WriteString(tf.pw, s); err != nil {
		return err
	}
	if err := tf.Eol(); err != nil {
		return err
	}
	s = str.RepeatRune('=', len(s))
	if _, err := iox.WriteString(tf.pw, s); err != nil {
		return err
	}
	return tf.Eol()
}

func (tf *Textifier) Thead(n *html.Node) error {
	th := tf.th
	tf.th = true
	if err := tf.Deep(n); err != nil {
		return err
	}
	tf.th = th
	return nil
}

func (tf *Textifier) Tr(n *html.Node) error {
	if _, err := iox.WriteString(tf.pw, "| "); err != nil {
		return err
	}
	tf.cw.Reset(true)
	tc := tf.tc
	tf.tc = 0
	if err := tf.Deep(n); err != nil {
		return err
	}
	if err := tf.Eol(); err != nil {
		return err
	}
	if tf.th {
		if _, err := tf.pw.Write([]byte{'|'}); err != nil {
			return err
		}
		for i := 0; i < tf.tc; i++ {
			if _, err := tf.pw.WriteString("---|"); err != nil {
				return err
			}
		}
		if err := tf.Eol(); err != nil {
			return err
		}
	}
	tf.tc = tc
	return nil
}

func (tf *Textifier) Td(n *html.Node) error {
	tf.tc++
	if err := tf.Deep(n); err != nil {
		return err
	}
	_, err := tf.pw.WriteString(" |")
	return err
}

func (tf *Textifier) Ol(n *html.Node) error {
	if err := tf.Eol(); err != nil {
		return err
	}
	tf.lv++
	ln := tf.ln
	tf.ln = 1
	if err := tf.Deep(n); err != nil {
		return err
	}
	tf.ln = ln
	tf.lv--
	return nil
}

func (tf *Textifier) Ul(n *html.Node) error {
	if err := tf.Eol(); err != nil {
		return err
	}
	tf.lv++
	ln := tf.ln
	tf.ln = 0
	if err := tf.Deep(n); err != nil {
		return err
	}
	tf.ln = ln
	tf.lv--
	return nil
}

func (tf *Textifier) Li(n *html.Node) error {
	p := str.RepeatRune('\t', tf.lv-1)
	if tf.ln > 0 {
		if _, err := tf.pw.WriteString(fmt.Sprintf("%s%d. ", p, tf.ln)); err != nil {
			return err
		}
		tf.ln++
	} else {
		if _, err := tf.pw.WriteString(p + "- "); err != nil {
			return err
		}
	}
	tf.cw.Reset(true)
	return tf.RbrDeep(n)
}

func (tf *Textifier) Dd(n *html.Node) error {
	if _, err := tf.pw.WriteString(": "); err != nil {
		return err
	}
	tf.cw.Reset(true)
	return tf.RbrDeep(n)
}

func (tf *Textifier) Blockquote(n *html.Node) error {
	if err := tf.Eol(); err != nil {
		return err
	}

	w := tf.pw.W
	tf.pw.W = iox.LinePrefixWriter(w, "> ")
	if err := tf.Deep(n); err != nil {
		return err
	}
	tf.pw.W = w

	return tf.Eol()
}

func (tf *Textifier) RawDeep(n *html.Node) error {
	if err := tf.Eol(); err != nil {
		return err
	}

	ow := tf.ow
	tf.ow = tf.pw
	if err := tf.Deep(n); err != nil {
		return err
	}
	tf.ow = ow
	return tf.Eol()
}

func (tf *Textifier) Deep(n *html.Node) error {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if err := tf.Textify(c); err != nil {
			return err
		}
	}
	return nil
}
