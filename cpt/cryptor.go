package cpt

import "github.com/askasoft/pango/str"

type Cryptor interface {
	SetKey(key string)
	SetIV(iv string)
	EncryptString(str string) (string, error)
	DecryptString(str string) (string, error)
	EncryptData(data []byte) ([]byte, error)
	DecryptData(data []byte) ([]byte, error)
}

// CutPadKey cut key if key's length is greater than 'size',
// or pad key with space if key's length is smaller than 'size'.
func CutPadKey(key string, size int) string {
	if len(key) > 16 {
		key = key[:16]
	} else if len(key) < 16 {
		key = str.PadRight(key, 16, " ")
	}
	return key
}
