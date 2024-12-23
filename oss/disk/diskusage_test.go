package disk

import (
	"fmt"
	"testing"

	"github.com/askasoft/pango/num"
)

func TestGetDiskUsage(t *testing.T) {
	du, err := GetDiskUsage(".")

	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("Total:", num.HumanSize(du.Total))
	fmt.Println("Avail:", num.HumanSize(du.Available))
	fmt.Println("Free:", num.HumanSize(du.Free))
	fmt.Println("Used:", num.HumanSize(du.Used()))
	fmt.Println("Usage:", num.FtoaWithDigits(du.Usage(), 2), "%")
}
