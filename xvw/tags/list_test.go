package tags

import (
	"testing"

	"github.com/askasoft/pango/cog"
	"github.com/askasoft/pango/cog/arraylist"
	"github.com/askasoft/pango/cog/linkedhashmap"
)

func TestAsList(t *testing.T) {
	cs := []any{
		linkedhashmap.NewLinkedHashMap(cog.KV("1", "a")),
		linkedhashmap.NewLinkedHashMap(cog.KV(1, "a")),
		map[string]string{"1": "a"},
		map[int]string{1: "a"},
		[]int{1, 2, 3},
		[]string{"a", "b"},
		arraylist.NewArrayList("a", "b"),
		arraylist.NewArrayList(1, 2),
	}

	for _, c := range cs {
		_ = AsList(c)
	}
}
