package binding

import (
	"testing"
)

func TestHeaderInterface(t *testing.T) {
	var _ setter = headerSource(nil)
}
