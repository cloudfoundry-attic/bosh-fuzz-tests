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
			NumberOfJobs:             []int{1, 2},
		}
		logger := boshlog.NewLogger(boshlog.LevelNone)
		inputRandomizer = NewSeededInputRandomizer(parameters, 2, 64, logger)

		inputs, err := inputRandomizer.Generate()
		Expect(err).ToNot(HaveOccurred())

		Expect(inputs).To(Equal([]Input{
			{
				Jobs: []Job{
					{
						Name:               "qNAwiIQ8el",
						Instances:          2,
						AvailabilityZones:  []string{"z1", "z2"},
						PersistentDiskType: "icN3O2GYdm",
						Network:            "default",
					},
					{
						Name:               "eagRjDTBs3",
						Instances:          4,
						AvailabilityZones:  []string{"z1", "z2"},
						PersistentDiskType: "rYND0xNg3R",
						Network:            "default",
					},
				},
				CloudConfig: CloudConfig{
					AvailabilityZones: []string{"z1", "z2"},
					PersistentDiskTypes: []DiskConfig{
						{Name: "icN3O2GYdm", Size: 200},
						{Name: "rYND0xNg3R", Size: 200},
					},
				},
			},
			{
				Jobs: []Job{
					{
						Name:               "mO5Qu",
						Instances:          4,
						AvailabilityZones:  []string{"z1", "z2"},
						PersistentDiskSize: 200,
						Network:            "default",
					},
				},
				CloudConfig: CloudConfig{
					AvailabilityZones: []string{"z1", "z2"},
				},
			},
		}))
	})
})
