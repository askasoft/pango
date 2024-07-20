package internal

import (
	"fmt"
	"os"
)

func Perror(err error) {
	fmt.Fprintln(os.Stderr, err.Error())
}

func Perrorf(format string, a ...any) {
	fmt.Fprintf(os.Stderr, format, a...)
	fmt.Fprintln(os.Stderr)
}
