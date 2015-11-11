package rundmc_test

import (
	"github.com/cloudfoundry-incubator/goci"
	"github.com/cloudfoundry-incubator/goci/specs"
	"github.com/cloudfoundry-incubator/guardian/gardener"
	"github.com/cloudfoundry-incubator/guardian/rundmc"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Bundle", func() {
	Context("when there is a network path", func() {
		It("adds the network path to the network namespace of the bundle", func() {
			base := goci.Bundle().WithNamespaces(specs.Namespace{Type: "network"})
			modifiedBundle := rundmc.BundleTemplate{Bndl: base}.Bundle(gardener.DesiredContainerSpec{
				NetworkPath: "/path/to/network",
			})

			Expect(modifiedBundle.RuntimeSpec.Linux.Namespaces).Should(ConsistOf(
				specs.Namespace{Type: specs.NetworkNamespace, Path: "/path/to/network"},
			))
		})

		It("does not modify the other fields", func() {
			base := goci.Bundle().WithProcess(goci.Process("potato"))
			modifiedBundle := rundmc.BundleTemplate{Bndl: base}.Bundle(gardener.DesiredContainerSpec{})
			Expect(modifiedBundle.Spec.Process.Args).Should(ConsistOf("potato"))
		})
	})

	Describe("Capabilities", func() {
		Context("when the container is not privileged", func() {
			It("should apply an array of capabilities", func() {
				base := goci.Bundle()
				modifiedBundle := rundmc.BundleTemplate{Bndl: base}.Bundle(gardener.DesiredContainerSpec{Privileged: false})

				capabilities := []string{
					"CAP_CHOWN",
					"CAP_DAC_OVERRIDE",
					"CAP_FSETID",
					"CAP_FOWNER",
					"CAP_MKNOD",
					"CAP_NET_RAW",
					"CAP_SETGID",
					"CAP_SETUID",
					"CAP_SETFCAP",
					"CAP_SETPCAP",
					"CAP_NET_BIND_SERVICE",
					"CAP_SYS_CHROOT",
					"CAP_KILL",
					"CAP_AUDIT_WRITE",
				}
				Expect(modifiedBundle.Spec.Linux.Capabilities).To(Equal(capabilities))
			})
		})

		Context("when the container is privileged", func() {
			It("should not apply capabilities", func() {
				base := goci.Bundle()
				modifiedBundle := rundmc.BundleTemplate{Bndl: base}.Bundle(gardener.DesiredContainerSpec{Privileged: true})
				Expect(modifiedBundle.Spec.Linux.Capabilities).To(BeNil())
			})
		})
	})
})
