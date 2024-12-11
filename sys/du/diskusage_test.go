package du

import (
	"fmt"
	"testing"

	"github.com/askasoft/pango/num"
)

func TestNewDiskUsage(t *testing.T) {
	usage := NewDiskUsage(".")

	fmt.Println("Total:", num.HumanSize(usage.Total()))
	fmt.Println("Avail:", num.HumanSize(usage.Available()))
	fmt.Println("Free:", num.HumanSize(usage.Free()))
	fmt.Println("Used:", num.HumanSize(usage.Used()))
	fmt.Println("Usage:", num.FtoaWithDigits(usage.Usage()*100, 2), "%")
}
