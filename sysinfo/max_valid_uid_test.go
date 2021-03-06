package sysinfo_test

import (
	"io/ioutil"
	"os"

	"github.com/cloudfoundry-incubator/guardian/sysinfo"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MaxValidUid", func() {
	var (
		tmpFiles []string
	)

	writeTmpFile := func(contents string) string {
		tmpFile, err := ioutil.TempFile("", "")
		Expect(err).ToNot(HaveOccurred())
		defer tmpFile.Close()

		_, err = tmpFile.WriteString(contents)
		Expect(err).ToNot(HaveOccurred())

		tmpFiles = append(tmpFiles, tmpFile.Name())
		return tmpFile.Name()
	}

	BeforeEach(func() {
		tmpFiles = make([]string, 0)
	})

	AfterEach(func() {
		for _, f := range tmpFiles {
			Expect(os.RemoveAll(f)).To(Succeed())
		}
	})

	Context("when the file has no entries", func() {
		It("should return 0", func() {
			Expect(sysinfo.IDMap(writeTmpFile("")).MaxValid()).To(Equal(0))
		})
	})

	Context("when the file has a single line", func() {
		It("returns ones less than the containerid column + the size column", func() {
			Expect(sysinfo.IDMap(writeTmpFile("12345 0 3")).MaxValid()).To(Equal(12347))
		})
	})

	Context("when the file has a multiple lines", func() {
		It("returns the largest value", func() {
			Expect(sysinfo.IDMap(writeTmpFile("12345 0 3\n44 0 1")).MaxValid()).To(Equal(12347))
			Expect(sysinfo.IDMap(writeTmpFile("44 0 1\n12345 0 1")).MaxValid()).To(Equal(12345))
		})
	})

	Context("when a line is invalid", func() {
		It("returns an error", func() {
			_, err := sysinfo.IDMap(writeTmpFile("cake")).MaxValid()
			Expect(err).To(MatchError(`expected integer while parsing line "cake"`))
		})
	})
})
