//go:build go1.18
// +build go1.18

package treemap

import (
	"fmt"
	"strings"

	"github.com/askasoft/pango/iox"
	"github.com/askasoft/pango/str"
)

// TreeMapNode is a node of red-black tree
type TreeMapNode[K any, V any] struct {
	color  color
	left   *TreeMapNode[K, V]
	right  *TreeMapNode[K, V]
	parent *TreeMapNode[K, V]
	key    K
	value  V
}

// Key returns the key
func (tn *TreeMapNode[K, V]) Key() K {
	return tn.key
}

// Value returns the key
func (tn *TreeMapNode[K, V]) Value() V {
	return tn.value
}

// SetValue sets the value
func (tn *TreeMapNode[K, V]) SetValue(v V) {
	tn.value = v
}

func (tn *TreeMapNode[K, V]) getLeft() *TreeMapNode[K, V] {
	if tn != nil {
		return tn.left
	}
	return nil
}

func (tn *TreeMapNode[K, V]) getRight() *TreeMapNode[K, V] {
	if tn != nil {
		return tn.right
	}
	return nil
}

func (tn *TreeMapNode[K, V]) getParent() *TreeMapNode[K, V] {
	if tn != nil {
		return tn.parent
	}
	return nil
}

func (tn *TreeMapNode[K, V]) getGrandParent() *TreeMapNode[K, V] {
	if tn != nil && tn.parent != nil {
		return tn.parent.parent
	}
	return nil
}

func (tn *TreeMapNode[K, V]) getColor() color {
	if tn == nil {
		return black
	}
	return tn.color
}

func (tn *TreeMapNode[K, V]) setColor(c color) {
	if tn != nil {
		tn.color = c
	}
}

// prev returns the previous node or nil.
func (tn *TreeMapNode[K, V]) prev() *TreeMapNode[K, V] {
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
func (tn *TreeMapNode[K, V]) next() *TreeMapNode[K, V] {
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
func (tn *TreeMapNode[K, V]) String() string {
	return fmt.Sprintf("%v => %v", tn.key, tn.value)
}

const (
	tmColor = 1 << iota
	tmPoint
	tmValue
)

func (tn *TreeMapNode[K, V]) graph(flag int) string {
	if tn == nil {
		return "(empty)"
	}

	sb := &strings.Builder{}
	tn.output(sb, "", true, flag)
	return sb.String()
}

func (tn *TreeMapNode[K, V]) output(sb *strings.Builder, prefix string, tail bool, flag int) {
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
