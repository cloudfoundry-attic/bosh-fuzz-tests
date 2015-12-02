package name_generator_test

import (
	"math/rand"

	. "github.com/cloudfoundry-incubator/bosh-fuzz-tests/name_generator"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NameGenerator", func() {
	var (
		nameGenerator NameGenerator
	)

	BeforeEach(func() {
		rand.Seed(5)
		nameGenerator = NewNameGenerator()
	})

	It("generates name of specified length", func() {
		name := nameGenerator.Generate(5)
		Expect(len(name)).To(Equal(5))
		Expect(name).To(Equal("qgvMT"))
	})
})
