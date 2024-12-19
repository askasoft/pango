package uptime

import (
	"syscall"
	"time"
)

var (
	kernel = syscall.MustLoadDLL("kernel32.dll")

	getTickCount = kernel.MustFindProc("GetTickCount64")
)

func GetUptime() (time.Duration, error) {
	ret, _, err := getTickCount.Call()
	if errno, ok := err.(syscall.Errno); !ok || errno != 0 { //nolint: errorlint
		return time.Duration(0), err
	}
	return time.Duration(ret) * time.Millisecond, nil
}
