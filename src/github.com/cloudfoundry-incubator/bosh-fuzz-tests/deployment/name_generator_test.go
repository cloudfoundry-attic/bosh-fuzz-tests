package deployment_test

import (
	. "github.com/cloudfoundry-incubator/bosh-fuzz-tests/deployment"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NameGenerator", func() {
	var (
		nameGenerator NameGenerator
	)

	BeforeEach(func() {
		nameGenerator = NewSeededNameGenerator(5)
	})

	It("generates name of specified length", func() {
		name := nameGenerator.Generate(5)
		Expect(len(name)).To(Equal(5))
		Expect(name).To(Equal("qgvMT"))
	})
})
