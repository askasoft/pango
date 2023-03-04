package iox

import "runtime"

// BOM '\uFEFF'
const BOM = '\uFEFF'

// CR "\r"
const CR = "\r"

// LF "\n"
const LF = "\n"

// CRLF "\r\n"
const CRLF = "\r\n"

// EOL windows: "\r\n" other: "\n"
var EOL = geteol()

func geteol() string {
	if runtime.GOOS == "windows" {
		return CRLF
	}
	return LF
}
