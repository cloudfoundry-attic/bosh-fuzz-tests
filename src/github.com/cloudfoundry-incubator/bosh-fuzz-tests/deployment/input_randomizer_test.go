package deployment_test

import (
	bftconfig "github.com/cloudfoundry-incubator/bosh-fuzz-tests/config"

	. "github.com/cloudfoundry-incubator/bosh-fuzz-tests/deployment"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("InputRandomizer", func() {
	var (
		inputRandomizer InputRandomizer
	)

	It("generates specified number of inputs", func() {
		parameters := bftconfig.Parameters{
			Name:              []string{"foo"},
			Instances:         []int{5},
			AvailabilityZones: [][]string{[]string{"z1"}},
		}
		inputRandomizer = NewSeededInputRandomizer(parameters, 5, 64)

		inputs, err := inputRandomizer.Generate()
		Expect(err).ToNot(HaveOccurred())
		expectedInput := Input{
			Name:              "foo",
			Instances:         5,
			AvailabilityZones: []string{"z1"},
		}
		Expect(inputs).To(Equal([]Input{expectedInput, expectedInput, expectedInput, expectedInput, expectedInput}))
	})

	It("generates inputs with parameters shuffled", func() {
		parameters := bftconfig.Parameters{
			Name:              []string{"foo", "bar"},
			Instances:         []int{2, 4},
			AvailabilityZones: [][]string{[]string{"z1"}, []string{"z1", "z2"}},
		}
		inputRandomizer = NewSeededInputRandomizer(parameters, 2, 64)

		inputs, err := inputRandomizer.Generate()
		Expect(err).ToNot(HaveOccurred())

		Expect(inputs).To(Equal([]Input{
			{
				Name:              "bar",
				Instances:         2,
				AvailabilityZones: []string{"z1", "z2"},
			},
			{
				Name:              "foo",
				Instances:         4,
				AvailabilityZones: []string{"z1"},
			},
		}))
	})
})
