package xmlx

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"path"
	"strings"
)

var (
	// ErrInvalidDocument invalid document err
	ErrInvalidDocument = errors.New("xmlx: invalid document")

	// ErrInvalidRoot data at the root level is invalid err
	ErrInvalidRoot = errors.New("xmlx: data at the root level is invalid")

	mapdec = NewMapDecoder("@", "#text")
)

type node struct {
	Parent  *node
	Value   map[string]any
	Attrs   []xml.Attr
	Label   string
	Space   string
	Text    string
	HasMany bool
}

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
	decoder := xml.NewDecoder(r)
	n := &node{}
	stack := make([]*node, 0)

	for {
		token, err := decoder.Token()
		if err != nil && !errors.Is(err, io.EOF) {
			return nil, err
		}

		if token == nil {
			break
		}

		switch tok := token.(type) {
		case xml.StartElement:
			label := tok.Name.Local
			if tok.Name.Space != "" {
				label = fmt.Sprintf("%s:%s", strings.ToLower(path.Base(tok.Name.Space)), tok.Name.Local)
			}
			n = &node{
				Label:  label,
				Space:  tok.Name.Space,
				Parent: n,
				Value:  map[string]any{label: map[string]any{}},
				Attrs:  tok.Attr,
			}

			d.setAttrs(n, &tok)
			stack = append(stack, n)

			if n.Parent != nil {
				n.Parent.HasMany = true
			}

		case xml.CharData:
			data := strings.TrimSpace(string(tok))
			if len(stack) > 0 {
				stack[len(stack)-1].Text = data
			} else if len(data) > 0 {
				return nil, ErrInvalidRoot
			}

		case xml.EndElement:
			length := len(stack)
			stack, n = stack[:length-1], stack[length-1]

			if !n.HasMany {
				if len(n.Attrs) > 0 {
					m := n.Value[n.Label].(map[string]any)
					m[d.t] = n.Text
				} else {
					n.Value[n.Label] = n.Text
				}
			}

			if len(stack) == 0 {
				return n.Value, nil
			}

			d.setNodeValue(n)
			n = n.Parent
		}
	}

	return nil, ErrInvalidDocument
}

func (d *MapDecoder) setAttrs(n *node, tok *xml.StartElement) {
	if len(tok.Attr) > 0 {
		m := make(map[string]any)
		for _, attr := range tok.Attr {
			if len(attr.Name.Space) > 0 {
				m[d.a+attr.Name.Space+":"+attr.Name.Local] = attr.Value
			} else {
				m[d.a+attr.Name.Local] = attr.Value
			}
		}
		n.Value[tok.Name.Local] = m
	}
}

func (d *MapDecoder) setNodeValue(n *node) {
	if v, ok := n.Parent.Value[n.Parent.Label]; ok {
		m := v.(map[string]any)
		if v, ok = m[n.Label]; ok {
			switch item := v.(type) {
			case string:
				m[n.Label] = []string{item, n.Value[n.Label].(string)}
			case []string:
				m[n.Label] = append(item, n.Value[n.Label].(string))
			case map[string]any:
				vm := d.getMap(n)
				if vm != nil {
					m[n.Label] = []map[string]any{item, vm}
				}
			case []map[string]any:
				vm := d.getMap(n)
				if vm != nil {
					m[n.Label] = append(item, vm)
				}
			}
		} else {
			m[n.Label] = n.Value[n.Label]
		}
	} else {
		n.Parent.Value[n.Parent.Label] = n.Value[n.Label]
	}
}

func (d *MapDecoder) getMap(node *node) map[string]any {
	if v, ok := node.Value[node.Label]; ok {
		switch v.(type) {
		case string:
			return map[string]any{node.Label: v}
		case map[string]any:
			return node.Value[node.Label].(map[string]any)
		}
	}

	return nil
}
