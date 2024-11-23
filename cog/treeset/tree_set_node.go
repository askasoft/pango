package treeset

import (
	"fmt"
	"strings"

	"github.com/askasoft/pango/iox"
)

// treeSetNode is a node of red-black tree
type treeSetNode[T any] struct {
	color  color
	left   *treeSetNode[T]
	right  *treeSetNode[T]
	parent *treeSetNode[T]
	value  T
}

func (tn *treeSetNode[T]) getLeft() *treeSetNode[T] {
	if tn != nil {
		return tn.left
	}
	return nil
}

func (tn *treeSetNode[T]) getRight() *treeSetNode[T] {
	if tn != nil {
		return tn.right
	}
	return nil
}

func (tn *treeSetNode[T]) getParent() *treeSetNode[T] {
	if tn != nil {
		return tn.parent
	}
	return nil
}

func (tn *treeSetNode[T]) getGrandParent() *treeSetNode[T] {
	if tn != nil && tn.parent != nil {
		return tn.parent.parent
	}
	return nil
}

func (tn *treeSetNode[T]) getColor() color {
	if tn == nil {
		return black
	}
	return tn.color
}

func (tn *treeSetNode[T]) setColor(c color) {
	if tn != nil {
		tn.color = c
	}
}

// prev returns the previous node or nil.
func (tn *treeSetNode[T]) prev() *treeSetNode[T] {
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
func (tn *treeSetNode[T]) next() *treeSetNode[T] {
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

// String print the set item to string
func (tn *treeSetNode[T]) String() string {
	return fmt.Sprint(tn.value)
}

const (
	tsColor = 1 << iota
	tsPoint
)

func (tn *treeSetNode[T]) graph(flag int) string {
	if tn == nil {
		return "(empty)"
	}

	sb := &strings.Builder{}
	tn.output(sb, "", true, flag)
	return sb.String()
}

func (tn *treeSetNode[T]) output(sb *strings.Builder, prefix string, tail bool, flag int) {
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

	if flag&tsColor == tsColor {
		fmt.Fprintf(sb, "(%v) ", tn.color)
	}
	fmt.Fprint(sb, tn.value)
	if flag&tsPoint == tsPoint {
		fmt.Fprintf(sb, " (%p)", tn)
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
