package rundmc_test

import (
	"io"
	"io/ioutil"
	"os/exec"
	"path/filepath"

	"github.com/cloudfoundry-incubator/guardian/rundmc"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("PivotRootDirMaker", func() {
	It("makes a sensible subdirectory in the directory we asked it to", func() {
		tmpDir, err := ioutil.TempDir("", "make-tester")
		Expect(err).NotTo(HaveOccurred())

		maker := rundmc.RootFsPreparer{}
		maker.Make(nil, false, tmpDir)
		Expect(filepath.Join(tmpDir, ".garden")).To(BeADirectory())
	})

	Context("when I'm building a privileged container", func() {
		It("should make the directory be owned by host root", func() {
			tmpDir, err := ioutil.TempDir("", "make-tester")
			Expect(err).NotTo(HaveOccurred())

			maker := rundmc.RootFsPreparer{
				OwnerUID: 12,
			}
			Expect(maker.Make(nil, true, tmpDir)).To(Succeed())

			buf := gbytes.NewBuffer()
			session, err := gexec.Start(
				exec.Command("stat", "-f", "%g", filepath.Join(tmpDir, ".garden")),
				io.MultiWriter(buf, GinkgoWriter), GinkgoWriter)

			Expect(err).NotTo(HaveOccurred())
			Eventually(session).Should(gexec.Exit(0))
			Expect(buf).To(gbytes.Say("12"))
		})
	})

	Context("when I'm building an unprivileged container", func() {
		It("should make the directory be owned by container root", func() {
		})
	})
})
