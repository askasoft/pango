package cpt

import (
	"bytes"
)

func Pkcs7Pad(data []byte, size int) []byte {
	i := size - (len(data) % size)
	return append(data, bytes.Repeat([]byte{byte(i)}, i)...)
}

func Pkcs7Unpad(data []byte) []byte {
	n := len(data)
	p := int(data[n-1])
	return data[:n-p]
}
