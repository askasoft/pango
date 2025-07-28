package xmw

import (
	"github.com/askasoft/pango/xin"
)

type testAccount struct {
	username, password string
}

func (ta *testAccount) GetUsername() string {
	return ta.username
}

func (ta *testAccount) GetPassword() string {
	return ta.password
}

type testAccounts map[string]*testAccount

func (tas testAccounts) FindUser(c *xin.Context, username, password string) (AuthUser, error) {
	if ta, ok := tas[username]; ok && (password == "" || ta.password == password) {
		return ta, nil
	}
	return nil, nil
}
