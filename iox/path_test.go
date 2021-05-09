package iox

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathMatch1(t *testing.T) {
	b, err := filepath.Match("dir/*.txt", "dir/a.txt")
	assert.Nil(t, err)
	assert.True(t, b)
}

func TestPathMatch2(t *testing.T) {
	b, err := filepath.Match("dir/**/*.txt", "dir/a.txt")
	assert.Nil(t, err)
	assert.False(t, b)
}

func TestPathMatch3(t *testing.T) {
	b, err := filepath.Match("dir/**/*.txt", "dir/3/a.txt")
	assert.Nil(t, err)
	assert.True(t, b)
}

func TestPathMatch4(t *testing.T) {
	b, err := filepath.Match("dir/**/*.txt", "dir/3/5/a.txt")
	assert.Nil(t, err)

	// why??
	if runtime.GOOS == "windows" {
		assert.True(t, b)
	} else {
		assert.False(t, b)
	}
}

func TestPathMatch5(t *testing.T) {
	b, err := filepath.Match("**/*.txt", "a.txt")
	assert.Nil(t, err)
	assert.False(t, b)
}

func TestPathMatch6(t *testing.T) {
	b, err := filepath.Match("**/*.txt", "a/a.txt")
	assert.Nil(t, err)
	assert.True(t, b)

	b, err = filepath.Match("**/*.txt", "a\\a.txt")
	assert.Nil(t, err)
	assert.False(t, b)

	if runtime.GOOS == "windows" {
		b, err = filepath.Match("**\\*.txt", "a\\a.txt")
		assert.Nil(t, err)
		assert.True(t, b)
	}
}
