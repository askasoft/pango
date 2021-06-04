package col

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapItemString(t *testing.T) {
	assert.Equal(t, "a => b", (&MapItem{"a", "b"}).String())
}
