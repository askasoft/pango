package col

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListItemString(t *testing.T) {
	assert.Equal(t, "a", (&ListItem{Value: "a"}).String())
}
