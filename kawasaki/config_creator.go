package kawasaki

import (
	"fmt"
	"net"
)

const maxHandleBeforeTruncation = 8

type NetworkConfig struct {
	HostIntf      string
	ContainerIntf string
	BridgeName    string
	BridgeIP      net.IP
	ContainerIP   net.IP
	Subnet        *net.IPNet
	Mtu           int
}

type Creator struct{}

func NewConfigCreator() *Creator {
	return &Creator{}
}

func (c *Creator) Create(handle, spec string) (NetworkConfig, error) {
	ip, subnet, err := net.ParseCIDR(spec)
	if err != nil {
		return NetworkConfig{}, err
	}

	return NetworkConfig{
		HostIntf:      fmt.Sprintf("w-%s-0", truncate(handle)),
		ContainerIntf: fmt.Sprintf("w-%s-1", truncate(handle)),
		BridgeName:    "guardian-bridge",
		ContainerIP:   ip,
		BridgeIP:      bridge(subnet),
		Subnet:        subnet,
		Mtu:           1500,
	}, nil
}

func truncate(handle string) string {
	if len(handle) > maxHandleBeforeTruncation {
		return handle[:maxHandleBeforeTruncation]
	}

	return handle
}

func bridge(subnet *net.IPNet) net.IP {
	return next(subnet.IP)
}

func next(ip net.IP) net.IP {
	next := clone(ip)
	for i := len(next) - 1; i >= 0; i-- {
		next[i]++
		if next[i] != 0 {
			return next
		}
	}

	panic("overflowed maximum IP")
}

func clone(ip net.IP) net.IP {
	clone := make([]byte, len(ip))
	copy(clone, ip)
	return clone
}
