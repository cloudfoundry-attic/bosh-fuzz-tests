package deployment_test

import (
	"math/rand"

	fakebftdecider "github.com/cloudfoundry-incubator/bosh-fuzz-tests/decider/fakes"
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
		inputGenerator        InputGenerator
		parameters            bftconfig.Parameters
		logger                boshlog.Logger
		nameGenerator         bftnamegen.NameGenerator
		fakeParameterProvider *fakebftparam.FakeParameterProvider
		decider               *fakebftdecider.FakeDecider
	)

	BeforeEach(func() {
		parameters = bftconfig.Parameters{
			NameLength:                 []int{5},
			Instances:                  []int{2},
			AvailabilityZones:          [][]string{[]string{"z1"}},
			PersistentDiskDefinition:   []string{"persistent_disk_size"},
			PersistentDiskSize:         []int{0},
			NumberOfJobs:               []int{1},
			MigratedFromCount:          []int{1},
			VmTypeDefinition:           []string{"vm_type"},
			StemcellDefinition:         []string{"name"},
			Templates:                  [][]string{[]string{"simple"}},
			NumberOfCompilationWorkers: []int{3},
			Canaries:                   []int{5},
			MaxInFlight:                []int{3},
			Serial:                     []string{"true"},
		}
		logger = boshlog.NewLogger(boshlog.LevelNone)
		nameGenerator = bftnamegen.NewNameGenerator()
		fakeParameterProvider = fakebftparam.NewFakeParameterProvider()
		decider = &fakebftdecider.FakeDecider{}
	})

	It("generates requested number of inputs", func() {
		parameters = bftconfig.Parameters{
			NameLength:                 []int{5},
			Instances:                  []int{2},
			AvailabilityZones:          [][]string{[]string{"z1"}, []string{"z1", "z2"}},
			PersistentDiskDefinition:   []string{"disk_pool"},
			PersistentDiskSize:         []int{100},
			NumberOfJobs:               []int{2},
			MigratedFromCount:          []int{0},
			VmTypeDefinition:           []string{"vm_type"},
			StemcellDefinition:         []string{"os"},
			Templates:                  [][]string{[]string{"simple"}},
			NumberOfCompilationWorkers: []int{3},
		}
		rand.Seed(64)
		inputGenerator = NewInputGenerator(parameters, fakeParameterProvider, 2, nameGenerator, decider, logger)

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
						Networks: []bftinput.JobNetworkConfig{
							{Name: "foo-network"},
						},
						Templates: []bftinput.Template{
							{Name: "simple"},
						},
					},
					{
						Name:               "gQ8el",
						Instances:          2,
						AvailabilityZones:  []string{"z1"},
						PersistentDiskPool: "fake-persistent-disk",
						VmType:             "fake-vm-type",
						Networks: []bftinput.JobNetworkConfig{
							{Name: "foo-network"},
						},
						Templates: []bftinput.Template{
							{Name: "simple"},
						},
					},
				},
				Update: bftinput.UpdateConfig{
					Canaries:    3,
					MaxInFlight: 5,
					Serial:      "true",
				},
				CloudConfig: bftinput.CloudConfig{
					Networks: []bftinput.NetworkConfig{
						{Name: "foo-network"},
					},
					AvailabilityZones: []bftinput.AvailabilityZone{
						{Name: "z1"},
					},
					PersistentDiskPools: []bftinput.DiskConfig{
						{
							Name: "fake-persistent-disk",
							Size: 1,
						},
					},
					VmTypes: []bftinput.VmTypeConfig{
						{Name: "fake-vm-type"},
					},
					NumberOfCompilationWorkers: 3,
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
						Networks: []bftinput.JobNetworkConfig{
							{Name: "foo-network"},
						},
						Templates: []bftinput.Template{
							{Name: "simple"},
						},
					},
					{
						Name:               "gQ8el",
						Instances:          2,
						AvailabilityZones:  []string{"z1"},
						PersistentDiskPool: "fake-persistent-disk",
						VmType:             "fake-vm-type",
						Networks: []bftinput.JobNetworkConfig{
							{Name: "foo-network"},
						},
						Templates: []bftinput.Template{
							{Name: "simple"},
						},
					},
				},
				Update: bftinput.UpdateConfig{
					Canaries:    3,
					MaxInFlight: 5,
					Serial:      "true",
				},
				CloudConfig: bftinput.CloudConfig{
					NumberOfCompilationWorkers: 3,
					AvailabilityZones: []bftinput.AvailabilityZone{
						{Name: "z1"},
					},
					PersistentDiskPools: []bftinput.DiskConfig{
						{Name: "fake-persistent-disk", Size: 1},
					},
					Networks: []bftinput.NetworkConfig{
						{Name: "foo-network"},
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
		rand.Seed(64)
		inputGenerator = NewInputGenerator(parameters, fakeParameterProvider, 1, nameGenerator, decider, logger)

		inputs, err := inputGenerator.Generate()
		Expect(err).ToNot(HaveOccurred())

		Expect(inputs).To(Equal([]bftinput.Input{
			{
				Jobs: []bftinput.Job{
					{
						Name:               "vgrKicN3O2",
						Instances:          2,
						VmType:             "fake-vm-type",
						AvailabilityZones:  []string{"z1"},
						PersistentDiskPool: "fake-persistent-disk",
						Networks: []bftinput.JobNetworkConfig{
							{Name: "foo-network"},
						},
						Templates: []bftinput.Template{
							{Name: "simple"},
						},
					},
				},
				Update: bftinput.UpdateConfig{
					Canaries:    3,
					MaxInFlight: 5,
					Serial:      "true",
				},
				CloudConfig: bftinput.CloudConfig{
					AvailabilityZones: []bftinput.AvailabilityZone{
						{Name: "z1"},
					},
					PersistentDiskPools: []bftinput.DiskConfig{
						{
							Name: "fake-persistent-disk",
							Size: 1,
						},
					},
					Networks: []bftinput.NetworkConfig{
						{Name: "foo-network"},
					},
					VmTypes: []bftinput.VmTypeConfig{
						{Name: "fake-vm-type"},
					},
					NumberOfCompilationWorkers: 3,
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
						Networks: []bftinput.JobNetworkConfig{
							{Name: "foo-network"},
						},
						MigratedFrom: []bftinput.MigratedFromConfig{
							{Name: "vgrKicN3O2"},
						},
						Templates: []bftinput.Template{
							{Name: "simple"},
						},
					},
				},
				Update: bftinput.UpdateConfig{
					Canaries:    3,
					MaxInFlight: 5,
					Serial:      "true",
				},
				CloudConfig: bftinput.CloudConfig{
					AvailabilityZones: []bftinput.AvailabilityZone{
						{Name: "z1"},
					},
					PersistentDiskPools: []bftinput.DiskConfig{
						{
							Name: "fake-persistent-disk",
							Size: 1,
						},
					},
					Networks: []bftinput.NetworkConfig{
						{Name: "foo-network"},
					},
					VmTypes: []bftinput.VmTypeConfig{
						{Name: "fake-vm-type"},
					},
					NumberOfCompilationWorkers: 3,
				},
				Stemcells: []bftinput.StemcellConfig{
					{Name: "fake-stemcell"},
				},
			},
		}))
	})
})
