package email

import (
	"fmt"
	"strings"
	"testing"
)

func TestEncodeString(t *testing.T) {
	fmt.Println(encodeString(strings.Repeat(" 一二三四五", 1)))
	fmt.Println(encodeString(strings.Repeat(" 一二三四五", 10)))
}
