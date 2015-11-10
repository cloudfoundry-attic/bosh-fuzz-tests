package manifest_test

import (
	bftconfig "github.com/cloudfoundry-incubator/bosh-fuzz-tests/config"

	. "github.com/cloudfoundry-incubator/bosh-fuzz-tests/manifest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("InputRandomizer", func() {
	var (
		inputRandomizer InputRandomizer
	)

	BeforeEach(func() {
		inputRandomizer = NewInputRandomizer()
	})

	It("generates specified number of inputs", func() {
		parameters := bftconfig.Parameters{
			Name:              []string{"foo"},
			Instances:         []int{5},
			AvailabilityZones: [][]string{[]string{"z1"}},
		}

		inputs, err := inputRandomizer.Generate(parameters, 5)
		Expect(err).ToNot(HaveOccurred())
		expectedInput := Input{
			Name:              "foo",
			Instances:         5,
			AvailabilityZones: []string{"z1"},
		}
		Expect(inputs).To(Equal([]Input{expectedInput, expectedInput, expectedInput, expectedInput, expectedInput}))
	})
})
