package num

import (
	"fmt"
)

func ExampleHumanSize() {
	fmt.Println(HumanSize(1000))
	fmt.Println(HumanSize(1024))
	fmt.Println(HumanSize(1000000))
	fmt.Println(HumanSize(1048576))
	fmt.Println(HumanSize(2 * MB))
	fmt.Println(HumanSize(float64(3.42 * GB)))
	fmt.Println(HumanSize(float64(5.372 * TB)))
	fmt.Println(HumanSize(float64(2.22 * PB)))
}

func ExampleParseSize() {
	fmt.Println(ParseSize("32"))
	fmt.Println(ParseSize("32b"))
	fmt.Println(ParseSize("32B"))
	fmt.Println(ParseSize("32k"))
	fmt.Println(ParseSize("32K"))
	fmt.Println(ParseSize("32kb"))
	fmt.Println(ParseSize("32Kb"))
	fmt.Println(ParseSize("32Mb"))
	fmt.Println(ParseSize("32Gb"))
	fmt.Println(ParseSize("32Tb"))
	fmt.Println(ParseSize("32Pb"))
}
