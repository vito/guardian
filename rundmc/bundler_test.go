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

	Context("when privileged is true", func() {
		It("adds CAP_SYS_ADMIN capability to bundler", func() {
			base := goci.Bundle()
			template := rundmc.BundleTemplate{Bndl: base}
			modifiedBundle := template.Bundle(gardener.DesiredContainerSpec{Privileged: true})
			Expect(modifiedBundle.Spec.Linux.Capabilities).Should(ConsistOf("CAP_SYS_ADMIN"))
		})
	})

	Context("when privileged is false", func() {
		It("adds the right set of capabilities to bundler", func() {
			base := goci.Bundle()
			template := rundmc.BundleTemplate{Bndl: base}
			modifiedBundle := template.Bundle(gardener.DesiredContainerSpec{Privileged: false})

			Expect(modifiedBundle.Spec.Linux.Capabilities).ShouldNot(ContainElement("CAP_SYS_ADMIN"))

			Expect(modifiedBundle.Spec.Linux.Capabilities).Should(ContainElement("CAP_CHOWN"))
			Expect(modifiedBundle.Spec.Linux.Capabilities).Should(ContainElement("CAP_DAC_OVERRIDE"))
			Expect(modifiedBundle.Spec.Linux.Capabilities).Should(ContainElement("CAP_FSETID"))
			Expect(modifiedBundle.Spec.Linux.Capabilities).Should(ContainElement("CAP_FOWNER"))
			Expect(modifiedBundle.Spec.Linux.Capabilities).Should(ContainElement("CAP_MKNOD"))
			Expect(modifiedBundle.Spec.Linux.Capabilities).Should(ContainElement("CAP_NET_RAW"))
			Expect(modifiedBundle.Spec.Linux.Capabilities).Should(ContainElement("CAP_SETGID"))
			Expect(modifiedBundle.Spec.Linux.Capabilities).Should(ContainElement("CAP_SETUID"))
			Expect(modifiedBundle.Spec.Linux.Capabilities).Should(ContainElement("CAP_SETFCAP"))
			Expect(modifiedBundle.Spec.Linux.Capabilities).Should(ContainElement("CAP_SETPCAP"))
			Expect(modifiedBundle.Spec.Linux.Capabilities).Should(ContainElement("CAP_NET_BIND_SERVICE"))
			Expect(modifiedBundle.Spec.Linux.Capabilities).Should(ContainElement("CAP_SYS_CHROOT"))
			Expect(modifiedBundle.Spec.Linux.Capabilities).Should(ContainElement("CAP_KILL"))
			Expect(modifiedBundle.Spec.Linux.Capabilities).Should(ContainElement("CAP_AUDIT_WRITE"))
		})
	})
})
