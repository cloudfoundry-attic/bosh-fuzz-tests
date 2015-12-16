package parameter_test

import (
	"math/rand"

	bftinput "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"
	. "github.com/cloudfoundry-incubator/bosh-fuzz-tests/parameter"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Compilation", func() {
	var (
		compilation Parameter
	)

	BeforeEach(func() {
		compilation = NewCompilation([]int{1, 2, 3})
	})

	It("generates random number of compilation workers", func() {
		rand.Seed(42)

		input := bftinput.Input{}

		result := compilation.Apply(input, bftinput.Input{})

		Expect(result).To(Equal(bftinput.Input{
			CloudConfig: bftinput.CloudConfig{
				NumberOfCompilationWorkers: 3,
			},
		}))
	})
})
