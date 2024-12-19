package disk

import (
	"fmt"
	"testing"

	"github.com/askasoft/pango/num"
)

func TestGetDiskStats(t *testing.T) {
	ds, err := GetDiskStats(".")

	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("Total:", num.HumanSize(ds.Total))
	fmt.Println("Avail:", num.HumanSize(ds.Available))
	fmt.Println("Free:", num.HumanSize(ds.Free))
	fmt.Println("Used:", num.HumanSize(ds.Used()))
	fmt.Println("Usage:", num.FtoaWithDigits(ds.Usage(), 2), "%")
}
