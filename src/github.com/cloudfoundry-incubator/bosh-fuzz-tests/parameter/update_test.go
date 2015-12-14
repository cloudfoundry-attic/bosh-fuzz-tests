package parameter_test

import (
	"math/rand"

	. "github.com/cloudfoundry-incubator/bosh-fuzz-tests/parameter"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	bftinput "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"
)

var _ = Describe("Update", func() {
	var (
		update Parameter
	)

	BeforeEach(func() {
		update = NewUpdate([]int{1, 5}, []int{1, 3}, []string{"not_specified", "true", "false"})
	})

	Context("when serial is not not_specified", func() {
		It("adds random settings for update", func() {
			rand.Seed(64)
			input := bftinput.Input{}
			result := update.Apply(input)

			Expect(result).To(Equal(bftinput.Input{
				Update: bftinput.UpdateConfig{
					Canaries:    5,
					MaxInFlight: 1,
					Serial:      "false",
				},
			}))
		})
	})

	Context("when serial is not_specified", func() {
		It("does not show serial", func() {
			rand.Seed(2)
			input := bftinput.Input{}
			result := update.Apply(input)

			Expect(result).To(Equal(bftinput.Input{
				Update: bftinput.UpdateConfig{
					Canaries:    1,
					MaxInFlight: 1,
					Serial:      "not_specified",
				},
			}))
		})
	})
})
