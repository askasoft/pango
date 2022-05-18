package col

import (
	"fmt"
	"strings"

	"github.com/pandafw/pango/iox"
	"github.com/pandafw/pango/str"
)

//-----------------------------------------------------

// TreeMapNode is a node of red-black tree
type TreeMapNode struct {
	color  color
	left   *TreeMapNode
	right  *TreeMapNode
	parent *TreeMapNode
	key    K
	value  V
}

// Key returns the key
func (tn *TreeMapNode) Key() K {
	return tn.key
}

// Value returns the key
func (tn *TreeMapNode) Value() V {
	return tn.value
}

// SetValue sets the value
func (tn *TreeMapNode) SetValue(v V) {
	tn.value = v
}

func (tn *TreeMapNode) getLeft() *TreeMapNode {
	if tn != nil {
		return tn.left
	}
	return nil
}

func (tn *TreeMapNode) getRight() *TreeMapNode {
	if tn != nil {
		return tn.right
	}
	return nil
}

func (tn *TreeMapNode) getParent() *TreeMapNode {
	if tn != nil {
		return tn.parent
	}
	return nil
}

func (tn *TreeMapNode) getGrandParent() *TreeMapNode {
	if tn != nil && tn.parent != nil {
		return tn.parent.parent
	}
	return nil
}

func (tn *TreeMapNode) getColor() color {
	if tn == nil {
		return black
	}
	return tn.color
}

func (tn *TreeMapNode) setColor(c color) {
	if tn != nil {
		tn.color = c
	}
}

// prev returns the previous node or nil.
func (tn *TreeMapNode) prev() *TreeMapNode {
	if tn == nil {
		return nil
	}

	if tn.left != nil {
		p := tn.left
		for p.right != nil {
			p = p.right
		}
		return p
	}

	c := tn
	p := tn.parent
	for p != nil && c == p.left {
		c = p
		p = p.parent
	}
	return p
}

// next returns the next node or nil.
func (tn *TreeMapNode) next() *TreeMapNode {
	if tn == nil {
		return nil
	}

	if tn.right != nil {
		n := tn.right
		for n.left != nil {
			n = n.left
		}
		return n
	}

	c := tn
	n := tn.parent
	for n != nil && c == n.right {
		c = n
		n = n.parent
	}
	return n
}

// String print the key/value node to string
func (tn *TreeMapNode) String() string {
	return fmt.Sprintf("%v => %v", tn.key, tn.value)
}

const (
	tmColor = 1 << iota
	tmPoint
	tmValue
)

func (tn *TreeMapNode) graph(flag int) string {
	if tn == nil {
		return "(empty)"
	}

	sb := &strings.Builder{}
	tn.output(sb, "", true, flag)
	return sb.String()
}

func (tn *TreeMapNode) output(sb *strings.Builder, prefix string, tail bool, flag int) {
	if tn.right != nil {
		newPrefix := prefix
		if tail {
			newPrefix += "│   "
		} else {
			newPrefix += "    "
		}
		tn.right.output(sb, newPrefix, false, flag)
	}

	sb.WriteString(prefix)
	if tail {
		sb.WriteString("└── ")
	} else {
		sb.WriteString("┌── ")
	}

	if flag&tmColor == tmColor {
		sb.WriteString(fmt.Sprintf("(%v) ", tn.color))
	}
	sb.WriteString(fmt.Sprint(tn.key))
	if flag&tmPoint == tmPoint {
		sb.WriteString(fmt.Sprintf(" (%p)", tn))
	}
	if flag&tmValue == tmValue {
		v := str.RemoveAny(fmt.Sprint(tn.value), "\r\n")
		sb.WriteString(" => ")
		sb.WriteString(v)
	}
	sb.WriteString(iox.EOL)

	if tn.left != nil {
		newPrefix := prefix
		if tail {
			newPrefix += "    "
		} else {
			newPrefix += "│   "
		}
		tn.left.output(sb, newPrefix, true, flag)
	}
}
