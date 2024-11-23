package xmw

import "github.com/askasoft/pango/cog/hashset"

type stringSet = hashset.HashSet[string]

func newStringSet(ss ...string) *stringSet {
	return hashset.NewHashSet(ss...)
}
