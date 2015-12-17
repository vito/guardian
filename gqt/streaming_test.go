package gqt_test

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"

	"github.com/cloudfoundry-incubator/garden"
	"github.com/cloudfoundry-incubator/guardian/gqt/runner"
	. "github.com/cloudfoundry-incubator/guardian/matchers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = FDescribe("Streaming", func() {
	var (
		client    *runner.RunningGarden
		container garden.Container
	)

	BeforeEach(func() {
		var err error

		client = startGarden()

		container, err = client.Create(garden.ContainerSpec{})
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		Expect(client.DestroyAndStop()).To(Succeed())
	})

	Describe("StreamIn", func() {
		var tarStream io.Reader

		BeforeEach(func() {
			tgzPath := os.Getenv("GARDEN_DORA_PATH")
			if tgzPath == "" {
				Skip("`GARDEN_DORA_PATH` environment variable was not found")
			}

			f, err := os.Open(tgzPath)
			Expect(err).NotTo(HaveOccurred())
			defer f.Close()

			tgz, err := gzip.NewReader(f)
			Expect(err).NotTo(HaveOccurred())
			defer tgz.Close()

			tarStream = tar.NewReader(tgz)
		})

		It("should stream in the files", func() {
			Expect(container.StreamIn(garden.StreamInSpec{
				Path:      "/root/test",
				User:      "root",
				TarStream: tarStream,
			})).To(Succeed())

			Expect(container).To(HaveFile("/root/test/staging_info.yml"))
			Expect(container).To(HaveFile("/root/test/app"))
			Expect(container).To(HaveFile("/root/test/logs"))
			Expect(container).To(HaveFile("/root/test/tmp"))
		})
	})
})
