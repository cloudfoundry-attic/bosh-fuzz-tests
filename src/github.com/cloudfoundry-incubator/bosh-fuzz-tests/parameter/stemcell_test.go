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
						{Alias: "stemcell-1", OS: "toronto-os", Version: "1"},
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
			stemcell = NewStemcell("name", []string{"1"})
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
						{Alias: "stemcell-1", Name: "ubuntu-stemcell", Version: "1"},
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
									Name: "ubuntu-stemcell", Version: "1",
								},
							},
						},
					},
				}))
			})
		})
	})

	Context("with multiple vm types and jobs", func() {
		BeforeEach(func() {
			rand.Seed(32)
			stemcell = NewStemcell("name", []string{"1", "2"})
		})

		It("generates stemcell version for each vm type and assigns stemcell to corresponding job", func() {
			input := &bftinput.Input{
				Jobs: []bftinput.Job{
					{
						Name:   "fake-job-1",
						VmType: "fake-vm-type-1",
					},
					{
						Name:   "fake-job-2",
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

			result := stemcell.Apply(input)
			Expect(result).To(Equal(&bftinput.Input{
				Jobs: []bftinput.Job{
					{
						Name:     "fake-job-1",
						VmType:   "fake-vm-type-1",
						Stemcell: "stemcell-1",
					},
					{
						Name:     "fake-job-2",
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
