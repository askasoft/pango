package htmlx

import (
	"fmt"
	"io"
	"strings"

	"github.com/askasoft/pango/iox"
	"github.com/askasoft/pango/str"
	"github.com/askasoft/pango/wcu"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func ExtractTextFromHTMLFile(name string, charsets ...string) (string, error) {
	sb := &strings.Builder{}
	err := HTMLFileTextify(sb, name, charsets...)
	return sb.String(), err
}

func ExtractTextFromHTMLNode(node *html.Node) string {
	sb := &strings.Builder{}
	_ = Textify(node, sb)
	return sb.String()
}

func HTMLFileTextify(w io.Writer, name string, charsets ...string) error {
	wf, _, err := wcu.DetectAndOpenFile(name, charsets...)
	if err != nil {
		return err
	}
	defer wf.Close()

	return HTMLReaderTextify(w, wf)
}

func HTMLReaderTextify(w io.Writer, r io.Reader) error {
	doc, err := html.Parse(r)
	if err != nil {
		return err
	}
	return Textify(doc, w)
}

type Textifier struct {
	w  io.Writer          // underlying writer
	o  io.Writer          // current writer
	cw *iox.CompactWriter // compact writer
	ln int                // ol>li number
	lv int                // ol/ul/dir/menu level
	th bool               // thead
	tc int                // td count
}

func NewTextifier(w io.Writer) *Textifier {
	tf := &Textifier{
		w:  w,
		cw: iox.NewCompactWriter(w, isSpace, ' '),
	}
	tf.o = tf.cw

	return tf
}

func (tf *Textifier) write(s string) error {
	_, err := iox.WriteString(tf.o, s)
	return err
}

func (tf *Textifier) eol() error {
	tf.cw.Reset(true)
	_, err := tf.w.Write([]byte{'\n'})
	return err
}

func (tf *Textifier) Textify(n *html.Node) error {
	switch n.Type {
	case html.CommentNode:
		return nil
	case html.ElementNode:
		switch n.DataAtom {
		case atom.Script, atom.Style, atom.Select, atom.Object, atom.Applet, atom.Iframe, atom.Frameset, atom.Frame:
			return nil
		case atom.Br:
			return tf.lbrDeep(n)
		case atom.H1:
			return tf.hbrDeep(n, 1)
		case atom.H2:
			return tf.hbrDeep(n, 2)
		case atom.H3:
			return tf.hbrDeep(n, 3)
		case atom.H4:
			return tf.hbrDeep(n, 4)
		case atom.H5:
			return tf.hbrDeep(n, 5)
		case atom.H6:
			return tf.hbrDeep(n, 6)
		case atom.Title:
			return tf.title(n)
		case atom.Body, atom.Div, atom.P, atom.Table:
			return tf.wbrDeep(n)
		case atom.Thead:
			return tf.thead(n)
		case atom.Tr:
			return tf.tr(n)
		case atom.Th, atom.Td:
			return tf.td(n)
		case atom.Ol:
			return tf.ol(n)
		case atom.Ul, atom.Dir, atom.Menu:
			return tf.ul(n)
		case atom.Li:
			return tf.li(n)
		case atom.Dl:
			return tf.rbrDeep(n)
		case atom.Dt:
			return tf.wbrDeep(n)
		case atom.Dd:
			return tf.dd(n)
		case atom.Code, atom.Pre, atom.Textarea, atom.Xmp:
			return tf.rawDeep(n)
		default:
			return tf.deep(n)
		}
	case html.TextNode:
		return tf.write(n.Data)
	default:
		return tf.deep(n)
	}
}

func (tf *Textifier) lbrDeep(n *html.Node) error {
	if err := tf.eol(); err != nil {
		return err
	}
	return tf.deep(n)
}

func (tf *Textifier) rbrDeep(n *html.Node) error {
	if err := tf.deep(n); err != nil {
		return err
	}
	return tf.eol()
}

func (tf *Textifier) hbrDeep(n *html.Node, x int) error {
	if err := tf.eol(); err != nil {
		return err
	}
	s := str.RepeatRune('#', x) + " "
	if _, err := tf.w.Write(str.UnsafeBytes(s)); err != nil {
		return err
	}
	tf.cw.Reset(true)
	if err := tf.deep(n); err != nil {
		return err
	}
	return tf.eol()
}

func (tf *Textifier) wbrDeep(n *html.Node) error {
	if err := tf.eol(); err != nil {
		return err
	}
	if err := tf.deep(n); err != nil {
		return err
	}
	return tf.eol()
}

func (tf *Textifier) title(n *html.Node) error {
	if err := tf.eol(); err != nil {
		return err
	}
	s := Stringify(n)
	if _, err := iox.WriteString(tf.w, s); err != nil {
		return err
	}
	if err := tf.eol(); err != nil {
		return err
	}
	s = str.RepeatRune('=', len(s))
	if _, err := iox.WriteString(tf.w, s); err != nil {
		return err
	}
	return tf.eol()
}

func (tf *Textifier) thead(n *html.Node) error {
	th := tf.th
	tf.th = true
	if err := tf.deep(n); err != nil {
		return err
	}
	tf.th = th
	return nil
}

func (tf *Textifier) tr(n *html.Node) error {
	if _, err := iox.WriteString(tf.w, "| "); err != nil {
		return err
	}
	tf.cw.Reset(true)
	tc := tf.tc
	tf.tc = 0
	if err := tf.deep(n); err != nil {
		return err
	}
	if err := tf.eol(); err != nil {
		return err
	}
	if tf.th {
		if _, err := tf.w.Write([]byte{'|'}); err != nil {
			return err
		}
		for i := 0; i < tf.tc; i++ {
			if _, err := iox.WriteString(tf.w, "---|"); err != nil {
				return err
			}
		}
		if err := tf.eol(); err != nil {
			return err
		}
	}
	tf.tc = tc
	return nil
}

func (tf *Textifier) td(n *html.Node) error {
	tf.tc++
	if err := tf.deep(n); err != nil {
		return err
	}
	_, err := iox.WriteString(tf.w, " |")
	return err
}

func (tf *Textifier) ol(n *html.Node) error {
	if err := tf.eol(); err != nil {
		return err
	}
	tf.lv++
	ln := tf.ln
	tf.ln = 1
	if err := tf.deep(n); err != nil {
		return err
	}
	tf.ln = ln
	tf.lv--
	return nil
}

func (tf *Textifier) ul(n *html.Node) error {
	if err := tf.eol(); err != nil {
		return err
	}
	tf.lv++
	ln := tf.ln
	tf.ln = 0
	if err := tf.deep(n); err != nil {
		return err
	}
	tf.ln = ln
	tf.lv--
	return nil
}

func (tf *Textifier) li(n *html.Node) error {
	p := str.RepeatRune('\t', tf.lv-1)
	if tf.ln > 0 {
		if _, err := iox.WriteString(tf.w, fmt.Sprintf("%s%d. ", p, tf.ln)); err != nil {
			return err
		}
		tf.ln++
	} else {
		if _, err := iox.WriteString(tf.w, p+"- "); err != nil {
			return err
		}
	}
	tf.cw.Reset(true)
	return tf.rbrDeep(n)
}

func (tf *Textifier) dd(n *html.Node) error {
	if _, err := iox.WriteString(tf.w, ": "); err != nil {
		return err
	}
	tf.cw.Reset(true)
	return tf.rbrDeep(n)
}

func (tf *Textifier) rawDeep(n *html.Node) error {
	if err := tf.eol(); err != nil {
		return err
	}
	o := tf.o
	tf.o = tf.w
	if err := tf.deep(n); err != nil {
		return err
	}
	tf.o = o
	return tf.eol()
}

func (tf *Textifier) deep(n *html.Node) error {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if err := tf.Textify(c); err != nil {
			return err
		}
	}
	return nil
}

func Textify(n *html.Node, w io.Writer) error {
	tf := NewTextifier(w)
	return tf.Textify(n)
}

func isSpace(r rune) bool {
	return r <= 0x20
}
