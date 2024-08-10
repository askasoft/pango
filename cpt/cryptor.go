package cpt

import "github.com/askasoft/pango/str"

type Cryptor interface {
	EncryptString(str string) (string, error)
	DecryptString(str string) (string, error)
	EncryptBytes(data []byte) ([]byte, error)
	DecryptBytes(data []byte) ([]byte, error)
}

// CutPadKey cut key if key's length is greater than 'size',
// or pad key with space if key's length is smaller than 'size'.
func CutPadKey(key string, size int) string {
	if len(key) > size {
		key = key[:size]
	} else if len(key) < size {
		key = str.PadRight(key, size, " ")
	}
	return key
}
