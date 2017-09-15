package config_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cloudfoundry-incubator/bosh-load-tests/config"
	"github.com/cloudfoundry/bosh-utils/system/fakes"
)

var _ = Describe("Config", func() {
	Describe("Load", func() {
		var c *config.Config

		BeforeEach(func() {
			fs := fakes.NewFakeFileSystem()
			c = config.NewConfig(fs)

			err := fs.WriteFileString("/some-config-file", `{"ruby_version": "2.3.1"}`)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when RUBY_VERSION is an env var", func() {
			AfterEach(func() {
				err := os.Unsetenv("RUBY_VERSION")
				Expect(err).NotTo(HaveOccurred())
			})

			It("overrides the file", func() {
				err := os.Setenv("RUBY_VERSION", "2.4.2")
				Expect(err).NotTo(HaveOccurred())

				err = c.Load("/some-config-file")
				Expect(err).NotTo(HaveOccurred())

				Expect(c.RubyVersion).To(Equal("2.4.2"))
			})
		})

		Context("when RUBY_VERSION is not an env var", func() {
			It("does not override the file", func() {
				err := c.Load("/some-config-file")
				Expect(err).NotTo(HaveOccurred())

				Expect(c.RubyVersion).To(Equal("2.3.1"))
			})
		})
	})
})
