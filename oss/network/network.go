package network

import (
	"fmt"
	"time"
)

// NetworkStats represents network statistics
type NetworkStats struct {
	Name        string `json:"name"`        // network interface name
	Received    uint64 `json:"received"`    // total number of bytes of data received
	Transmitted uint64 `json:"transmitted"` // total number of bytes of data transmitted
}

func (ns *NetworkStats) Subtract(s *NetworkStats) {
	ns.Received -= s.Received
	ns.Transmitted -= s.Transmitted
}

func (ns *NetworkStats) String() string {
	return fmt.Sprintf("(%q R: %d, T: %d)", ns.Name, ns.Received, ns.Transmitted)
}

type NetworksStats []NetworkStats

func (nss NetworksStats) Subtract(ss NetworksStats) {
	for _, ns := range nss {
		for _, s := range ss {
			if ns.Name == s.Name {
				ns.Subtract(&s)
				break
			}
		}
	}
}

type NetworkUsage struct {
	NetworkStats
	Delta time.Duration
}

// ReceiveSpeed get receive speed bytes/second
func (nu *NetworkUsage) ReceiveSpeed() float64 {
	if nu.Delta == 0 {
		return 0
	}
	return float64(nu.Received) / nu.Delta.Seconds()
}

// TransmitSpeed get transmit speed bytes/second
func (nu *NetworkUsage) TransmitSpeed() float64 {
	if nu.Delta == 0 {
		return 0
	}
	return float64(nu.Transmitted) / nu.Delta.Seconds()
}

func (nu *NetworkUsage) String() string {
	return fmt.Sprintf("(%q R: %d, T: %d, D: %s)", nu.Name, nu.Received, nu.Transmitted, nu.Delta)
}

type NetworksUsage []NetworkUsage

// GetNetworksUsage get networks usage between delta duration
func GetNetworksUsage(delta time.Duration) (nsu NetworksUsage, err error) {
	var nss1, nss2 NetworksStats

	nss1, err = GetNetworksStats()
	if err != nil {
		return
	}

	time.Sleep(delta)

	nss2, err = GetNetworksStats()
	if err != nil {
		return
	}

	nss2.Subtract(nss1)

	nsu = make(NetworksUsage, len(nss2))
	for i, ns := range nss2 {
		nsu[i].NetworkStats = ns
		nsu[i].Delta = delta
	}
	return
}
