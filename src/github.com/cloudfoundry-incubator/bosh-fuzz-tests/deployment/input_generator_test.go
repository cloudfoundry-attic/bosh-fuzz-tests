package deployment_test

import (
	"math/rand"

	fakebftparam "github.com/cloudfoundry-incubator/bosh-fuzz-tests/parameter/fakes"

	bftconfig "github.com/cloudfoundry-incubator/bosh-fuzz-tests/config"
	bftinput "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"
	bftnamegen "github.com/cloudfoundry-incubator/bosh-fuzz-tests/name_generator"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"

	. "github.com/cloudfoundry-incubator/bosh-fuzz-tests/deployment"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("InputGenerator", func() {
	var (
		inputGenerator InputGenerator
	)

	It("generates requested number of inputs", func() {
		parameters := bftconfig.Parameters{
			NameLength:               []int{5},
			Instances:                []int{2},
			AvailabilityZones:        [][]string{[]string{"z1"}, []string{"z1", "z2"}},
			PersistentDiskDefinition: []string{"disk_pool"},
			PersistentDiskSize:       []int{100},
			NumberOfJobs:             []int{2},
			MigratedFromCount:        []int{0},
			VmTypeDefinition:         []string{"vm_type"},
			StemcellDefinition:       []string{"os"},
		}
		logger := boshlog.NewLogger(boshlog.LevelNone)
		rand.Seed(64)
		nameGenerator := bftnamegen.NewNameGenerator()
		fakeParameterProvider := fakebftparam.NewFakeParameterProvider()
		inputGenerator = NewInputGenerator(parameters, fakeParameterProvider, 2, nameGenerator, logger)

		inputs, err := inputGenerator.Generate()
		Expect(err).ToNot(HaveOccurred())

		Expect(inputs).To(Equal([]bftinput.Input{
			{
				Jobs: []bftinput.Job{
					{
						Name:               "joNAw",
						Instances:          2,
						AvailabilityZones:  []string{"z1"},
						PersistentDiskPool: "fake-persistent-disk",
						VmType:             "fake-vm-type",
					},
					{
						Name:               "gQ8el",
						Instances:          2,
						AvailabilityZones:  []string{"z1"},
						PersistentDiskPool: "fake-persistent-disk",
						VmType:             "fake-vm-type",
					},
				},
				CloudConfig: bftinput.CloudConfig{
					AvailabilityZones: []string{"z1"},
					PersistentDiskPools: []bftinput.DiskConfig{
						{
							Name: "fake-persistent-disk",
							Size: 1,
						},
					},
					VmTypes: []bftinput.VmTypeConfig{
						{Name: "fake-vm-type"},
					},
				},
				Stemcells: []bftinput.StemcellConfig{
					{Name: "fake-stemcell"},
				},
			},
			{
				Jobs: []bftinput.Job{
					{
						Name:               "joNAw",
						Instances:          2,
						AvailabilityZones:  []string{"z1"},
						PersistentDiskPool: "fake-persistent-disk",
						VmType:             "fake-vm-type",
					},
					{
						Name:               "gQ8el",
						Instances:          2,
						AvailabilityZones:  []string{"z1"},
						PersistentDiskPool: "fake-persistent-disk",
						VmType:             "fake-vm-type",
					},
				},
				CloudConfig: bftinput.CloudConfig{
					AvailabilityZones: []string{"z1"},
					PersistentDiskPools: []bftinput.DiskConfig{
						{Name: "fake-persistent-disk", Size: 1},
					},
					VmTypes: []bftinput.VmTypeConfig{
						{Name: "fake-vm-type"},
					},
				},
				Stemcells: []bftinput.StemcellConfig{
					{Name: "fake-stemcell"},
				},
			},
		}))
	})

	It("when migrated job does not have az it sets random az in migrated_from", func() {
		parameters := bftconfig.Parameters{
			NameLength:               []int{5},
			Instances:                []int{2},
			AvailabilityZones:        [][]string{[]string{"z1"}},
			PersistentDiskDefinition: []string{"persistent_disk_size"},
			PersistentDiskSize:       []int{0},
			NumberOfJobs:             []int{1},
			MigratedFromCount:        []int{1},
			VmTypeDefinition:         []string{"vm_type"},
			StemcellDefinition:       []string{"name"},
		}
		logger := boshlog.NewLogger(boshlog.LevelNone)
		rand.Seed(64)
		nameGenerator := bftnamegen.NewNameGenerator()
		fakeParameterProvider := fakebftparam.NewFakeParameterProvider()
		inputGenerator = NewInputGenerator(parameters, fakeParameterProvider, 1, nameGenerator, logger)

		inputs, err := inputGenerator.Generate()
		Expect(err).ToNot(HaveOccurred())

		Expect(inputs).To(Equal([]bftinput.Input{
			{
				Jobs: []bftinput.Job{
					{
						Name:               "oelgrKicN3",
						Instances:          2,
						VmType:             "fake-vm-type",
						AvailabilityZones:  []string{"z1"},
						PersistentDiskPool: "fake-persistent-disk",
					},
				},
				CloudConfig: bftinput.CloudConfig{
					AvailabilityZones: []string{"z1"},
					PersistentDiskPools: []bftinput.DiskConfig{
						{
							Name: "fake-persistent-disk",
							Size: 1,
						},
					},
					VmTypes: []bftinput.VmTypeConfig{
						{Name: "fake-vm-type"},
					},
				},
				Stemcells: []bftinput.StemcellConfig{
					{Name: "fake-stemcell"},
				},
			},
			{
				Jobs: []bftinput.Job{
					{
						Name:               "joNAw",
						Instances:          2,
						AvailabilityZones:  []string{"z1"},
						VmType:             "fake-vm-type",
						PersistentDiskPool: "fake-persistent-disk",
						MigratedFrom: []bftinput.MigratedFromConfig{
							{Name: "oelgrKicN3"},
						},
					},
				},
				CloudConfig: bftinput.CloudConfig{
					AvailabilityZones: []string{"z1"},
					PersistentDiskPools: []bftinput.DiskConfig{
						{
							Name: "fake-persistent-disk",
							Size: 1,
						},
					},
					VmTypes: []bftinput.VmTypeConfig{
						{Name: "fake-vm-type"},
					},
				},
				Stemcells: []bftinput.StemcellConfig{
					{Name: "fake-stemcell"},
				},
			},
		}))
	})
})
