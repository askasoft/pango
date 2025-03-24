package netutil

import (
	"net"

	"github.com/askasoft/pango/str"
)

// parseIP parse a string representation of an IP and returns a net.IP with the
// minimum byte representation or nil if input is invalid.
func parseIP(ip string) net.IP {
	parsedIP := net.ParseIP(ip)

	if ipv4 := parsedIP.To4(); ipv4 != nil {
		// return ip in a 4-byte representation
		return ipv4
	}

	// return ip in a 16-byte representation or nil
	return parsedIP
}

func ParseCIDRs(cidrs []string) ([]*net.IPNet, error) {
	ipnets := make([]*net.IPNet, 0, len(cidrs))
	for _, cidr := range cidrs {
		if !str.ContainsByte(cidr, '/') {
			ip := parseIP(cidr)
			if ip == nil {
				return nil, &net.ParseError{Type: "IP address", Text: cidr}
			}

			switch len(ip) {
			case net.IPv4len:
				cidr += "/32"
			case net.IPv6len:
				cidr += "/128"
			}
		}

		_, ipnet, err := net.ParseCIDR(cidr)
		if err != nil {
			return nil, err
		}
		ipnets = append(ipnets, ipnet)
	}

	return ipnets, nil
}
