package netx

import (
	"encoding/binary"
	"net"
	"reflect"
	"testing"
)

func TestIPv4ToInt(t *testing.T) {
	tests := []struct {
		name string
		ip   net.IP
		want uint32
	}{
		{"IPv4 address", net.ParseIP("192.168.0.1"), binary.BigEndian.Uint32(net.ParseIP("192.168.0.1").To4())},
		{"IPv6 address returns 0", net.ParseIP("2001:db8::1"), 0},
		{"nil IP returns 0", nil, 0},
		{"invalid IPv4 (too short) returns 0", net.IP{127, 0, 0}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IPv4ToInt(tt.ip)
			if got != tt.want {
				t.Errorf("IPv4ToInt(%v) = %v; want %v", tt.ip, got, tt.want)
			}
		})
	}
}

func TestParseIP(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want net.IP
	}{
		{"valid IPv4", "192.168.1.1", net.ParseIP("192.168.1.1").To4()},
		{"valid IPv6", "2001:db8::1", net.ParseIP("2001:db8::1")},
		{"invalid IP", "999.999.999.999", nil},
		{"empty string", "", nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseIP(tt.in)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseIP(%q) = %v; want %v", tt.in, got, tt.want)
			}
		})
	}
}

func TestParseCIDR(t *testing.T) {
	tests := []struct {
		name     string
		in       string
		wantIP   string
		wantMask string
		wantErr  bool
	}{
		{"IPv4 with /24", "192.168.1.1/24", "192.168.1.1", "ffffff00", false},
		{"IPv6 with /64", "2001:db8::1/64", "2001:db8::1", "ffffffffffffffff0000000000000000", false},
		{"IPv4 without mask", "10.0.0.1", "10.0.0.1", "ffffffff", false},
		{"IPv6 without mask", "2001:db8::1", "2001:db8::1", "ffffffffffffffffffffffffffffffff", false},
		{"invalid IP", "999.999.999.999/24", "", "", true},
		{"empty string", "", "", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ip, ipnet, err := ParseCIDR(tt.in)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseCIDR(%q) error = %v; wantErr %v", tt.in, err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if ip.String() != tt.wantIP {
				t.Errorf("ParseCIDR(%q) got IP = %v; want %v", tt.in, ip, tt.wantIP)
			}
			if ipnet == nil {
				t.Errorf("ParseCIDR(%q) returned nil ipnet", tt.in)
			} else if ipnet.Mask.String() != tt.wantMask {
				t.Errorf("ParseCIDR(%q) got NetMask = %v; want %v", tt.in, ipnet.Mask.String(), tt.wantMask)
			}
		})
	}
}

func TestParseCIDRs(t *testing.T) {
	tests := []struct {
		name    string
		in      []string
		wantLen int
		wantErr bool
	}{
		{"multiple valid CIDRs", []string{"192.168.1.0/24", "10.0.0.0/8"}, 2, false},
		{"mixed IPv4 and IPv6", []string{"10.0.0.0/8", "2001:db8::/32"}, 2, false},
		{"invalid CIDR", []string{"192.168.1.0/33"}, 0, true},
		{"empty list", []string{}, 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseCIDRs(tt.in)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseCIDRs(%v) error = %v; wantErr %v", tt.in, err, tt.wantErr)
			}
			if !tt.wantErr && len(got) != tt.wantLen {
				t.Errorf("ParseCIDRs(%v) returned %d items; want %d", tt.in, len(got), tt.wantLen)
			}
		})
	}
}
