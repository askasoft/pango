//go:build go1.18
// +build go1.18

package treeset

import (
	"fmt"
)

// color node color
type color byte

// Red Black color
const (
	red   color = 'R'
	black color = 'B'
)

func (c color) String() string {
	return fmt.Sprint(string(c))
}
