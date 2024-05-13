package tags

import "github.com/askasoft/pango/cog"

type List[K any, V any] interface {
	cog.Iterable2[K, V]
}
