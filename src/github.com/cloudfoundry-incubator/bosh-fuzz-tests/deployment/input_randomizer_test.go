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

	It("generates inputs with parameters shuffled", func() {
		parameters := bftconfig.Parameters{
			NameLength:        []int{5, 10},
			Instances:         []int{2, 4},
			AvailabilityZones: [][]string{[]string{"z1"}, []string{"z1", "z2"}},
		}
		inputRandomizer = NewSeededInputRandomizer(parameters, 3, 64)

		inputs, err := inputRandomizer.Generate()
		Expect(err).ToNot(HaveOccurred())

		Expect(inputs).To(Equal([]Input{
			{
				Name:              "izSREqw4Qe",
				Instances:         2,
				AvailabilityZones: []string{"z1", "z2"},
			},
			{
				Name:              "jaWC7",
				Instances:         4,
				AvailabilityZones: []string{"z1"},
			},
			{
				Name:              "c_VC5",
				Instances:         2,
				AvailabilityZones: []string{"z1"},
			},
		}))
	})
})
