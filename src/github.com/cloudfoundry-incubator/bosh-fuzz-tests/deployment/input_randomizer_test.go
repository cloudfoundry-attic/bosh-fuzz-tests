package deployment_test

import (
	bftconfig "github.com/cloudfoundry-incubator/bosh-fuzz-tests/config"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"

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
			NameLength:               []int{5, 10},
			Instances:                []int{2, 4},
			AvailabilityZones:        [][]string{[]string{"z1"}, []string{"z1", "z2"}},
			PersistentDiskDefinition: []string{"disk_pool", "disk_type", "persistent_disk_size"},
			PersistentDiskSize:       []int{0, 100, 200},
		}
		logger := boshlog.NewLogger(boshlog.LevelNone)
		inputRandomizer = NewSeededInputRandomizer(parameters, 3, 64, logger)

		inputs, err := inputRandomizer.Generate()
		Expect(err).ToNot(HaveOccurred())

		Expect(inputs).To(Equal([]Input{
			{
				Name:                     "iHoNAwiIQ8",
				Instances:                2,
				AvailabilityZones:        []string{"z1", "z2"},
				PersistentDiskDefinition: "disk_pool",
				PersistentDiskSize:       100,
			},
			{
				Name:                     "icN3O",
				Instances:                2,
				AvailabilityZones:        []string{"z1"},
				PersistentDiskDefinition: "disk_pool",
				PersistentDiskSize:       200,
			},
			{
				Name:                     "v6agR",
				Instances:                4,
				AvailabilityZones:        []string{"z1", "z2"},
				PersistentDiskDefinition: "disk_type",
				PersistentDiskSize:       200,
			},
		}))
	})
})
