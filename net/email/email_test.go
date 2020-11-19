package email

import (
	"testing"
)

func TestEmailSetFrom(t *testing.T) {
	email := &Email{}
	err := email.SetFrom("王 <ou-x@test.com>")
	if err != nil {
		t.Error(err)
	}
}
