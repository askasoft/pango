package htmlx

import (
	"io"
	"strings"

	"github.com/askasoft/pango/iox"
	"github.com/askasoft/pango/num"
	"github.com/askasoft/pango/str"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func HTMLFileTextifyString(name string, detect int, charsets ...string) (string, error) {
	sb := &strings.Builder{}
	err := HTMLFileTextify(sb, name, detect, charsets...)
	return sb.String(), err
}

func HTMLFileTextify(w io.Writer, name string, detect int, charsets ...string) error {
	doc, err := ParseHTMLFile(name, detect, charsets...)
	if err != nil {
		return err
	}

	return HTMLNodeTextify(w, doc)
}

func HTMLTextifyString(html string) (string, error) {
	sb := &strings.Builder{}
	err := HTMLTextify(sb, html)
	return sb.String(), err
}

func HTMLTextify(w io.Writer, html string) error {
	r := strings.NewReader(html)
	return HTMLReaderTextify(w, r)
}

func HTMLReaderTextifyString(r io.Reader) (string, error) {
	sb := &strings.Builder{}
	err := HTMLReaderTextify(sb, r)
	return sb.String(), err
}

func HTMLReaderTextify(w io.Writer, r io.Reader) error {
	doc, err := html.Parse(r)
	if err != nil {
		return err
	}
	return HTMLNodeTextify(w, doc)
}

func HTMLNodeTextifyString(node *html.Node) string {
	sb := &strings.Builder{}
	_ = HTMLNodeTextify(sb, node)
	return sb.String()
}

func HTMLNodeTextify(w io.Writer, n *html.Node) error {
	tf := NewTextifier(w)
	return tf.Textify(n)
}

func isSpace(r rune) bool {
	return r <= 0x20
}

type textWriter interface {
	io.Writer
	io.StringWriter
}

type Textifier struct {
	Custom func(tf *Textifier, n *html.Node) (bool, error)
	Escape func(s string) string

	tw textWriter         // text writer
	pw *iox.ProxyWriter   // proxy writer
	cw *iox.CompactWriter // compact writer
	lv int                // ol/ul/dir/menu level
	lt string             // ol list type
	ln int                // ol>li number
	th bool               // thead
	td int                // td count
}

func NewTextifier(w io.Writer) *Textifier {
	pw := &iox.ProxyWriter{W: w}
	cw := iox.NewCompactWriter(pw, isSpace, ' ')
	tf := &Textifier{
		Escape: noescape,
		tw:     cw,
		pw:     pw,
		cw:     cw,
	}
	return tf
}

func noescape(s string) string {
	return s
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
		case atom.Noscript, atom.Script, atom.Style, atom.Select, atom.Object, atom.Applet:
			return nil
		case atom.Iframe, atom.Frameset, atom.Frame, atom.Rb, atom.Rp, atom.Rt, atom.Rtc:
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
		case atom.Blockquote:
			return tf.Blockquote(n)
		case atom.A:
			return tf.A(n)
		case atom.Img:
			return tf.Img(n)
		case atom.B, atom.Strong:
			return tf.Bold(n)
		case atom.Em, atom.I:
			return tf.Italic(n)
		case atom.S, atom.Strike:
			return tf.Strike(n)
		case atom.Div, atom.P, atom.Section:
			return tf.RbrDeep(n)
		case atom.Code, atom.Pre, atom.Textarea, atom.Xmp:
			return tf.RawDeep(n)
		default:
			return tf.Deep(n)
		}
	case html.TextNode:
		return tf.Text(n.Data)
	default:
		return tf.Deep(n)
	}
}

func (tf *Textifier) Text(ss ...string) error {
	for _, s := range ss {
		if _, err := tf.tw.WriteString(tf.Escape(s)); err != nil {
			return err
		}
	}
	return nil
}

func (tf *Textifier) Eol() error {
	tf.cw.Reset(true)
	_, err := tf.pw.WriteString("\n")
	return err
}

func (tf *Textifier) Bold(n *html.Node) error {
	return tf.WrapDeep(n, "**")
}

func (tf *Textifier) Italic(n *html.Node) error {
	return tf.WrapDeep(n, "*")
}

func (tf *Textifier) Strike(n *html.Node) error {
	return tf.WrapDeep(n, "~")
}

