package rundmc_test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/cf-guardian/specs"
	"github.com/cloudfoundry-incubator/guardian/rundmc"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Depot", func() {
	var (
		tmpDir string
		depot  *rundmc.DirectoryDepot
	)

	BeforeEach(func() {
		var err error

		tmpDir, err = ioutil.TempDir("", "depot-test")
		Expect(err).NotTo(HaveOccurred())

		depot = &rundmc.DirectoryDepot{
			Dir: tmpDir,
		}
	})

	Describe("lookup", func() {
		Context("when a subdirectory with the given name does not exist", func() {
			It("returns an ErrDoesNotExist", func() {
				_, err := depot.Lookup("potato")
				Expect(err).To(MatchError(rundmc.ErrDoesNotExist))
			})
		})

		Context("when a subdirectory with the given name exists", func() {
			It("returns the absolute path of the directory", func() {
				os.Mkdir(filepath.Join(tmpDir, "potato"), 0700)
				Expect(depot.Lookup("potato")).To(Equal(filepath.Join(tmpDir, "potato")))
			})
		})
	})

	Describe("create", func() {
		BeforeEach(func() {
			Expect(depot.Create("aardvaark")).To(Succeed())
		})

		It("should create a directory", func() {
			Expect(filepath.Join(tmpDir, "aardvaark")).To(BeADirectory())
		})

		Describe("the container directory", func() {
			It("should contain a config.json", func() {
				Expect(filepath.Join(tmpDir, "aardvaark", "config.json")).To(BeARegularFile())
			})

			It("should have a config.json with a process which echos 'Pid 1 Running'", func() {
				file, err := os.Open(filepath.Join(tmpDir, "aardvaark", "config.json"))
				Expect(err).NotTo(HaveOccurred())

				var target specs.Spec
				Expect(json.NewDecoder(file).Decode(&target)).To(Succeed())

				Expect(target.Process).To(Equal(specs.Process{
					Terminal: true,
					Args: []string{
						"/bin/sh", "-c", `echo "Pid 1 Running"; read`,
					},
				}))
			})

			It("should have a config.json which specifies the spec version", func() {
				file, err := os.Open(filepath.Join(tmpDir, "aardvaark", "config.json"))
				Expect(err).NotTo(HaveOccurred())

				var target specs.Spec
				Expect(json.NewDecoder(file).Decode(&target)).To(Succeed())

				Expect(target.Version).To(Equal("pre-draft"))
			})
		})
	})
})
