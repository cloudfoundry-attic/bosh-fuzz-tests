package parameter_test

import (
	fakebftdecider "github.com/cloudfoundry-incubator/bosh-fuzz-tests/decider/fakes"
	bftinput "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"
	fakebftnamegen "github.com/cloudfoundry-incubator/bosh-fuzz-tests/name_generator/fakes"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"

	. "github.com/cloudfoundry-incubator/bosh-fuzz-tests/parameter"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("VmType", func() {
	var (
		fakeNameGenerator *fakebftnamegen.FakeNameGenerator
		fakeDecider       *fakebftdecider.FakeDecider
		logger            boshlog.Logger
		vmType            Parameter
		availabilityZones []bftinput.AvailabilityZone
	)

	BeforeEach(func() {
		fakeNameGenerator = &fakebftnamegen.FakeNameGenerator{
			Names: []string{"fake-vm-type"},
		}
		fakeDecider = &fakebftdecider.FakeDecider{}
		logger = boshlog.NewLogger(boshlog.LevelNone)
	})

	Context("when using cloud-config for vm_type", func() {
		BeforeEach(func() {
			fakeDecider.IsYesYes = false
			vmType = NewVmType(fakeNameGenerator, fakeDecider, logger)
			availabilityZones = []bftinput.AvailabilityZone{{Name: "z1"}}
		})

		It("adds vm_types to the input", func() {
			input := bftinput.Input{
				InstanceGroups: []bftinput.InstanceGroup{
					{
						Name: "fake-instance-group",
					},
				},
				CloudConfig: bftinput.CloudConfig{
					AvailabilityZones: availabilityZones,
				},
			}

			result := vmType.Apply(input, bftinput.Input{})

			Expect(result).To(Equal(bftinput.Input{
				InstanceGroups: []bftinput.InstanceGroup{
					{
						Name:   "fake-instance-group",
						VmType: "fake-vm-type",
					},
				},
				CloudConfig: bftinput.CloudConfig{
					AvailabilityZones: availabilityZones,
					VmTypes: []bftinput.VmTypeConfig{
						{
							Name: "fake-vm-type",
						},
					},
				},
			}))
		})
	})

	Context("when it is decided to keep previous input", func() {
		BeforeEach(func() {
			fakeDecider.IsYesYes = true
			vmType = NewVmType(fakeNameGenerator, fakeDecider, logger)
		})

		It("uses previous input", func() {
			input := bftinput.Input{
				InstanceGroups: []bftinput.InstanceGroup{
					{
						Name:   "fake-instance-group",
						VmType: "previous-vm-type",
					},
				},
				CloudConfig: bftinput.CloudConfig{
					AvailabilityZones: availabilityZones,
					VmTypes: []bftinput.VmTypeConfig{
						{
							Name: "previous-vm-type",
						},
					},
				},
			}

			result := vmType.Apply(input, bftinput.Input{})

			Expect(result).To(Equal(bftinput.Input{
				InstanceGroups: []bftinput.InstanceGroup{
					{
						Name:   "fake-instance-group",
						VmType: "previous-vm-type",
					},
				},
				CloudConfig: bftinput.CloudConfig{
					AvailabilityZones: availabilityZones,
					VmTypes: []bftinput.VmTypeConfig{
						{
							Name: "previous-vm-type",
						},
					},
				},
			}))
		})
	})

	Context("when it is decided to share vm types", func() {
		BeforeEach(func() {
			fakeDecider.IsYesYes = true
			vmType = NewVmType(fakeNameGenerator, fakeDecider, logger)
		})

		It("sets same vm type on input instance groups", func() {
			input := bftinput.Input{
				InstanceGroups: []bftinput.InstanceGroup{
					{
						Name: "fake-instance-group-1",
					},
					{
						Name: "fake-instance-group-2",
					},
				},
				CloudConfig: bftinput.CloudConfig{
					AvailabilityZones: availabilityZones,
				},
			}

			result := vmType.Apply(input, bftinput.Input{})

			Expect(result).To(Equal(bftinput.Input{
				InstanceGroups: []bftinput.InstanceGroup{
					{
						Name:   "fake-instance-group-1",
						VmType: "fake-vm-type",
					},
					{
						Name:   "fake-instance-group-2",
						VmType: "fake-vm-type",
					},
				},
				CloudConfig: bftinput.CloudConfig{
					AvailabilityZones: availabilityZones,
					VmTypes: []bftinput.VmTypeConfig{
						{
							Name: "fake-vm-type",
						},
					},
				},
			}))
		})
	})

	Context("when it is decided to share resource pool", func() {
		BeforeEach(func() {
			fakeDecider.IsYesYes = true
			vmType = NewVmType(fakeNameGenerator, fakeDecider, logger)
		})

		It("sets same vm type on input instance groups", func() {
			input := bftinput.Input{
				InstanceGroups: []bftinput.InstanceGroup{
					{
						Name: "fake-instance-group-1",
					},
					{
						Name: "fake-instance-group-2",
					},
				},
			}

			result := vmType.Apply(input, bftinput.Input{})

			Expect(result).To(Equal(bftinput.Input{
				InstanceGroups: []bftinput.InstanceGroup{
					{
						Name:         "fake-instance-group-1",
						ResourcePool: "fake-vm-type",
					},
					{
						Name:         "fake-instance-group-2",
						ResourcePool: "fake-vm-type",
					},
				},
				CloudConfig: bftinput.CloudConfig{
					ResourcePools: []bftinput.ResourcePoolConfig{
						{
							Name: "fake-vm-type",
						},
					},
				},
			}))
		})
	})
})
