package loadavg

import (
	"errors"
)

func GetLoadAvg() (la LoadAvg, err error) {
	err = errors.New("loadavg for windows is not supported")
	return
}
