package configure_test

import (
	"errors"
	"net"

	"github.com/cloudfoundry-incubator/guardian/kawasaki"
	"github.com/cloudfoundry-incubator/guardian/kawasaki/configure"
	"github.com/cloudfoundry-incubator/guardian/kawasaki/devices/fakedevices"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Container", func() {
	var (
		linkApplyr *fakedevices.FakeLink
		configurer *configure.Container
		config     kawasaki.NetworkConfig
	)

	BeforeEach(func() {
		linkApplyr = &fakedevices.FakeLink{AddIPReturns: make(map[string]error)}
		configurer = &configure.Container{
			Link: linkApplyr,
		}
	})

	Context("when the loopback device does not exist", func() {
		var eth *net.Interface
		BeforeEach(func() {
			linkApplyr.InterfaceByNameFunc = func(name string) (*net.Interface, bool, error) {
				if name != "lo" {
					return eth, true, nil
				}

				return nil, false, nil
			}
		})

		It("returns a wrapped error", func() {
			err := configurer.Apply(config)
			Expect(err).To(MatchError(&configure.FindLinkError{nil, "loopback", "lo"}))
		})

		It("does not attempt to configure other devices", func() {
			Expect(configurer.Apply(config)).ToNot(Succeed())
			Expect(linkApplyr.SetUpCalledWith).ToNot(ContainElement(eth))
		})
	})

	Context("when the loopback exists", func() {
		var lo *net.Interface

		BeforeEach(func() {
			lo = &net.Interface{Name: "lo"}
			linkApplyr.InterfaceByNameFunc = func(name string) (*net.Interface, bool, error) {
				return &net.Interface{Name: name}, true, nil
			}
		})

		It("adds 127.0.0.1/8 as an address", func() {
			ip, subnet, _ := net.ParseCIDR("127.0.0.1/8")
			Expect(configurer.Apply(config)).To(Succeed())
			Expect(linkApplyr.AddIPCalledWith).To(ContainElement(fakedevices.InterfaceIPAndSubnet{lo, ip, subnet}))
		})

		Context("when adding the IP address fails", func() {
			It("returns a wrapped error", func() {
				linkApplyr.AddIPReturns["lo"] = errors.New("o no")
				err := configurer.Apply(config)
				ip, subnet, _ := net.ParseCIDR("127.0.0.1/8")
				Expect(err).To(MatchError(&configure.ConfigureLinkError{errors.New("o no"), "loopback", lo, ip, subnet}))
			})
		})

		It("brings it up", func() {
			Expect(configurer.Apply(config)).To(Succeed())
			Expect(linkApplyr.SetUpCalledWith).To(ContainElement(lo))
		})

		Context("when bringing the link up fails", func() {
			It("returns a wrapped error", func() {
				linkApplyr.SetUpFunc = func(intf *net.Interface) error {
					return errors.New("o no")
				}

				err := configurer.Apply(config)
				Expect(err).To(MatchError(&configure.LinkUpError{errors.New("o no"), lo, "loopback"}))
			})
		})
	})

	Context("when the container interface does not exist", func() {
		BeforeEach(func() {
			linkApplyr.InterfaceByNameFunc = func(name string) (*net.Interface, bool, error) {
				if name == "lo" {
					return &net.Interface{Name: name}, true, nil
				}

				return nil, false, nil
			}
		})

		It("returns a wrapped error", func() {
			config.ContainerIntf = "foo"
			err := configurer.Apply(config)
			Expect(err).To(MatchError(&configure.FindLinkError{nil, "container", "foo"}))
		})
	})

	Context("when the container interface exists", func() {
		BeforeEach(func() {
			linkApplyr.InterfaceByNameFunc = func(name string) (*net.Interface, bool, error) {
				return &net.Interface{Name: name}, true, nil
			}
		})

		It("Adds the requested IP", func() {
			config.ContainerIntf = "foo"
			config.ContainerIP, config.Subnet, _ = net.ParseCIDR("2.3.4.5/6")

			Expect(configurer.Apply(config)).To(Succeed())
			Expect(linkApplyr.AddIPCalledWith).To(ContainElement(fakedevices.InterfaceIPAndSubnet{
				&net.Interface{Name: "foo"},
				config.ContainerIP,
				config.Subnet,
			}))
		})

		Context("when adding the IP fails", func() {
			It("returns a wrapped error", func() {
				linkApplyr.AddIPReturns["foo"] = errors.New("o no")

				config.ContainerIntf = "foo"
				config.ContainerIP, config.Subnet, _ = net.ParseCIDR("2.3.4.5/6")
				err := configurer.Apply(config)
				Expect(err).To(MatchError(&configure.ConfigureLinkError{
					errors.New("o no"),
					"container",
					&net.Interface{Name: "foo"},
					config.ContainerIP,
					config.Subnet,
				}))
			})
		})

		It("Brings the link up", func() {
			config.ContainerIntf = "foo"
			Expect(configurer.Apply(config)).To(Succeed())
			Expect(linkApplyr.SetUpCalledWith).To(ContainElement(&net.Interface{Name: "foo"}))
		})

		Context("when bringing the link up fails", func() {
			It("returns a wrapped error", func() {
				cause := errors.New("who ate my pie?")
				linkApplyr.SetUpFunc = func(iface *net.Interface) error {
					if iface.Name == "foo" {
						return cause
					}

					return nil
				}

				config.ContainerIntf = "foo"
				err := configurer.Apply(config)
				Expect(err).To(MatchError(&configure.LinkUpError{cause, &net.Interface{Name: "foo"}, "container"}))
			})
		})

		It("sets the mtu", func() {
			config.ContainerIntf = "foo"
			config.Mtu = 1234
			Expect(configurer.Apply(config)).To(Succeed())
			Expect(linkApplyr.SetMTUCalledWith.Interface).To(Equal(&net.Interface{Name: "foo"}))
			Expect(linkApplyr.SetMTUCalledWith.MTU).To(Equal(1234))
		})

		Context("when setting the mtu fails", func() {
			It("returns a wrapped error", func() {
				linkApplyr.SetMTUReturns = errors.New("this is NOT the right potato")

				config.ContainerIntf = "foo"
				config.Mtu = 1234
				err := configurer.Apply(config)
				Expect(err).To(MatchError(&configure.MTUError{linkApplyr.SetMTUReturns, &net.Interface{Name: "foo"}, 1234}))
			})
		})

		It("adds a default gateway with the requested IP", func() {
			config.ContainerIntf = "foo"
			config.BridgeIP = net.ParseIP("2.3.4.5")
			Expect(configurer.Apply(config)).To(Succeed())
			Expect(linkApplyr.AddDefaultGWCalledWith.Interface).To(Equal(&net.Interface{Name: "foo"}))
			Expect(linkApplyr.AddDefaultGWCalledWith.IP).To(Equal(net.ParseIP("2.3.4.5")))
		})

		Context("when adding a default gateway fails", func() {
			It("returns a wrapped error", func() {
				linkApplyr.AddDefaultGWReturns = errors.New("this is NOT the right potato")

				config.ContainerIntf = "foo"
				config.BridgeIP = net.ParseIP("2.3.4.5")
				err := configurer.Apply(config)
				Expect(err).To(MatchError(&configure.ConfigureDefaultGWError{linkApplyr.AddDefaultGWReturns, &net.Interface{Name: "foo"}, net.ParseIP("2.3.4.5")}))
			})
		})
	})
})
