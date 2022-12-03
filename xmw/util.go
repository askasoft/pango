//go:build go1.18
// +build go1.18

package xmw

import (
	"github.com/pandafw/pango/cog"
)

type stringSet = cog.HashSet[string]

func newStringSet(ss ...string) *stringSet {
	return cog.NewHashSet(ss...)
}
