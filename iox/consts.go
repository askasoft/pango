package iox

import "runtime"

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
