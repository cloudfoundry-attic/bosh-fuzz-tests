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
			VmTypeDefinition:         []string{"vm_type", "resource_pool"},
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
						Instances:         2,
						AvailabilityZones: []string{"z1"},
						ResourcePool:      "h3O2GYdmz6",
					},
					{
						Name:               "gQ8el",
						Instances:          2,
						AvailabilityZones:  []string{"z1", "z2"},
						PersistentDiskPool: "pTBs3VXU3Y",
						ResourcePool:       "xD0xNg3RWD",
					},
				},
				CloudConfig: CloudConfig{
					AvailabilityZones: []string{"z1", "z2"},
					PersistentDiskPools: []DiskConfig{
						{Name: "pTBs3VXU3Y", Size: 100},
					},
					ResourcePools: []VmTypeConfig{
						{Name: "h3O2GYdmz6"},
						{Name: "xD0xNg3RWD"},
					},
				},
			},
			{
				Jobs: []Job{
					{
						Name:              "joNAw",
						Instances:         2,
						AvailabilityZones: []string{"z1"},
						ResourcePool:      "fqDqBUnQKB",
					},
				},
				CloudConfig: CloudConfig{
					AvailabilityZones: []string{"z1"},
					ResourcePools: []VmTypeConfig{
						{Name: "fqDqBUnQKB"},
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
			VmTypeDefinition:         []string{"vm_type"},
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
						Name:      "qdmz6agRjD",
						Instances: 2,
						VmType:    "rU3YND0xNg",
					},
				},
				CloudConfig: CloudConfig{
					VmTypes: []VmTypeConfig{
						{Name: "rU3YND0xNg"},
					},
				},
			},
			{
				Jobs: []Job{
					{
						Name:              "joNAw",
						Instances:         2,
						AvailabilityZones: []string{"z1"},
						VmType:            "vgrKicN3O2",
						MigratedFrom: []MigratedFromConfig{
							{Name: "qdmz6agRjD", AvailabilityZone: "z1"},
						},
					},
				},
				CloudConfig: CloudConfig{
					AvailabilityZones: []string{"z1"},
					VmTypes: []VmTypeConfig{
						{Name: "vgrKicN3O2"},
					},
				},
			},
		}))
	})
})
