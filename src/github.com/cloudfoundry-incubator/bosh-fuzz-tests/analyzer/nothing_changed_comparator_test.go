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
		previousInput            bftinput.Input
		currentInput             bftinput.Input
	)

	BeforeEach(func() {
		nothingChangedComparator = NewNothingChangedComparator()
	})

	Context("when there are same jobs", func() {
		BeforeEach(func() {
			previousInput = bftinput.Input{
				Jobs: []bftinput.Job{
					{
						Name: "foo-job",
						Networks: []bftinput.JobNetworkConfig{
							{
								Name: "network-1",
							},
						},
					},
				},
			}

			currentInput = bftinput.Input{
				Jobs: []bftinput.Job{
					{
						Name: "foo-job",
						Networks: []bftinput.JobNetworkConfig{
							{
								Name: "network-1",
							},
						},
					},
				},
			}
		})

		It("returns debug log expectation", func() {
			expectations := nothingChangedComparator.Compare(previousInput, currentInput)
			expectedDebugLogExpectation := bftexpectation.NewDebugLog("No instances to update for 'foo-job'")
			Expect(expectations).To(ContainElement(expectedDebugLogExpectation))
		})
	})

	Context("when there are jobs that have different properties", func() {
		BeforeEach(func() {
			previousInput = bftinput.Input{
				Jobs: []bftinput.Job{
					{
						Name: "foo-job",
						Networks: []bftinput.JobNetworkConfig{
							{
								Name: "network-1",
							},
						},
					},
				},
			}

			currentInput = bftinput.Input{
				Jobs: []bftinput.Job{
					{
						Name: "foo-job",
						Networks: []bftinput.JobNetworkConfig{
							{
								Name: "network-2",
							},
						},
					},
				},
			}
		})

		It("returns no expectations", func() {
			expectations := nothingChangedComparator.Compare(previousInput, currentInput)
			Expect(expectations).To(BeEmpty())
		})
	})

	Context("when az properties was changed", func() {
		BeforeEach(func() {
			previousInput = bftinput.Input{
				Jobs: []bftinput.Job{
					{
						Name:              "foo-job",
						AvailabilityZones: []string{"z1", "z2"},
					},
				},
				CloudConfig: bftinput.CloudConfig{
					AvailabilityZones: []bftinput.AvailabilityZone{
						{Name: "z1"},
						{Name: "z2"},
					},
				},
			}

			currentInput = bftinput.Input{
				Jobs: []bftinput.Job{
					{
						Name:              "foo-job",
						AvailabilityZones: []string{"z1", "z2"},
					},
				},
				CloudConfig: bftinput.CloudConfig{
					AvailabilityZones: []bftinput.AvailabilityZone{
						{Name: "z1"},
						{
							Name: "z2",
							CloudProperties: map[string]interface{}{
								"fake-key": "fake-property",
							},
						},
					},
				},
			}
		})

		It("returns no expectations", func() {
			expectations := nothingChangedComparator.Compare(previousInput, currentInput)
			Expect(expectations).To(BeEmpty())
		})
	})

	Context("when PersistentDiskPool properties was changed", func() {
		BeforeEach(func() {
			previousInput = bftinput.Input{
				Jobs: []bftinput.Job{
					{
						Name:               "foo-job",
						PersistentDiskPool: "foo-disk-pool",
					},
				},
				CloudConfig: bftinput.CloudConfig{
					PersistentDiskPools: []bftinput.DiskConfig{
						{Name: "foo-disk-pool", Size: 200},
					},
				},
			}

			currentInput = bftinput.Input{
				Jobs: []bftinput.Job{
					{
						Name:               "foo-job",
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
			expectations := nothingChangedComparator.Compare(previousInput, currentInput)
			Expect(expectations).To(BeEmpty())
		})
	})

	Context("when PersistentDiskType properties was changed", func() {
		BeforeEach(func() {
			previousInput = bftinput.Input{
				Jobs: []bftinput.Job{
					{
						Name:               "foo-job",
						PersistentDiskType: "foo-disk-type",
					},
				},
				CloudConfig: bftinput.CloudConfig{
					PersistentDiskTypes: []bftinput.DiskConfig{
						{Name: "foo-disk-type", Size: 200},
					},
				},
			}

			currentInput = bftinput.Input{
				Jobs: []bftinput.Job{
					{
						Name:               "foo-job",
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
			expectations := nothingChangedComparator.Compare(previousInput, currentInput)
			Expect(expectations).To(BeEmpty())
		})
	})

	Context("when Networks properties was changed", func() {
		BeforeEach(func() {
			previousInput = bftinput.Input{
				Jobs: []bftinput.Job{
					{
						Name: "foo-job",
						Networks: []bftinput.JobNetworkConfig{
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
			}

			currentInput = bftinput.Input{
				Jobs: []bftinput.Job{
					{
						Name: "foo-job",
						Networks: []bftinput.JobNetworkConfig{
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
			expectations := nothingChangedComparator.Compare(previousInput, currentInput)
			Expect(expectations).To(BeEmpty())
		})
	})

	Context("when ResourcePool property was changed", func() {
		BeforeEach(func() {
			previousInput = bftinput.Input{
				Jobs: []bftinput.Job{
					{
						Name:         "foo-job",
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
			}

			currentInput = bftinput.Input{
				Jobs: []bftinput.Job{
					{
						Name:         "foo-job",
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
			expectations := nothingChangedComparator.Compare(previousInput, currentInput)
			Expect(expectations).To(BeEmpty())
		})
	})

	Context("when VmType property was changed", func() {
		BeforeEach(func() {
			previousInput = bftinput.Input{
				Jobs: []bftinput.Job{
					{
						Name:   "foo-job",
						VmType: "foo-vm-type",
					},
				},
				CloudConfig: bftinput.CloudConfig{
					VmTypes: []bftinput.VmTypeConfig{
						{
							Name: "foo-vm-type",
							CloudProperties: map[string]interface{}{
								"fake-key": "fake-property",
							},
						},
					},
				},
			}

			currentInput = bftinput.Input{
				Jobs: []bftinput.Job{
					{
						Name:   "foo-job",
						VmType: "foo-vm-type",
					},
				},
				CloudConfig: bftinput.CloudConfig{
					VmTypes: []bftinput.VmTypeConfig{
						{
							Name: "foo-vm-type",
							CloudProperties: map[string]interface{}{
								"fake-key": "fake-updated-property",
							},
						},
					},
				},
			}
		})

		It("returns no expectations", func() {
			expectations := nothingChangedComparator.Compare(previousInput, currentInput)
			Expect(expectations).To(BeEmpty())
		})
	})

	Context("when Stemcell property was changed", func() {
		BeforeEach(func() {
			previousInput = bftinput.Input{
				Jobs: []bftinput.Job{
					{
						Name:     "foo-job",
						Stemcell: "foo-stemcell",
					},
				},
				Stemcells: []bftinput.StemcellConfig{
					{
						Name:    "foo-stemcell",
						Version: "1",
					},
				},
			}

			currentInput = bftinput.Input{
				Jobs: []bftinput.Job{
					{
						Name:     "foo-job",
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
			expectations := nothingChangedComparator.Compare(previousInput, currentInput)
			Expect(expectations).To(BeEmpty())
		})
	})
})
