package parameter_test

import (
	"math/rand"

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
			stemcell = NewStemcell("os", []string{"1"})
		})

		Context("when input has vm types", func() {
			It("adds stemcells to the input", func() {
				input := bftinput.Input{
					CloudConfig: bftinput.CloudConfig{
						VmTypes: []bftinput.VmTypeConfig{
							{Name: "fake-vm-type-1"},
						},
					},
				}

				result := stemcell.Apply(input, bftinput.Input{})
				Expect(result).To(Equal(bftinput.Input{
					CloudConfig: bftinput.CloudConfig{
						VmTypes: []bftinput.VmTypeConfig{
							{Name: "fake-vm-type-1"},
						},
					},
					Stemcells: []bftinput.StemcellConfig{
						{Alias: "stemcell-1", OS: "toronto-os", Version: "1"},
					},
				}))
			})
		})
	})

	Context("when definition is name", func() {
		BeforeEach(func() {
			stemcell = NewStemcell("name", []string{"1"})
		})

		Context("when input has vm types", func() {
			It("adds stemcells to the input", func() {
				input := bftinput.Input{
					CloudConfig: bftinput.CloudConfig{
						VmTypes: []bftinput.VmTypeConfig{
							{Name: "fake-vm-type-1"},
						},
					},
				}

				result := stemcell.Apply(input, bftinput.Input{})
				Expect(result).To(Equal(bftinput.Input{
					CloudConfig: bftinput.CloudConfig{
						VmTypes: []bftinput.VmTypeConfig{
							{Name: "fake-vm-type-1"},
						},
					},
					Stemcells: []bftinput.StemcellConfig{
						{Alias: "stemcell-1", Name: "ubuntu-stemcell", Version: "1"},
					},
				}))
			})
		})
	})

	Context("with multiple vm types and instance groups", func() {
		BeforeEach(func() {
			rand.Seed(32)
			stemcell = NewStemcell("name", []string{"1", "2"})
		})

		It("generates stemcell version for each vm type and assigns stemcell to corresponding instance group", func() {
			input := bftinput.Input{
				InstanceGroups: []bftinput.InstanceGroup{
					{
						Name:   "fake-instance-group-1",
						VmType: "fake-vm-type-1",
					},
					{
						Name:   "fake-instance-group-2",
						VmType: "fake-vm-type-2",
					},
				},
				CloudConfig: bftinput.CloudConfig{
					VmTypes: []bftinput.VmTypeConfig{
						{Name: "fake-vm-type-1"},
						{Name: "fake-vm-type-2"},
					},
				},
			}

			result := stemcell.Apply(input, bftinput.Input{})
			Expect(result).To(Equal(bftinput.Input{
				InstanceGroups: []bftinput.InstanceGroup{
					{
						Name:     "fake-instance-group-1",
						VmType:   "fake-vm-type-1",
						Stemcell: "stemcell-1",
					},
					{
						Name:     "fake-instance-group-2",
						VmType:   "fake-vm-type-2",
						Stemcell: "stemcell-2",
					},
				},
				CloudConfig: bftinput.CloudConfig{
					VmTypes: []bftinput.VmTypeConfig{
						{Name: "fake-vm-type-1"},
						{Name: "fake-vm-type-2"},
					},
				},
				Stemcells: []bftinput.StemcellConfig{
					{Alias: "stemcell-1", Name: "ubuntu-stemcell", Version: "1"},
					{Alias: "stemcell-2", Name: "ubuntu-stemcell", Version: "2"},
				},
			}))
		})
	})
})
