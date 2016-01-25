package rundmc_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/cloudfoundry-incubator/goci"
	"github.com/cloudfoundry-incubator/goci/specs"
	"github.com/cloudfoundry-incubator/guardian/gardener"
	"github.com/cloudfoundry-incubator/guardian/rundmc"
	"github.com/cloudfoundry-incubator/guardian/rundmc/fakes"
)

var _ = Describe("BaseTemplateRule", func() {
	var (
		privilegeBndl, unprivilegeBndl *goci.Bndl

		rule rundmc.BaseTemplateRule
	)

	BeforeEach(func() {
		privilegeBndl = goci.Bndl{}.WithNamespace(goci.NetworkNamespace)
		unprivilegeBndl = goci.Bndl{}.WithNamespace(goci.UserNamespace)

		rule = rundmc.BaseTemplateRule{
			PrivilegedBase:   privilegeBndl,
			UnprivilegedBase: unprivilegeBndl,
		}
	})

	Context("when it is privileged", func() {
		It("should use the correct base", func() {
			retBndl := rule.Apply(nil, gardener.DesiredContainerSpec{
				Privileged: true,
			})

			Expect(retBndl).To(Equal(privilegeBndl))
		})
	})

	Context("when it is not privileged", func() {
		It("should use the correct base", func() {
			retBndl := rule.Apply(nil, gardener.DesiredContainerSpec{
				Privileged: false,
			})

			Expect(retBndl).To(Equal(unprivilegeBndl))
		})
	})
})

var _ = Describe("RootFSRule", func() {
	var (
		fakeMkdirChowner *fakes.FakeMkdirChowner
		rule             rundmc.RootFSRule
	)

	BeforeEach(func() {
		fakeMkdirChowner = new(fakes.FakeMkdirChowner)
		rule = rundmc.RootFSRule{
			ContainerRootUID: 999,
			ContainerRootGID: 888,

			MkdirChowner: fakeMkdirChowner,
		}
	})

	It("applies the rootfs to the passed bundle", func() {
		newBndl := rule.Apply(goci.Bundle(), gardener.DesiredContainerSpec{
			RootFSPath: "/path/to/banana/rootfs",
		})

		Expect(newBndl.Spec.Root.Path).To(Equal("/path/to/banana/rootfs"))
	})

	// this is a workaround for our current aufs code not properly changing the
	// ownership of / to container-root. Without this step runC is unable to
	// pivot root in user-namespaced containers.
	Describe("creating the .pivot_root directory", func() {
		It("pre-creates the /.pivot_root directory with the correct ownership", func() {
			rule.Apply(goci.Bundle(), gardener.DesiredContainerSpec{
				RootFSPath: "/path/to/banana",
			})

			Expect(fakeMkdirChowner.MkdirChownCallCount()).To(Equal(1))
			path, perms, uid, gid := fakeMkdirChowner.MkdirChownArgsForCall(0)
			Expect(path).To(Equal("/path/to/banana/.pivot_root"))
			Expect(perms).To(Equal(os.FileMode(0700)))
			Expect(uid).To(Equal(999))
			Expect(gid).To(Equal(888))
		})

		PIt("what do we about errors", func() {})
	})
})

var _ = Describe("NetworkHookRule", func() {
	DescribeTable("the envirionment should contain", func(envVar string) {
		rule := rundmc.NetworkHookRule{LogFilePattern: "/path/to/%s.log"}

		newBndl := rule.Apply(goci.Bundle(), gardener.DesiredContainerSpec{
			Handle: "fred",
		})

		Expect(newBndl.RuntimeSpec.Hooks.Prestart[0].Env).To(
			ContainElement(envVar),
		)
	},
		Entry("the GARDEN_LOG_FILE path", "GARDEN_LOG_FILE=/path/to/fred.log"),
		Entry("a sensible PATH", "PATH="+os.Getenv("PATH")),
	)

	It("add the hook to the pre-start hooks of the passed bundle", func() {
		newBndl := rundmc.NetworkHookRule{}.Apply(goci.Bundle(), gardener.DesiredContainerSpec{
			NetworkHook: gardener.Hook{
				Path: "/path/to/bananas/network",
				Args: []string{"arg", "barg"},
			},
		})

		Expect(pathAndArgsOf(newBndl.RuntimeSpec.Hooks.Prestart)).To(ContainElement(PathAndArgs{
			Path: "/path/to/bananas/network",
			Args: []string{"arg", "barg"},
		}))
	})
})

func pathAndArgsOf(a []specs.Hook) (b []PathAndArgs) {
	for _, h := range a {
		b = append(b, PathAndArgs{h.Path, h.Args})
	}

	return
}

type PathAndArgs struct {
	Path string
	Args []string
}
