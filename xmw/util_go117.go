//go:build !go1.18
// +build !go1.18

package xmw

import "github.com/pandafw/pango/col"

type stringSet = col.HashSet

func newStringSet(ss ...string) *stringSet {
	return col.NewStringHashSet(ss...)
}
