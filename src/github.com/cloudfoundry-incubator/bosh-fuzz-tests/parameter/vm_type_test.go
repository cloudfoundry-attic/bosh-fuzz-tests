package parameter_test

import (
	bftinput "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"
	fakebftnamegen "github.com/cloudfoundry-incubator/bosh-fuzz-tests/name_generator/fakes"

	. "github.com/cloudfoundry-incubator/bosh-fuzz-tests/parameter"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("VmType", func() {
	var (
		vmType Parameter
	)

	Context("when definition is vm_type", func() {
		BeforeEach(func() {
			fakeNameGenerator := &fakebftnamegen.FakeNameGenerator{
				Names: []string{"fake-vm-type"},
			}
			vmType = NewVmType("vm_type", fakeNameGenerator)
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
		It("uses previous input", func() {})
	})
})
