package parameter_test

import (
	fakebftdecider "github.com/cloudfoundry-incubator/bosh-fuzz-tests/decider/fakes"
	bftinput "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"
	fakebftnamegen "github.com/cloudfoundry-incubator/bosh-fuzz-tests/name_generator/fakes"

	. "github.com/cloudfoundry-incubator/bosh-fuzz-tests/parameter"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("VmType", func() {
	var (
		fakeNameGenerator *fakebftnamegen.FakeNameGenerator
		fakeDecider       *fakebftdecider.FakeDecider
		vmType            Parameter
	)

	BeforeEach(func() {
		fakeNameGenerator = &fakebftnamegen.FakeNameGenerator{
			Names: []string{"fake-vm-type"},
		}
		fakeDecider = &fakebftdecider.FakeDecider{}
	})

	Context("when definition is vm_type", func() {
		BeforeEach(func() {
			fakeDecider.IsYesYes = false
			vmType = NewVmType("vm_type", fakeNameGenerator, fakeDecider)
		})

		It("adds vm_types to the input", func() {
			input := bftinput.Input{
				Jobs: []bftinput.Job{
					{
						Name: "fake-job",
					},
				},
			}

			result := vmType.Apply(input)

			Expect(result).To(Equal(bftinput.Input{
				Jobs: []bftinput.Job{
					{
						Name:   "fake-job",
						VmType: "fake-vm-type",
					},
				},
				CloudConfig: bftinput.CloudConfig{
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
			vmType = NewVmType("vm_type", fakeNameGenerator, fakeDecider)
		})

		It("uses previous input", func() {
			input := bftinput.Input{
				Jobs: []bftinput.Job{
					{
						Name:   "fake-job",
						VmType: "previous-vm-type",
					},
				},
			}

			result := vmType.Apply(input)

			Expect(result).To(Equal(bftinput.Input{
				Jobs: []bftinput.Job{
					{
						Name:   "fake-job",
						VmType: "previous-vm-type",
					},
				},
				CloudConfig: bftinput.CloudConfig{
					VmTypes: []bftinput.VmTypeConfig{
						{
							Name: "previous-vm-type",
						},
					},
				},
			}))
		})
	})
})
