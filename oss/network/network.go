package network

import "time"

// NetworkStats represents network statistics
type NetworkStats struct {
	Name        string
	Received    uint64 // total number of bytes of data received
	Transmitted uint64 // total number of bytes of data transmitted
}

func (ns *NetworkStats) Subtract(s *NetworkStats) {
	ns.Received -= s.Received
	ns.Transmitted -= s.Transmitted
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

type NetworkStatsDelta struct {
	NetworkStats
	Delta time.Duration
}

type NetworksStatsDelta []NetworkStatsDelta

// GetNetworkStatsDelta get network statistics between delta duration
func GetNetworkStatsDelta(delta time.Duration) (nssd NetworksStatsDelta, err error) {
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

	nssd = make(NetworksStatsDelta, len(nss2))
	for i, ns := range nss2 {
		nssd[i].NetworkStats = ns
		nssd[i].Delta = delta
	}
	return
}
