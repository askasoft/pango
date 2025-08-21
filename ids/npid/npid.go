// Package npid provides a ID generator based on machine network private IP and process ID.
package npid

import (
	"math"
	"net"
	"os"

	"github.com/askasoft/pango/net/netx"
)

// New returns a ID based on private IPv4 and process ID.
// The default ID format:
// +--------------------------------------------------+
// | (64 - x - y) Bit Unused | x Bit IPv4 | y Bit PID |
// +--------------------------------------------------+
func New(ip4Bits, pidBits int) int64 {
	if ip4Bits+pidBits > 63 {
		panic("npid: IPv4+PID bits must less than 64")
	}

	var ip4 uint32
	var pid int

	if ip4Bits > 0 {
		ip4 = getFirstPrivateIPv4Addr()
	}
	if pidBits > 0 {
		pid = os.Getpid()
	}

	return Build(ip4, pid, ip4Bits, pidBits)
}

func Build(ip4 uint32, pid int, ip4Bits, pidBits int) int64 {
	ip4Mask := ^(int64(-1) << ip4Bits)
	pidMask := ^(int64(-1) << pidBits)

	return ((int64(ip4) & ip4Mask) << pidBits) | (int64(pid) & pidMask)
}

func getFirstPrivateIPv4Addr() (ip4 uint32) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return
	}

	idx := math.MaxInt
	for _, i := range ifaces {
		if ip4 != 0 && i.Index > idx {
			continue
		}
		if i.Flags&net.FlagLoopback != 0 || i.Flags&net.FlagRunning == 0 {
			continue
		}

		addrs, err := i.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && ipnet.IP.IsPrivate() {
				if ipv4 := ipnet.IP.To4(); ipv4 != nil {
					idx, ip4 = i.Index, netx.IPv4ToInt(ipv4)
					break
				}
			}
		}
	}
	return
}
