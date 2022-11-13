package cpt

type Cryptor interface {
	SetKey(key string)
	SetIV(iv string)
	EncryptString(str string) (string, error)
	DecryptString(str string) (string, error)
	EncryptData(data []byte) ([]byte, error)
	DecryptData(data []byte) ([]byte, error)
}
