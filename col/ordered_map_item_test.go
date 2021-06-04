package col

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrderedMapItemString(t *testing.T) {
	assert.Equal(t, "a => b", (&OrderedMapItem{MapItem{"a", "b"}, nil}).String())
}
