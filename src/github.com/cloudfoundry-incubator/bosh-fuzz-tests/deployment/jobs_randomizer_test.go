package deployment_test

import (
	"math/rand"

	bftconfig "github.com/cloudfoundry-incubator/bosh-fuzz-tests/config"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"

	. "github.com/cloudfoundry-incubator/bosh-fuzz-tests/deployment"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("JobsRandomizer", func() {
	var (
		jobsRandomizer JobsRandomizer
	)

	It("generates extra input for migrated jobs", func() {
		parameters := bftconfig.Parameters{
			NameLength:               []int{5, 10},
			Instances:                []int{2, 4},
			AvailabilityZones:        [][]string{[]string{"z1"}, []string{"z1", "z2"}},
			PersistentDiskDefinition: []string{"disk_pool", "disk_type", "persistent_disk_size"},
			PersistentDiskSize:       []int{0, 100, 200},
			NumberOfJobs:             []int{1, 2},
			MigratedFromCount:        []int{0, 2},
		}
		logger := boshlog.NewLogger(boshlog.LevelNone)
		rand.Seed(64)
		nameGenerator := NewNameGenerator()
		jobsRandomizer = NewJobsRandomizer(parameters, 2, nameGenerator, logger)

		inputs, err := jobsRandomizer.Generate()
		Expect(err).ToNot(HaveOccurred())

		Expect(inputs).To(Equal([]Input{
			{
				Jobs: []Job{
					{
						Name:              "joNAw",
						Instances:         4,
						AvailabilityZones: []string{"z1"},
					},
					{
						Name:              "gQ8el",
						Instances:         4,
						AvailabilityZones: []string{"z1", "z2"},
					},
				},
				CloudConfig: CloudConfig{
					AvailabilityZones: []string{"z1", "z2"},
				},
			},
			{
				Jobs: []Job{
					{
						Name:               "rU3YND0xNg",
						Instances:          4,
						AvailabilityZones:  []string{"z1"},
						PersistentDiskPool: "gBUnQKBYoE",
					},
					{
						Name:               "pRWDsiO5Qu",
						Instances:          4,
						AvailabilityZones:  []string{"z1"},
						PersistentDiskPool: "a5gmsYqE7Y",
					},
				},
				CloudConfig: CloudConfig{
					AvailabilityZones: []string{"z1"},
					PersistentDiskPools: []DiskConfig{
						{Name: "gBUnQKBYoE", Size: 100},
						{Name: "a5gmsYqE7Y", Size: 100},
					},
				},
			},
			{
				Jobs: []Job{
					{
						Name:               "joNAw",
						Instances:          4,
						AvailabilityZones:  []string{"z1"},
						PersistentDiskPool: "eagRjDTBs3",
						MigratedFrom: []MigratedFromConfig{
							{Name: "rU3YND0xNg"},
							{Name: "pRWDsiO5Qu"},
						},
					},
				},
				CloudConfig: CloudConfig{
					AvailabilityZones: []string{"z1"},
					PersistentDiskPools: []DiskConfig{
						{Name: "eagRjDTBs3", Size: 100},
					},
				},
			},
		}))
	})

	It("when migrated job does not have az it sets random az in migrated_from", func() {
		parameters := bftconfig.Parameters{
			NameLength:               []int{5},
			Instances:                []int{2},
			AvailabilityZones:        [][]string{[]string{"z1"}, nil},
			PersistentDiskDefinition: []string{"persistent_disk_size"},
			PersistentDiskSize:       []int{0},
			NumberOfJobs:             []int{1},
			MigratedFromCount:        []int{1},
		}
		logger := boshlog.NewLogger(boshlog.LevelNone)
		rand.Seed(64)
		nameGenerator := NewNameGenerator()
		jobsRandomizer = NewJobsRandomizer(parameters, 1, nameGenerator, logger)

		inputs, err := jobsRandomizer.Generate()
		Expect(err).ToNot(HaveOccurred())

		Expect(inputs).To(Equal([]Input{
			{
				Jobs: []Job{
					{
						Name:      "vgrKicN3O2",
						Instances: 2,
					},
				},
			},
			{
				Jobs: []Job{
					{
						Name:              "joNAw",
						Instances:         2,
						AvailabilityZones: []string{"z1"},
						MigratedFrom: []MigratedFromConfig{
							{Name: "vgrKicN3O2", AvailabilityZone: "z1"},
						},
					},
				},
				CloudConfig: CloudConfig{
					AvailabilityZones: []string{"z1"},
				},
			},
		}))
	})
})
