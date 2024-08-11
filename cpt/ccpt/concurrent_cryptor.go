package ccpt

import "github.com/askasoft/pango/cpt"

type cryptor struct {
	fe func() cpt.Encryptor
	fd func() cpt.Decryptor
}

func (c *cryptor) EncryptString(src string) (string, error) {
	return c.fe().EncryptString(src)
}

func (c *cryptor) EncryptBytes(src []byte) ([]byte, error) {
	return c.fe().EncryptBytes(src)
}

func (c *cryptor) DecryptString(src string) (string, error) {
	return c.fd().DecryptString(src)
}

func (c *cryptor) DecryptBytes(src []byte) ([]byte, error) {
	return c.fd().DecryptBytes(src)
}
