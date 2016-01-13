package devices_test

import (
	"net"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vishvananda/netlink"

	"testing"
)

func TestDevices(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Devices Suite")
}

func cleanup(intfName string) error {
	if _, err := net.InterfaceByName(intfName); err == nil {
		link, err := netlink.LinkByName(intfName)
		if err != nil {
			return err
		}
		return netlink.LinkDel(link)
	}
	return nil
}
