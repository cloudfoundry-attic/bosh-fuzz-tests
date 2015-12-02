package parameter_test

import (
	bftinput "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"

	. "github.com/cloudfoundry-incubator/bosh-fuzz-tests/parameter"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Stemcell", func() {
	var (
		stemcell Parameter
	)

	Context("when definition is os", func() {
		BeforeEach(func() {
			stemcell = NewStemcell("os")
		})

		Context("when input has vm types", func() {
			It("adds stemcells to the input", func() {
				input := &bftinput.Input{
					CloudConfig: bftinput.CloudConfig{
						VmTypes: []bftinput.VmTypeConfig{
							{Name: "fake-vm-type-1"},
						},
					},
				}

				result := stemcell.Apply(input)
				Expect(result).To(Equal(&bftinput.Input{
					CloudConfig: bftinput.CloudConfig{
						VmTypes: []bftinput.VmTypeConfig{
							{Name: "fake-vm-type-1"},
						},
					},
					Stemcells: []bftinput.StemcellConfig{
						{Alias: "default", OS: "toronto-os", Version: "1"},
					},
				}))
			})
		})

		Context("when input has resource pools", func() {
			It("adds stemcell each resource pool", func() {
				input := &bftinput.Input{
					CloudConfig: bftinput.CloudConfig{
						ResourcePools: []bftinput.ResourcePoolConfig{
							{Name: "fake-vm-type-1"},
						},
					},
				}

				result := stemcell.Apply(input)
				Expect(result).To(Equal(&bftinput.Input{
					CloudConfig: bftinput.CloudConfig{
						ResourcePools: []bftinput.ResourcePoolConfig{
							{
								Name: "fake-vm-type-1",
								Stemcell: bftinput.StemcellConfig{
									OS: "toronto-os", Version: "1",
								},
							},
						},
					},
				}))
			})
		})
	})

	Context("when definition is name", func() {
		BeforeEach(func() {
			stemcell = NewStemcell("name")
		})
	})
})
