package xmlx

import (
	"bytes"
	"encoding/xml"
	"errors"
	"io"
	"path"
	"strings"

	"github.com/askasoft/pango/str"
)

var (
	// ErrInvalidDocument invalid document err
	ErrInvalidDocument = errors.New("xmlx: invalid document")

	// ErrInvalidRoot data at the root level is invalid err
	ErrInvalidRoot = errors.New("xmlx: data at the root level is invalid")

	mapdec = NewMapDecoder("@", "#text")
)

// MapDecoder a xml decoder for map[string]any
type MapDecoder struct {
	a string
	t string
}

// NewMapDecoder create new decoder instance with custom attribute prefix and text key
func NewMapDecoder(attrPrefix, textKey string) *MapDecoder {
	return &MapDecoder{a: attrPrefix, t: textKey}
}

// Decode xml reader to map[string]any
func Decode(r io.Reader) (map[string]any, error) {
	return mapdec.Decode(r)
}

// DecodeBytes xml bytes to map[string]any
func DecodeBytes(bs []byte) (map[string]any, error) {
	return mapdec.Decode(bytes.NewReader(bs))
}

// DecodeString xml string to map[string]any
func DecodeString(s string) (map[string]any, error) {
	return mapdec.Decode(strings.NewReader(s))
}

// Decode xml string to map[string]any
func (d *MapDecoder) Decode(r io.Reader) (map[string]any, error) {
	n := &node{}

	xd := xml.NewDecoder(r)
	for {
		token, err := xd.Token()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return n.Value, nil
			}
			return nil, err
		}

		switch tok := token.(type) {
		case xml.StartElement:
			label := tok.Name.Local
			if tok.Name.Space != "" {
				label = strings.ToLower(path.Base(tok.Name.Space)) + ":" + tok.Name.Local
			}

			n = &node{
				Label:  label,
				Parent: n,
			}

			if len(tok.Attr) > 0 {
				m := make(map[string]any, len(tok.Attr))
				for _, attr := range tok.Attr {
					if len(attr.Name.Space) > 0 {
						m[d.a+attr.Name.Space+":"+attr.Name.Local] = attr.Value
					} else {
						m[d.a+attr.Name.Local] = attr.Value
					}
				}
				n.Value = m
			}

		case xml.CharData:
			if n.Parent != nil {
				n.Text += string(tok)
			} else if !str.IsWhitespace(string(tok)) {
				return nil, ErrInvalidRoot
			}

		case xml.EndElement:
			if len(n.Value) > 0 {
				if !str.IsWhitespace(n.Text) {
					n.Value[d.t] = n.Text
				}
			}
			n.Parent.AddChild(n)
			n = n.Parent
		}
	}
}

type node struct {
	Label  string
	Value  map[string]any
	Text   string
	Parent *node
}

func (n *node) Data() any {
	if len(n.Value) != 0 {
		return n.Value
	}
	return n.Text
}

func (n *node) AddChild(c *node) {
	if n.Value == nil {
		n.Value = make(map[string]any)
	}

	if v, ok := n.Value[c.Label]; ok {
		switch item := v.(type) {
		case []any:
			n.Value[c.Label] = append(item, c.Data())
		default:
			n.Value[c.Label] = []any{item, c.Data()}
		}
	} else {
		n.Value[c.Label] = c.Data()
	}
}