func (tf *Textifier) WrapDeep(n *html.Node, w string) error {
	if _, err := tf.pw.WriteString(w); err != nil {
		return err
	}
	if err := tf.Deep(n); err != nil {
		return err
	}
	_, err := tf.pw.WriteString(w)
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
	if _, err := iox.RepeatWriteString(tf.pw, "#", x); err != nil {
		return err
	}
	if _, err := tf.pw.WriteString(" "); err != nil {
		return err
	}
	tf.cw.Reset(true)
	if err := tf.Deep(n); err != nil {
		return err
	}
	return tf.Eol()
}

func (tf *Textifier) A(n *html.Node) error {
	href := str.Strip(GetNodeAttrValue(n, "href"))
	if href == "" {
		return tf.Deep(n)
	}

	var sa strings.Builder
	ht := NewTextifier(&sa)
	if err := ht.Deep(n); err != nil {
		return err
	}
	text := str.Strip(sa.String())

	if href == text {
		return tf.Text(" ", href, " ")
	}

	if _, err := tf.pw.WriteString(" ["); err != nil {
		return err
	}
	if _, err := tf.pw.WriteString(text); err != nil {
		return err
	}
	if _, err := tf.pw.WriteString(" ]("); err != nil {
		return err
	}
	if err := tf.Text(href); err != nil {
		return err
	}
	if _, err := tf.pw.WriteString(") "); err != nil {
		return err
	}
	return nil
}

func (tf *Textifier) Img(n *html.Node) error {
	src := str.Strip(GetNodeAttrValue(n, "src"))
	if src != "" {
		alt := str.Strip(GetNodeAttrValue(n, "alt"))

		if _, err := tf.pw.WriteString(" !["); err != nil {
			return err
		}
		if err := tf.Text(alt); err != nil {
			return err
		}
		if _, err := tf.pw.WriteString(" ]("); err != nil {
			return err
		}
		if err := tf.Text(src); err != nil {
			return err
		}
		if _, err := tf.pw.WriteString(") "); err != nil {
			return err
		}
	}

	return tf.Deep(n)
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
	if err := tf.Text(s); err != nil {
		return err
	}
	if err := tf.Eol(); err != nil {
		return err
	}
	if _, err := iox.RepeatWriteString(tf.pw, "=", len(s)); err != nil {
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
	if _, err := tf.pw.WriteString("| "); err != nil {
		return err
	}
	tf.cw.Reset(true)
	td := tf.td
	tf.td = 0
	if err := tf.Deep(n); err != nil {
		return err
	}
	if err := tf.Eol(); err != nil {
		return err
	}
	if tf.th {
		if _, err := tf.pw.WriteString("|"); err != nil {
			return err
		}
		for range tf.td {
			if _, err := tf.pw.WriteString("---|"); err != nil {
				return err
			}
		}
		if err := tf.Eol(); err != nil {
			return err
		}
	}
	tf.td = td
	return nil
}

func (tf *Textifier) Td(n *html.Node) error {
	tf.td++
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
	lt, ln := tf.lt, tf.ln
	tf.lt = GetNodeAttrValue(n, "type")
	tf.ln = num.Atoi(GetNodeAttrValue(n, "start"), 1)
	if err := tf.Deep(n); err != nil {
		return err
	}
	tf.lt, tf.ln = lt, ln
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
	if _, err := iox.RepeatWriteString(tf.pw, "\t", tf.lv-1); err != nil {
		return err
	}
	if tf.ln > 0 {
		var s string
		switch tf.lt {
		case "a":
			s = str.ToLower(num.IntToAlpha(tf.ln))
		case "A":
			s = num.IntToAlpha(tf.ln)
		case "i":
			s = str.ToLower(num.IntToRoman(tf.ln))
		case "I":
			s = num.IntToRoman(tf.ln)
		default:
			s = num.Itoa(tf.ln)
		}
		if _, err := tf.pw.WriteString(s); err != nil {
			return err
		}
		if _, err := tf.pw.WriteString("."); err != nil {
			return err
		}
		tf.ln++
	} else {
		if _, err := tf.pw.WriteString("-"); err != nil {
			return err
		}
	}
	if _, err := tf.pw.WriteString(" "); err != nil {
		return err
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

	tw := tf.tw
	tf.tw = tf.pw
	if err := tf.Deep(n); err != nil {
		return err
	}
	tf.tw = tw
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
