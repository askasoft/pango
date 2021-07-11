package gingzip

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGzipProxiedFlag(t *testing.T) {
	i := 0
	assert.Equal(t, i, int(ProxiedOff))
	i = 1
	assert.Equal(t, i, int(ProxiedAny))
	i <<= 1
	assert.Equal(t, i, int(ProxiedAuth))
	i <<= 1
	assert.Equal(t, i, int(ProxiedExpired))
	i <<= 1
	assert.Equal(t, i, int(ProxiedNoCache))
	i <<= 1
	assert.Equal(t, i, int(ProxiedNoStore))
	i <<= 1
	assert.Equal(t, i, int(ProxiedPrivate))
	i <<= 1
	assert.Equal(t, i, int(ProxiedNoLastModified))
	i <<= 1
	assert.Equal(t, i, int(ProxiedNoETag))
}
