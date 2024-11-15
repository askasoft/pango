package internal

import (
	"fmt"
	"os"
)

func Perror(a any) {
	fmt.Fprintln(os.Stderr, a)
}

func Perrorf(format string, a ...any) {
	fmt.Fprintf(os.Stderr, format, a...)
	fmt.Fprintln(os.Stderr)
}
