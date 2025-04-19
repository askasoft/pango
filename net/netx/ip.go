package netx

import (
	"encoding/binary"
	"net"

	"github.com/askasoft/pango/str"
)

// IPv4ToInt converts IP address of version 4 from net.IP to uint32 representation.
func IPv4ToInt(ip net.IP) uint32 {
	if ipv4 := ip.To4(); ipv4 != nil {
		return binary.BigEndian.Uint32(ipv4)
	}
	return 0
}

// ParseIP parse a string representation of an IP and returns a net.IP with the
// minimum byte representation or nil if input is invalid.
func ParseIP(ip string) net.IP {
	parsedIP := net.ParseIP(ip)

	if ipv4 := parsedIP.To4(); ipv4 != nil {
		// return ip in a 4-byte representation
		return ipv4
	}

	// return ip in a 16-byte representation or nil
	return parsedIP
}

// ParseCIDR parses s as a CIDR notation IP address and prefix length,
// like "192.0.2.0/24" or "2001:db8::/32", as defined in
// RFC 4632 and RFC 4291.
// If no mask supplyed, the default "/32" or "/128" will be appended.
// It returns the IP address and the network implied by the IP and
// prefix length.
// For example, ParseCIDR("192.0.2.1/24") returns the IP address
// 192.0.2.1 and the network 192.0.2.0/24.
func ParseCIDR(cidr string) (net.IP, *net.IPNet, error) {
	if !str.ContainsByte(cidr, '/') {
		ip := ParseIP(cidr)
		if ip == nil {
			return nil, nil, &net.ParseError{Type: "IP address", Text: cidr}
		}

		switch len(ip) {
		case net.IPv4len:
			cidr += "/32"
		case net.IPv6len:
			cidr += "/128"
		}
	}

	return net.ParseCIDR(cidr)
}

// ParseCIDRs parse a string representation of an IP and returns a net.IP with the
// minimum byte representation or nil if input is invalid.
func ParseCIDRs(cidrs []string) ([]*net.IPNet, error) {
	ipnets := make([]*net.IPNet, 0, len(cidrs))

	for _, cidr := range cidrs {
		_, ipnet, err := ParseCIDR(cidr)
		if err != nil {
			return nil, err
		}
		ipnets = append(ipnets, ipnet)
	}

	return ipnets, nil
}
