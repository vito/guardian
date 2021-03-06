package iptables_test

import (
	. "github.com/onsi/ginkgo"
)

var _ = Describe("FilterChain", func() {
	// var fakeRunner *fake_command_runner.FakeCommandRunner
	// var subject Chain
	// var useKernelLogging bool

	// JustBeforeEach(func() {
	// 	fakeRunner = fake_command_runner.New()
	// 	subject = NewLoggingChain("foo-bar-baz", useKernelLogging, fakeRunner, lagertest.NewTestLogger("test"))
	// })

	// Describe("PrependFilterRule", func() {
	// 	Context("when all parameters are defaulted", func() {
	// 		It("runs iptables with appropriate parameters", func() {
	// 			Expect(subject.PrependFilterRule(garden.NetOutRule{})).To(Succeed())
	// 			Expect(fakeRunner).To(HaveExecutedSerially(fake_command_runner.CommandSpec{
	// 				Path: "/sbin/iptables",
	// 				Args: []string{"-w", "-I", "foo-bar-baz", "1", "--protocol", "all", "--jump", "RETURN"},
	// 			}))
	// 		})
	// 	})

	// 	Describe("Network", func() {
	// 		Context("when an empty IPRange is specified", func() {
	// 			It("does not limit the range", func() {
	// 				Expect(subject.PrependFilterRule(garden.NetOutRule{
	// 					Networks: []garden.IPRange{
	// 						garden.IPRange{},
	// 					},
	// 				})).To(Succeed())

	// 				Expect(fakeRunner).To(HaveExecutedSerially(fake_command_runner.CommandSpec{
	// 					Path: "/sbin/iptables",
	// 					Args: []string{"-w", "-I", "foo-bar-baz", "1", "--protocol", "all", "--jump", "RETURN"},
	// 				}))
	// 			})
	// 		})

	// 		Context("when a single destination IP is specified", func() {
	// 			It("opens only that IP", func() {
	// 				Expect(subject.PrependFilterRule(garden.NetOutRule{
	// 					Networks: []garden.IPRange{
	// 						{
	// 							Start: net.ParseIP("1.2.3.4"),
	// 						},
	// 					},
	// 				})).To(Succeed())

	// 				Expect(fakeRunner).To(HaveExecutedSerially(fake_command_runner.CommandSpec{
	// 					Path: "/sbin/iptables",
	// 					Args: []string{"-w", "-I", "foo-bar-baz", "1", "--protocol", "all", "--destination", "1.2.3.4", "--jump", "RETURN"},
	// 				}))
	// 			})
	// 		})

	// 		Context("when a multiple destination networks are specified", func() {
	// 			It("opens only that IP", func() {
	// 				Expect(subject.PrependFilterRule(garden.NetOutRule{
	// 					Networks: []garden.IPRange{
	// 						{
	// 							Start: net.ParseIP("1.2.3.4"),
	// 						},
	// 						{
	// 							Start: net.ParseIP("2.2.3.4"),
	// 							End:   net.ParseIP("2.2.3.9"),
	// 						},
	// 					},
	// 				})).To(Succeed())

	// 				Expect(fakeRunner.ExecutedCommands()).To(HaveLen(2))
	// 				Expect(fakeRunner).To(HaveExecutedSerially(
	// 					fake_command_runner.CommandSpec{
	// 						Path: "/sbin/iptables",
	// 						Args: []string{"-w", "-I", "foo-bar-baz", "1", "--protocol", "all", "--destination", "1.2.3.4", "--jump", "RETURN"},
	// 					},
	// 					fake_command_runner.CommandSpec{
	// 						Path: "/sbin/iptables",
	// 						Args: []string{"-w", "-I", "foo-bar-baz", "1", "--protocol", "all", "-m", "iprange", "--dst-range", "2.2.3.4-2.2.3.9", "--jump", "RETURN"},
	// 					},
	// 				))
	// 			})
	// 		})

	// 		Context("when a EndIP is specified without a StartIP", func() {
	// 			It("opens only that IP", func() {
	// 				Expect(subject.PrependFilterRule(garden.NetOutRule{
	// 					Networks: []garden.IPRange{
	// 						{
	// 							End: net.ParseIP("1.2.3.4"),
	// 						},
	// 					},
	// 				})).To(Succeed())

	// 				Expect(fakeRunner).To(HaveExecutedSerially(fake_command_runner.CommandSpec{
	// 					Path: "/sbin/iptables",
	// 					Args: []string{"-w", "-I", "foo-bar-baz", "1", "--protocol", "all", "--destination", "1.2.3.4", "--jump", "RETURN"},
	// 				}))
	// 			})
	// 		})

	// 		Context("when a range of IPs is specified", func() {
	// 			It("opens only the range", func() {
	// 				Expect(subject.PrependFilterRule(garden.NetOutRule{
	// 					Networks: []garden.IPRange{
	// 						{
	// 							net.ParseIP("1.2.3.4"), net.ParseIP("2.3.4.5"),
	// 						},
	// 					},
	// 				})).To(Succeed())

	// 				Expect(fakeRunner).To(HaveExecutedSerially(fake_command_runner.CommandSpec{
	// 					Path: "/sbin/iptables",
	// 					Args: []string{"-w", "-I", "foo-bar-baz", "1", "--protocol", "all", "-m", "iprange", "--dst-range", "1.2.3.4-2.3.4.5", "--jump", "RETURN"},
	// 				}))
	// 			})
	// 		})
	// 	})

	// 	Describe("Ports", func() {
	// 		Context("when a single port is specified", func() {
	// 			It("opens only that port", func() {
	// 				Expect(subject.PrependFilterRule(garden.NetOutRule{
	// 					Protocol: garden.ProtocolTCP,
	// 					Ports: []garden.PortRange{
	// 						garden.PortRangeFromPort(22),
	// 					},
	// 				})).To(Succeed())

	// 				Expect(fakeRunner).To(HaveExecutedSerially(fake_command_runner.CommandSpec{
	// 					Path: "/sbin/iptables",
	// 					Args: []string{"-w", "-I", "foo-bar-baz", "1", "--protocol", "tcp", "--destination-port", "22", "--jump", "RETURN"},
	// 				}))
	// 			})
	// 		})

	// 		Context("when a port range is specified", func() {
	// 			It("opens that port range", func() {
	// 				Expect(subject.PrependFilterRule(garden.NetOutRule{
	// 					Protocol: garden.ProtocolTCP,
	// 					Ports: []garden.PortRange{
	// 						{12, 24},
	// 					},
	// 				})).To(Succeed())

	// 				Expect(fakeRunner).To(HaveExecutedSerially(fake_command_runner.CommandSpec{
	// 					Path: "/sbin/iptables",
	// 					Args: []string{"-w", "-I", "foo-bar-baz", "1", "--protocol", "tcp", "--destination-port", "12:24", "--jump", "RETURN"},
	// 				}))
	// 			})
	// 		})

	// 		Context("when multiple port ranges are specified", func() {
	// 			It("opens those port ranges", func() {
	// 				Expect(subject.PrependFilterRule(garden.NetOutRule{
	// 					Protocol: garden.ProtocolTCP,
	// 					Ports: []garden.PortRange{
	// 						{12, 24},
	// 						{64, 942},
	// 					},
	// 				})).To(Succeed())

	// 				Expect(fakeRunner).To(HaveExecutedSerially(
	// 					fake_command_runner.CommandSpec{
	// 						Path: "/sbin/iptables",
	// 						Args: []string{"-w", "-I", "foo-bar-baz", "1", "--protocol", "tcp", "--destination-port", "12:24", "--jump", "RETURN"},
	// 					},
	// 					fake_command_runner.CommandSpec{
	// 						Path: "/sbin/iptables",
	// 						Args: []string{"-w", "-I", "foo-bar-baz", "1", "--protocol", "tcp", "--destination-port", "64:942", "--jump", "RETURN"},
	// 					},
	// 				))
	// 			})
	// 		})
	// 	})

	// 	Describe("Protocol", func() {
	// 		Context("when tcp protocol is specified", func() {
	// 			It("passes tcp protocol to iptables", func() {
	// 				Expect(subject.PrependFilterRule(garden.NetOutRule{
	// 					Protocol: garden.ProtocolTCP,
	// 				})).To(Succeed())

	// 				Expect(fakeRunner).To(HaveExecutedSerially(fake_command_runner.CommandSpec{
	// 					Path: "/sbin/iptables",
	// 					Args: []string{"-w", "-I", "foo-bar-baz", "1", "--protocol", "tcp", "--jump", "RETURN"},
	// 				}))
	// 			})
	// 		})

	// 		Context("when udp protocol is specified", func() {
	// 			It("passes udp protocol to iptables", func() {
	// 				Expect(subject.PrependFilterRule(garden.NetOutRule{
	// 					Protocol: garden.ProtocolUDP,
	// 				})).To(Succeed())

	// 				Expect(fakeRunner).To(HaveExecutedSerially(fake_command_runner.CommandSpec{
	// 					Path: "/sbin/iptables",
	// 					Args: []string{"-w", "-I", "foo-bar-baz", "1", "--protocol", "udp", "--jump", "RETURN"},
	// 				}))
	// 			})
	// 		})

	// 		Context("when icmp protocol is specified", func() {
	// 			It("passes icmp protocol to iptables", func() {
	// 				Expect(subject.PrependFilterRule(garden.NetOutRule{
	// 					Protocol: garden.ProtocolICMP,
	// 				})).To(Succeed())

	// 				Expect(fakeRunner).To(HaveExecutedSerially(fake_command_runner.CommandSpec{
	// 					Path: "/sbin/iptables",
	// 					Args: []string{"-w", "-I", "foo-bar-baz", "1", "--protocol", "icmp", "--jump", "RETURN"},
	// 				}))
	// 			})

	// 			Context("when icmp type is specified", func() {
	// 				It("passes icmp protcol type to iptables", func() {
	// 					Expect(subject.PrependFilterRule(garden.NetOutRule{
	// 						Protocol: garden.ProtocolICMP,
	// 						ICMPs: &garden.ICMPControl{
	// 							Type: 99,
	// 						},
	// 					})).To(Succeed())

	// 					Expect(fakeRunner).To(HaveExecutedSerially(fake_command_runner.CommandSpec{
	// 						Path: "/sbin/iptables",
	// 						Args: []string{"-w", "-I", "foo-bar-baz", "1", "--protocol", "icmp", "--icmp-type", "99", "--jump", "RETURN"},
	// 					}))
	// 				})
	// 			})

	// 			Context("when icmp type and code are specified", func() {
	// 				It("passes icmp protcol type and code to iptables", func() {
	// 					Expect(subject.PrependFilterRule(garden.NetOutRule{
	// 						Protocol: garden.ProtocolICMP,
	// 						ICMPs: &garden.ICMPControl{
	// 							Type: 99,
	// 							Code: garden.ICMPControlCode(11),
	// 						},
	// 					})).To(Succeed())

	// 					Expect(fakeRunner).To(HaveExecutedSerially(fake_command_runner.CommandSpec{
	// 						Path: "/sbin/iptables",
	// 						Args: []string{"-w", "-I", "foo-bar-baz", "1", "--protocol", "icmp", "--icmp-type", "99/11", "--jump", "RETURN"},
	// 					}))
	// 				})
	// 			})
	// 		})
	// 	})

	// 	Describe("Log", func() {
	// 		It("redirects via the log chain if log is specified", func() {
	// 			Expect(subject.PrependFilterRule(garden.NetOutRule{
	// 				Log: true,
	// 			})).To(Succeed())

	// 			Expect(fakeRunner).To(HaveExecutedSerially(fake_command_runner.CommandSpec{
	// 				Path: "/sbin/iptables",
	// 				Args: []string{"-w", "-I", "foo-bar-baz", "1", "--protocol", "all", "--goto", "foo-bar-baz-log"},
	// 			}))
	// 		})
	// 	})

	// 	Context("when multiple port ranges and multiple networks are specified", func() {
	// 		It("opens the permutations of those port ranges and networks", func() {
	// 			Expect(subject.PrependFilterRule(garden.NetOutRule{
	// 				Protocol: garden.ProtocolTCP,
	// 				Networks: []garden.IPRange{
	// 					{
	// 						Start: net.ParseIP("1.2.3.4"),
	// 					},
	// 					{
	// 						Start: net.ParseIP("2.2.3.4"),
	// 						End:   net.ParseIP("2.2.3.9"),
	// 					},
	// 				},
	// 				Ports: []garden.PortRange{
	// 					{12, 24},
	// 					{64, 942},
	// 				},
	// 			})).To(Succeed())

	// 			Expect(fakeRunner.ExecutedCommands()).To(HaveLen(4))
	// 			Expect(fakeRunner).To(HaveExecutedSerially(
	// 				fake_command_runner.CommandSpec{
	// 					Path: "/sbin/iptables",
	// 					Args: []string{"-w", "-I", "foo-bar-baz", "1", "--protocol", "tcp", "--destination", "1.2.3.4", "--destination-port", "12:24", "--jump", "RETURN"},
	// 				},
	// 				fake_command_runner.CommandSpec{
	// 					Path: "/sbin/iptables",
	// 					Args: []string{"-w", "-I", "foo-bar-baz", "1", "--protocol", "tcp", "--destination", "1.2.3.4", "--destination-port", "64:942", "--jump", "RETURN"},
	// 				},
	// 				fake_command_runner.CommandSpec{
	// 					Path: "/sbin/iptables",
	// 					Args: []string{"-w", "-I", "foo-bar-baz", "1", "--protocol", "tcp", "-m", "iprange", "--dst-range", "2.2.3.4-2.2.3.9", "--destination-port", "12:24", "--jump", "RETURN"},
	// 				},
	// 				fake_command_runner.CommandSpec{
	// 					Path: "/sbin/iptables",
	// 					Args: []string{"-w", "-I", "foo-bar-baz", "1", "--protocol", "tcp", "-m", "iprange", "--dst-range", "2.2.3.4-2.2.3.9", "--destination-port", "64:942", "--jump", "RETURN"},
	// 				},
	// 			))
	// 		})
	// 	})

	// 	Context("when a portrange is specified for ProtocolALL", func() {
	// 		It("returns a nice error message", func() {
	// 			Expect(subject.PrependFilterRule(garden.NetOutRule{
	// 				Protocol: garden.ProtocolAll,
	// 				Ports:    []garden.PortRange{{Start: 1, End: 5}},
	// 			})).To(MatchError("Ports cannot be specified for Protocol ALL"))
	// 		})

	// 		It("does not run iptables", func() {
	// 			subject.PrependFilterRule(garden.NetOutRule{
	// 				Protocol: garden.ProtocolAll,
	// 				Ports:    []garden.PortRange{{Start: 1, End: 5}},
	// 			})

	// 			Expect(fakeRunner.ExecutedCommands()).To(HaveLen(0))
	// 		})
	// 	})

	// 	Context("when a portrange is specified for ProtocolICMP", func() {
	// 		It("returns a nice error message", func() {
	// 			Expect(subject.PrependFilterRule(garden.NetOutRule{
	// 				Protocol: garden.ProtocolICMP,
	// 				Ports:    []garden.PortRange{{Start: 1, End: 5}},
	// 			})).To(MatchError("Ports cannot be specified for Protocol ICMP"))
	// 		})

	// 		It("does not run iptables", func() {
	// 			subject.PrependFilterRule(garden.NetOutRule{
	// 				Protocol: garden.ProtocolICMP,
	// 				Ports:    []garden.PortRange{{Start: 1, End: 5}},
	// 			})

	// 			Expect(fakeRunner.ExecutedCommands()).To(HaveLen(0))
	// 		})
	// 	})

	// 	Context("when an invaild protocol is specified", func() {
	// 		It("returns an error", func() {
	// 			err := subject.PrependFilterRule(garden.NetOutRule{
	// 				Protocol: garden.Protocol(52),
	// 			})
	// 			Expect(err).To(HaveOccurred())
	// 			Expect(err).To(MatchError("invalid protocol: 52"))
	// 		})
	// 	})

	// 	Context("when the command returns an error", func() {
	// 		It("returns a wrapped error, including stderr", func() {
	// 			someError := errors.New("badly laid iptable")
	// 			fakeRunner.WhenRunning(
	// 				fake_command_runner.CommandSpec{Path: "/sbin/iptables"},
	// 				func(cmd *exec.Cmd) error {
	// 					cmd.Stderr.Write([]byte("stderr contents"))
	// 					return someError
	// 				},
	// 			)

	// 			Expect(subject.PrependFilterRule(garden.NetOutRule{})).To(MatchError("iptables: badly laid iptable, stderr contents"))
	// 		})
	// 	})
	// })
})
