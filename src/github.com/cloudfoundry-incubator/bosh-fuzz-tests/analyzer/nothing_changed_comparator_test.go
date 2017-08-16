package analyzer_test

import (
	bftexpectation "github.com/cloudfoundry-incubator/bosh-fuzz-tests/expectation"
	bftinput "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"

	. "github.com/cloudfoundry-incubator/bosh-fuzz-tests/analyzer"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NothingChangedComparator", func() {
	var (
		nothingChangedComparator Comparator
		previousInputs           []bftinput.Input
		currentInput             bftinput.Input
	)

	BeforeEach(func() {
		nothingChangedComparator = NewNothingChangedComparator()
	})

	Context("when there are same instance groups", func() {
		BeforeEach(func() {
			previousInputs = []bftinput.Input{
				{
					InstanceGroups: []bftinput.InstanceGroup{
						{
							Name: "foo-instance-group",
							Networks: []bftinput.InstanceGroupNetworkConfig{
								{
									Name: "network-1",
								},
							},
						},
					},
				},
			}

			currentInput = bftinput.Input{
				InstanceGroups: []bftinput.InstanceGroup{
					{
						Name: "foo-instance-group",
						Networks: []bftinput.InstanceGroupNetworkConfig{
							{
								Name: "network-1",
							},
						},
					},
				},
			}
		})

		It("returns debug log expectation", func() {
			expectations := nothingChangedComparator.Compare(previousInputs, currentInput)
			expectedDebugLogExpectation := bftexpectation.NewDebugLog("No instances to update for 'foo-instance-group'")
			Expect(expectations).To(ContainElement(expectedDebugLogExpectation))
		})
	})

	Context("when there are instance groups that have different properties", func() {
		BeforeEach(func() {
			previousInputs = []bftinput.Input{
				{
					InstanceGroups: []bftinput.InstanceGroup{
						{
							Name: "foo-instance-group",
							Networks: []bftinput.InstanceGroupNetworkConfig{
								{
									Name: "network-1",
								},
							},
						},
					},
				},
			}

			currentInput = bftinput.Input{
				InstanceGroups: []bftinput.InstanceGroup{
					{
						Name: "foo-instance-group",
						Networks: []bftinput.InstanceGroupNetworkConfig{
							{
								Name: "network-2",
							},
						},
					},
				},
			}
		})

		It("returns no expectations", func() {
			expectations := nothingChangedComparator.Compare(previousInputs, currentInput)
			Expect(expectations).To(BeEmpty())
		})
	})

	Context("when az properties was changed", func() {
		BeforeEach(func() {
			previousInputs = []bftinput.Input{
				{
					InstanceGroups: []bftinput.InstanceGroup{
						{
							Name:              "foo-instance-group",
							AvailabilityZones: []string{"z1", "z2"},
						},
					},
					CloudConfig: bftinput.CloudConfig{
						AvailabilityZones: []bftinput.AvailabilityZone{
							{Name: "z1"},
							{Name: "z2"},
						},
					},
				},
			}

			currentInput = bftinput.Input{
				InstanceGroups: []bftinput.InstanceGroup{
					{
						Name:              "foo-instance-group",
						AvailabilityZones: []string{"z1", "z2"},
					},
				},
				CloudConfig: bftinput.CloudConfig{
					AvailabilityZones: []bftinput.AvailabilityZone{
						{Name: "z1"},
						{
							Name: "z2",
							CloudProperties: map[string]string{
								"fake-key": "fake-property",
							},
						},
					},
				},
			}
		})

		It("returns no expectations", func() {
			expectations := nothingChangedComparator.Compare(previousInputs, currentInput)
			Expect(expectations).To(BeEmpty())
		})
	})

	Context("when PersistentDiskPool properties was changed", func() {
		BeforeEach(func() {
			previousInputs = []bftinput.Input{
				{
					InstanceGroups: []bftinput.InstanceGroup{
						{
							Name:               "foo-instance-group",
							PersistentDiskPool: "foo-disk-pool",
						},
					},
					CloudConfig: bftinput.CloudConfig{
						PersistentDiskPools: []bftinput.DiskConfig{
							{Name: "foo-disk-pool", Size: 200},
						},
					},
				},
			}

			currentInput = bftinput.Input{
				InstanceGroups: []bftinput.InstanceGroup{
					{
						Name:               "foo-instance-group",
						PersistentDiskPool: "foo-disk-pool",
					},
				},
				CloudConfig: bftinput.CloudConfig{
					PersistentDiskPools: []bftinput.DiskConfig{
						{Name: "foo-disk-pool", Size: 100},
					},
				},
			}
		})

		It("returns no expectations", func() {
			expectations := nothingChangedComparator.Compare(previousInputs, currentInput)
			Expect(expectations).To(BeEmpty())
		})
	})

	Context("when PersistentDiskType properties was changed", func() {
		BeforeEach(func() {
			previousInputs = []bftinput.Input{
				{
					InstanceGroups: []bftinput.InstanceGroup{
						{
							Name:               "foo-instance-group",
							PersistentDiskType: "foo-disk-type",
						},
					},
					CloudConfig: bftinput.CloudConfig{
						PersistentDiskTypes: []bftinput.DiskConfig{
							{Name: "foo-disk-type", Size: 200},
						},
					},
				},
			}

			currentInput = bftinput.Input{
				InstanceGroups: []bftinput.InstanceGroup{
					{
						Name:               "foo-instance-group",
						PersistentDiskType: "foo-disk-type",
					},
				},
				CloudConfig: bftinput.CloudConfig{
					PersistentDiskTypes: []bftinput.DiskConfig{
						{Name: "foo-disk-type", Size: 100},
					},
				},
			}
		})

		It("returns no expectations", func() {
			expectations := nothingChangedComparator.Compare(previousInputs, currentInput)
			Expect(expectations).To(BeEmpty())
		})
	})

	Context("when Networks properties was changed", func() {
		BeforeEach(func() {
			previousInputs = []bftinput.Input{
				{
					InstanceGroups: []bftinput.InstanceGroup{
						{
							Name: "foo-instance-group",
							Networks: []bftinput.InstanceGroupNetworkConfig{
								{
									Name: "foo-network",
								},
							},
						},
					},
					CloudConfig: bftinput.CloudConfig{
						Networks: []bftinput.NetworkConfig{
							{
								Name: "foo-network",
								Subnets: []bftinput.SubnetConfig{
									{
										IpPool: &bftinput.IpPool{
											IpRange: "192.168.0.0/24",
										},
									},
								},
							},
						},
					},
				},
			}

			currentInput = bftinput.Input{
				InstanceGroups: []bftinput.InstanceGroup{
					{
						Name: "foo-instance-group",
						Networks: []bftinput.InstanceGroupNetworkConfig{
							{
								Name: "foo-network",
							},
						},
					},
				},
				CloudConfig: bftinput.CloudConfig{
					Networks: []bftinput.NetworkConfig{
						{
							Name: "foo-network",
							Subnets: []bftinput.SubnetConfig{
								{
									IpPool: &bftinput.IpPool{
										IpRange: "192.168.10.0/24",
									},
								},
							},
						},
					},
				},
			}
		})

		It("returns no expectations", func() {
			expectations := nothingChangedComparator.Compare(previousInputs, currentInput)
			Expect(expectations).To(BeEmpty())
		})
	})

	Context("when ResourcePool property was changed", func() {
		BeforeEach(func() {
			previousInputs = []bftinput.Input{
				{
					InstanceGroups: []bftinput.InstanceGroup{
						{
							Name:         "foo-instance-group",
							ResourcePool: "foo-resource-pool",
						},
					},
					CloudConfig: bftinput.CloudConfig{
						ResourcePools: []bftinput.ResourcePoolConfig{
							{
								Name: "foo-resource-pool",
								Stemcell: bftinput.StemcellConfig{
									Name: "foo-name-one",
								},
							},
						},
					},
				},
			}

			currentInput = bftinput.Input{
				InstanceGroups: []bftinput.InstanceGroup{
					{
						Name:         "foo-instance-group",
						ResourcePool: "foo-resource-pool",
					},
				},
				CloudConfig: bftinput.CloudConfig{
					ResourcePools: []bftinput.ResourcePoolConfig{
						{
							Name: "foo-resource-pool",
							Stemcell: bftinput.StemcellConfig{
								Name: "foo-name-two",
							},
						},
					},
				},
			}
		})

		It("returns no expectations", func() {
			expectations := nothingChangedComparator.Compare(previousInputs, currentInput)
			Expect(expectations).To(BeEmpty())
		})
	})

	Context("when VmType property was changed", func() {
		BeforeEach(func() {
			previousInputs = []bftinput.Input{
				{
					InstanceGroups: []bftinput.InstanceGroup{
						{
							Name:   "foo-instance-group",
							VmType: "foo-vm-type",
						},
					},
					CloudConfig: bftinput.CloudConfig{
						VmTypes: []bftinput.VmTypeConfig{
							{
								Name: "foo-vm-type",
								CloudProperties: map[string]string{
									"fake-key": "fake-property",
								},
							},
						},
					},
				},
			}

			currentInput = bftinput.Input{
				InstanceGroups: []bftinput.InstanceGroup{
					{
						Name:   "foo-instance-group",
						VmType: "foo-vm-type",
					},
				},
				CloudConfig: bftinput.CloudConfig{
					VmTypes: []bftinput.VmTypeConfig{
						{
							Name: "foo-vm-type",
							CloudProperties: map[string]string{
								"fake-key": "fake-updated-property",
							},
						},
					},
				},
			}
		})

		It("returns no expectations", func() {
			expectations := nothingChangedComparator.Compare(previousInputs, currentInput)
			Expect(expectations).To(BeEmpty())
		})
	})

	Context("when Stemcell property was changed", func() {
		BeforeEach(func() {
			previousInputs = []bftinput.Input{
				{
					InstanceGroups: []bftinput.InstanceGroup{
						{
							Name:     "foo-instance-group",
							Stemcell: "foo-stemcell",
						},
					},
					Stemcells: []bftinput.StemcellConfig{
						{
							Name:    "foo-stemcell",
							Version: "1",
						},
					},
				},
			}

			currentInput = bftinput.Input{
				InstanceGroups: []bftinput.InstanceGroup{
					{
						Name:     "foo-instance-group",
						Stemcell: "foo-stemcell",
					},
				},
				Stemcells: []bftinput.StemcellConfig{
					{
						Name:    "foo-stemcell",
						Version: "2",
					},
				},
			}
		})

		It("returns no expectations", func() {
			expectations := nothingChangedComparator.Compare(previousInputs, currentInput)
			Expect(expectations).To(BeEmpty())
		})
	})

	Context("when persistent disk was removed in previous input", func() {
		BeforeEach(func() {
			previousInputs = []bftinput.Input{
				{
					InstanceGroups: []bftinput.InstanceGroup{
						{
							Name:               "foo-instance-group",
							PersistentDiskPool: "foo-disk-pool",
						},
					},
					CloudConfig: bftinput.CloudConfig{
						PersistentDiskPools: []bftinput.DiskConfig{
							{
								Name: "foo-disk-pool",
								Size: 100,
							},
						},
					},
				},
				{
					InstanceGroups: []bftinput.InstanceGroup{
						{
							Name: "foo-instance-group",
						},
					},
				},
			}

			currentInput = bftinput.Input{
				InstanceGroups: []bftinput.InstanceGroup{
					{
						Name: "foo-instance-group",
					},
				},
			}
		})

		It("returns no expectations", func() {
			expectations := nothingChangedComparator.Compare(previousInputs, currentInput)
			Expect(expectations).To(BeEmpty())
		})
	})

	Context("when persistent disk was removed and instance group was migrated in previous input", func() {
		BeforeEach(func() {
			previousInputs = []bftinput.Input{
				{
					InstanceGroups: []bftinput.InstanceGroup{
						{
							Name:               "bar-instance-group",
							PersistentDiskPool: "foo-disk-pool",
						},
					},
					CloudConfig: bftinput.CloudConfig{
						PersistentDiskPools: []bftinput.DiskConfig{
							{
								Name: "foo-disk-pool",
								Size: 100,
							},
						},
					},
				},
				{
					InstanceGroups: []bftinput.InstanceGroup{
						{
							Name: "foo-instance-group",
							MigratedFrom: []bftinput.MigratedFromConfig{
								{
									Name: "bar-instance-group",
								},
							},
						},
					},
				},
			}

			currentInput = bftinput.Input{
				InstanceGroups: []bftinput.InstanceGroup{
					{
						Name: "foo-instance-group",
						MigratedFrom: []bftinput.MigratedFromConfig{
							{
								Name: "bar-instance-group",
							},
						},
					},
				},
			}
		})

		It("returns no expectations", func() {
			expectations := nothingChangedComparator.Compare(previousInputs, currentInput)
			Expect(expectations).To(BeEmpty())
		})
	})

	Context("when Lifecycle property is errand", func() {
		BeforeEach(func() {
			previousInputs = []bftinput.Input{
				{
					InstanceGroups: []bftinput.InstanceGroup{
						{
							Name:      "foo-instance-group",
							Lifecycle: "errand",
						},
					},
				},
			}

			currentInput = bftinput.Input{
				InstanceGroups: []bftinput.InstanceGroup{
					{
						Name:      "foo-instance-group",
						Lifecycle: "errand",
					},
				},
			}
		})

		It("returns no expectations", func() {
			expectations := nothingChangedComparator.Compare(previousInputs, currentInput)
			Expect(expectations).To(BeEmpty())
		})
	})
})
