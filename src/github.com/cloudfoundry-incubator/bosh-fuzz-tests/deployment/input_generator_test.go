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
			NameLength:               []int{5},
			Instances:                []int{2},
			AvailabilityZones:        [][]string{[]string{"z1"}},
			PersistentDiskDefinition: []string{"persistent_disk_size"},
			PersistentDiskSize:       []int{0},
			NumberOfInstanceGroups:   []int{1},
			MigratedFromCount:        []int{1},
			VmTypeDefinition:         []string{"vm_type"},
			StemcellDefinition:       []string{"name"},
			Jobs:                     [][]string{[]string{"simple"}},
			NumberOfCompilationWorkers: []int{3},
			Canaries:                   []int{5},
			MaxInFlight:                []int{3},
			Serial:                     []string{"true"},
			NumOfCloudProperties:       []int{2},
			NumOfVariables:             []int{0, 10, 500},
			VariableTypes:              []string{"rsa", "ssh", "password"},
		}
		logger = boshlog.NewLogger(boshlog.LevelNone)
		nameGenerator = bftnamegen.NewNameGenerator()
		fakeParameterProvider = fakebftparam.NewFakeParameterProvider("disk_pool", "vm_type")
		decider = &fakebftdecider.FakeDecider{}
	})

	It("generates requested number of inputs", func() {
		parameters = bftconfig.Parameters{
			NameLength:             []int{5},
			Instances:              []int{2},
			NumberOfInstanceGroups: []int{2},
			MigratedFromCount:      []int{0},
			NumOfCloudProperties:   []int{2},
		}

		rand.Seed(64)
		inputGenerator = NewInputGenerator(parameters, fakeParameterProvider, 2, nameGenerator, decider, logger)

		inputs, err := inputGenerator.Generate()
		Expect(err).ToNot(HaveOccurred())

		Expect(inputs).To(Equal([]bftinput.Input{
			{
				InstanceGroups: []bftinput.InstanceGroup{
					{
						Name:               "joNAw",
						Instances:          2,
						AvailabilityZones:  []string{"z1"},
						Lifecycle:          "mufasa",
						PersistentDiskPool: "fake-persistent-disk",
						VmType:             "fake-vm-type",
						Networks: []bftinput.InstanceGroupNetworkConfig{
							{Name: "foo-network"},
						},
						Jobs: []bftinput.Job{
							{Name: "simple"},
						},
					},
					{
						Name:               "gQ8el",
						Instances:          2,
						Lifecycle:          "mufasa",
						AvailabilityZones:  []string{"z1"},
						PersistentDiskPool: "fake-persistent-disk",
						VmType:             "fake-vm-type",
						Networks: []bftinput.InstanceGroupNetworkConfig{
							{Name: "foo-network"},
						},
						Jobs: []bftinput.Job{
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
						{
							Name: "foo-network",
							Subnets: []bftinput.SubnetConfig{
								{
									IpPool: &bftinput.IpPool{
										IpRange: "10.0.0.0/24",
									},
									CloudProperties: map[string]string{
										"foo": "bar",
										"baz": "qux",
									},
								},
							},
						},
					},
					AvailabilityZones: []bftinput.AvailabilityZone{
						{
							Name: "z1",
							CloudProperties: map[string]string{
								"foo": "bar",
								"baz": "qux",
							},
						},
					},
					PersistentDiskPools: []bftinput.DiskConfig{
						{
							Name: "fake-persistent-disk",
							Size: 1,
							CloudProperties: map[string]string{
								"foo": "bar",
								"baz": "qux",
							},
						},
					},
					VmTypes: []bftinput.VmTypeConfig{
						{
							Name: "fake-vm-type",
							CloudProperties: map[string]string{
								"foo": "bar",
								"baz": "qux",
							},
						},
					},
					Compilation: bftinput.CompilationConfig{
						NumberOfWorkers: 3,
						CloudProperties: map[string]string{
							"foo": "bar",
							"baz": "qux",
						},
					},
				},
				Stemcells: []bftinput.StemcellConfig{
					{Name: "fake-stemcell"},
				},
			},
			{
				InstanceGroups: []bftinput.InstanceGroup{
					{
						Name:               "joNAw",
						Instances:          2,
						Lifecycle:          "mufasa",
						AvailabilityZones:  []string{"z1"},
						PersistentDiskPool: "fake-persistent-disk",
						VmType:             "fake-vm-type",
						Networks: []bftinput.InstanceGroupNetworkConfig{
							{Name: "foo-network"},
						},
						Jobs: []bftinput.Job{
							{Name: "simple"},
						},
					},
					{
						Name:               "gQ8el",
						Instances:          2,
						Lifecycle:          "mufasa",
						AvailabilityZones:  []string{"z1"},
						PersistentDiskPool: "fake-persistent-disk",
						VmType:             "fake-vm-type",
						Networks: []bftinput.InstanceGroupNetworkConfig{
							{Name: "foo-network"},
						},
						Jobs: []bftinput.Job{
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
					Compilation: bftinput.CompilationConfig{
						NumberOfWorkers: 3,
						CloudProperties: map[string]string{
							"foo": "bar",
							"baz": "qux",
						},
					},
					AvailabilityZones: []bftinput.AvailabilityZone{
						{
							Name: "z1",
							CloudProperties: map[string]string{
								"foo": "bar",
								"baz": "qux",
							},
						},
					},
					PersistentDiskPools: []bftinput.DiskConfig{
						{
							Name: "fake-persistent-disk",
							Size: 1,
							CloudProperties: map[string]string{
								"foo": "bar",
								"baz": "qux",
							},
						},
					},
					Networks: []bftinput.NetworkConfig{
						{
							Name: "foo-network",
							Subnets: []bftinput.SubnetConfig{
								{
									IpPool: &bftinput.IpPool{
										IpRange: "10.0.0.0/24",
									},
									CloudProperties: map[string]string{
										"foo": "bar",
										"baz": "qux",
									},
								},
							},
						},
					},
					VmTypes: []bftinput.VmTypeConfig{
						{
							Name: "fake-vm-type",
							CloudProperties: map[string]string{
								"foo": "bar",
								"baz": "qux",
							},
						},
					},
				},
				Stemcells: []bftinput.StemcellConfig{
					{Name: "fake-stemcell"},
				},
			},
		}))
	})

	It("generates manifest for dry run", func() {
		decider = &fakebftdecider.FakeDecider{IsYesYes: true}
		parameters = bftconfig.Parameters{
			NameLength:             []int{5},
			Instances:              []int{2},
			NumberOfInstanceGroups: []int{2},
			MigratedFromCount:      []int{0},
			NumOfCloudProperties:   []int{2},
		}

		rand.Seed(64)
		inputGenerator = NewInputGenerator(parameters, fakeParameterProvider, 2, nameGenerator, decider, logger)

		inputs, err := inputGenerator.Generate()
		Expect(err).ToNot(HaveOccurred())

		Expect(inputs[1].IsDryRun).To(BeTrue())
	})

	It("generates requested number of inputs using disk_type", func() {
		parameters = bftconfig.Parameters{
			NameLength:             []int{5},
			Instances:              []int{2},
			NumberOfInstanceGroups: []int{1},
			MigratedFromCount:      []int{0},
			NumOfCloudProperties:   []int{2},
		}

		rand.Seed(64)
		fakeParameterProvider = fakebftparam.NewFakeParameterProvider("disk_type", "vm_type")
		inputGenerator = NewInputGenerator(parameters, fakeParameterProvider, 1, nameGenerator, decider, logger)

		inputs, err := inputGenerator.Generate()
		Expect(err).ToNot(HaveOccurred())

		Expect(inputs).To(Equal([]bftinput.Input{
			{
				InstanceGroups: []bftinput.InstanceGroup{
					{
						Name:               "joNAw",
						Instances:          2,
						Lifecycle:          "mufasa",
						AvailabilityZones:  []string{"z1"},
						PersistentDiskType: "fake-persistent-disk",
						VmType:             "fake-vm-type",
						Networks: []bftinput.InstanceGroupNetworkConfig{
							{Name: "foo-network"},
						},
						Jobs: []bftinput.Job{
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
						{
							Name: "foo-network",
							Subnets: []bftinput.SubnetConfig{
								{
									IpPool: &bftinput.IpPool{
										IpRange: "10.0.0.0/24",
									},
									CloudProperties: map[string]string{
										"foo": "bar",
										"baz": "qux",
									},
								},
							},
						},
					},
					AvailabilityZones: []bftinput.AvailabilityZone{
						{
							Name: "z1",
							CloudProperties: map[string]string{
								"foo": "bar",
								"baz": "qux",
							},
						},
					},
					PersistentDiskTypes: []bftinput.DiskConfig{
						{
							Name: "fake-persistent-disk",
							Size: 1,
							CloudProperties: map[string]string{
								"foo": "bar",
								"baz": "qux",
							},
						},
					},
					VmTypes: []bftinput.VmTypeConfig{
						{
							Name: "fake-vm-type",
							CloudProperties: map[string]string{
								"foo": "bar",
								"baz": "qux",
							},
						},
					},
					Compilation: bftinput.CompilationConfig{
						NumberOfWorkers: 3,
						CloudProperties: map[string]string{
							"foo": "bar",
							"baz": "qux",
						},
					},
				},
				Stemcells: []bftinput.StemcellConfig{
					{Name: "fake-stemcell"},
				},
			},
		}))
	})

	It("generates requested number of inputs using disk_type and resource pool", func() {
		parameters = bftconfig.Parameters{
			NameLength:             []int{5},
			Instances:              []int{2},
			NumberOfInstanceGroups: []int{1},
			MigratedFromCount:      []int{0},
			NumOfCloudProperties:   []int{2},
		}

		rand.Seed(64)
		fakeParameterProvider = fakebftparam.NewFakeParameterProvider("disk_type", "resource_pool")
		inputGenerator = NewInputGenerator(parameters, fakeParameterProvider, 1, nameGenerator, decider, logger)

		inputs, err := inputGenerator.Generate()
		Expect(err).ToNot(HaveOccurred())

		Expect(inputs).To(Equal([]bftinput.Input{
			{
				InstanceGroups: []bftinput.InstanceGroup{
					{
						Name:               "joNAw",
						Instances:          2,
						Lifecycle:          "mufasa",
						AvailabilityZones:  []string{"z1"},
						PersistentDiskType: "fake-persistent-disk",
						Networks: []bftinput.InstanceGroupNetworkConfig{
							{Name: "foo-network"},
						},
						ResourcePool: "fake-resource-pool",
						Jobs: []bftinput.Job{
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
						{
							Name: "foo-network",
							Subnets: []bftinput.SubnetConfig{
								{
									IpPool: &bftinput.IpPool{
										IpRange: "10.0.0.0/24",
									},
									CloudProperties: map[string]string{
										"foo": "bar",
										"baz": "qux",
									},
								},
							},
						},
					},
					AvailabilityZones: []bftinput.AvailabilityZone{
						{
							Name: "z1",
							CloudProperties: map[string]string{
								"foo": "bar",
								"baz": "qux",
							},
						},
					},
					PersistentDiskTypes: []bftinput.DiskConfig{
						{
							Name: "fake-persistent-disk",
							Size: 1,
							CloudProperties: map[string]string{
								"foo": "bar",
								"baz": "qux",
							},
						},
					},
					ResourcePools: []bftinput.ResourcePoolConfig{
						{
							Name: "fake-resource-pool",
							Stemcell: bftinput.StemcellConfig{
								Name:    "foo-stemcell",
								Version: "1",
							},
							CloudProperties: map[string]string{
								"foo": "bar",
								"baz": "qux",
							},
						},
					},
					Compilation: bftinput.CompilationConfig{
						NumberOfWorkers: 3,
						CloudProperties: map[string]string{
							"foo": "bar",
							"baz": "qux",
						},
					},
				},
				Stemcells: []bftinput.StemcellConfig{
					{Name: "fake-stemcell"},
				},
			},
		}))
	})

	It("when migrated instance group does not have az it sets random az in migrated_from", func() {
		rand.Seed(64)
		inputGenerator = NewInputGenerator(parameters, fakeParameterProvider, 1, nameGenerator, decider, logger)

		inputs, err := inputGenerator.Generate()
		Expect(err).ToNot(HaveOccurred())

		Expect(inputs).To(Equal([]bftinput.Input{
			{
				InstanceGroups: []bftinput.InstanceGroup{
					{
						Name:               "vgrKicN3O2",
						Instances:          2,
						Lifecycle:          "mufasa",
						VmType:             "fake-vm-type",
						AvailabilityZones:  []string{"z1"},
						PersistentDiskPool: "fake-persistent-disk",
						Networks: []bftinput.InstanceGroupNetworkConfig{
							{Name: "foo-network"},
						},
						Jobs: []bftinput.Job{
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
						{
							Name: "z1",
							CloudProperties: map[string]string{
								"foo": "bar",
								"baz": "qux",
							},
						},
					},
					PersistentDiskPools: []bftinput.DiskConfig{
						{
							Name: "fake-persistent-disk",
							Size: 1,
							CloudProperties: map[string]string{
								"foo": "bar",
								"baz": "qux",
							},
						},
					},
					Networks: []bftinput.NetworkConfig{
						{
							Name: "foo-network",
							Subnets: []bftinput.SubnetConfig{
								{
									IpPool: &bftinput.IpPool{
										IpRange: "10.0.0.0/24",
									},
									CloudProperties: map[string]string{
										"foo": "bar",
										"baz": "qux",
									},
								},
							},
						},
					},
					VmTypes: []bftinput.VmTypeConfig{
						{
							Name: "fake-vm-type",
							CloudProperties: map[string]string{
								"foo": "bar",
								"baz": "qux",
							},
						},
					},
					Compilation: bftinput.CompilationConfig{
						NumberOfWorkers: 3,
						CloudProperties: map[string]string{
							"foo": "bar",
							"baz": "qux",
						},
					},
				},
				Stemcells: []bftinput.StemcellConfig{
					{Name: "fake-stemcell"},
				},
			},
			{
				InstanceGroups: []bftinput.InstanceGroup{
					{
						Name:               "joNAw",
						Instances:          2,
						Lifecycle:          "mufasa",
						AvailabilityZones:  []string{"z1"},
						VmType:             "fake-vm-type",
						PersistentDiskPool: "fake-persistent-disk",
						Networks: []bftinput.InstanceGroupNetworkConfig{
							{Name: "foo-network"},
						},
						MigratedFrom: []bftinput.MigratedFromConfig{
							{Name: "vgrKicN3O2"},
						},
						Jobs: []bftinput.Job{
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
						{
							Name: "z1",
							CloudProperties: map[string]string{
								"foo": "bar",
								"baz": "qux",
							},
						},
					},
					PersistentDiskPools: []bftinput.DiskConfig{
						{
							Name: "fake-persistent-disk",
							Size: 1,
							CloudProperties: map[string]string{
								"foo": "bar",
								"baz": "qux",
							},
						},
					},
					Networks: []bftinput.NetworkConfig{
						{
							Name: "foo-network",
							Subnets: []bftinput.SubnetConfig{
								{
									IpPool: &bftinput.IpPool{
										IpRange: "10.0.0.0/24",
									},
									CloudProperties: map[string]string{
										"foo": "bar",
										"baz": "qux",
									},
								},
							},
						},
					},
					VmTypes: []bftinput.VmTypeConfig{
						{
							Name: "fake-vm-type",
							CloudProperties: map[string]string{
								"foo": "bar",
								"baz": "qux",
							},
						},
					},
					Compilation: bftinput.CompilationConfig{
						NumberOfWorkers: 3,
						CloudProperties: map[string]string{
							"foo": "bar",
							"baz": "qux",
						},
					},
				},
				Stemcells: []bftinput.StemcellConfig{
					{Name: "fake-stemcell"},
				},
			},
		}))
	})
})
