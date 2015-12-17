package parameter_test

import (
	bftinput "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"

	fakebftdecider "github.com/cloudfoundry-incubator/bosh-fuzz-tests/decider/fakes"
	fakebftnamegen "github.com/cloudfoundry-incubator/bosh-fuzz-tests/name_generator/fakes"

	. "github.com/cloudfoundry-incubator/bosh-fuzz-tests/parameter"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"math/rand"
)

var _ = Describe("CloudProperties", func() {
	var (
		cloudProperties   Parameter
		fakeNameGenerator fakebftnamegen.FakeNameGenerator
	)

	BeforeEach(func() {
		fakeNameGenerator = fakebftnamegen.FakeNameGenerator{
			Names: []string{"steve", "alvin", "jack", "bob", "anderson", "robinson", "hook", "molimo"},
		}
	})

	Context("Adds random cloud properties to input", func() {
		It("fuzz AvailabilityZones", func() {
			rand.Seed(64)
			fakeReuseDecider := &fakebftdecider.FakeDecider{}
			cloudProperties = NewCloudProperties([]int{2}, &fakeNameGenerator, fakeReuseDecider)

			input := bftinput.Input{
				CloudConfig: bftinput.CloudConfig{
					AvailabilityZones: []bftinput.AvailabilityZone{
						{
							Name:            "z1",
							CloudProperties: map[string]string{},
						},
					},
				},
			}
			result := cloudProperties.Apply(input, bftinput.Input{})

			Expect(result).To(Equal(bftinput.Input{
				CloudConfig: bftinput.CloudConfig{
					AvailabilityZones: []bftinput.AvailabilityZone{
						{
							Name: "z1",
							CloudProperties: map[string]string{
								"steve": "alvin",
								"jack":  "bob",
							},
						},
					},
					Compilation: bftinput.CompilationConfig{
						CloudProperties: map[string]string{
							"anderson": "robinson",
							"hook":     "molimo",
						},
					},
				},
			}))
		})

		It("fuzz VM types", func() {
			rand.Seed(64)
			fakeReuseDecider := &fakebftdecider.FakeDecider{}
			cloudProperties = NewCloudProperties([]int{2}, &fakeNameGenerator, fakeReuseDecider)

			input := bftinput.Input{
				CloudConfig: bftinput.CloudConfig{
					VmTypes: []bftinput.VmTypeConfig{
						{
							Name:            "vm1",
							CloudProperties: map[string]string{},
						},
					},
				},
			}
			result := cloudProperties.Apply(input, bftinput.Input{})

			Expect(result).To(Equal(bftinput.Input{
				CloudConfig: bftinput.CloudConfig{
					VmTypes: []bftinput.VmTypeConfig{
						{
							Name: "vm1",
							CloudProperties: map[string]string{
								"steve": "alvin",
								"jack":  "bob",
							},
						},
					},
					Compilation: bftinput.CompilationConfig{
						CloudProperties: map[string]string{
							"anderson": "robinson",
							"hook":     "molimo",
						},
					},
				},
			}))
		})

		It("fuzz Disk Pools", func() {
			rand.Seed(64)
			fakeReuseDecider := &fakebftdecider.FakeDecider{}
			cloudProperties = NewCloudProperties([]int{2}, &fakeNameGenerator, fakeReuseDecider)

			input := bftinput.Input{
				CloudConfig: bftinput.CloudConfig{
					PersistentDiskPools: []bftinput.DiskConfig{
						{
							Name:            "z1",
							CloudProperties: map[string]string{},
						},
					},
				},
			}
			result := cloudProperties.Apply(input, bftinput.Input{})

			Expect(result).To(Equal(bftinput.Input{
				CloudConfig: bftinput.CloudConfig{
					PersistentDiskPools: []bftinput.DiskConfig{
						{
							Name: "z1",
							CloudProperties: map[string]string{
								"steve": "alvin",
								"jack":  "bob",
							},
						},
					},
					Compilation: bftinput.CompilationConfig{
						CloudProperties: map[string]string{
							"anderson": "robinson",
							"hook":     "molimo",
						},
					},
				},
			}))
		})

		It("fuzz Disk Types", func() {
			rand.Seed(64)
			fakeReuseDecider := &fakebftdecider.FakeDecider{}
			cloudProperties = NewCloudProperties([]int{2}, &fakeNameGenerator, fakeReuseDecider)

			input := bftinput.Input{
				CloudConfig: bftinput.CloudConfig{
					PersistentDiskTypes: []bftinput.DiskConfig{
						{
							Name:            "z1",
							CloudProperties: map[string]string{},
						},
					},
				},
			}
			result := cloudProperties.Apply(input, bftinput.Input{})

			Expect(result).To(Equal(bftinput.Input{
				CloudConfig: bftinput.CloudConfig{
					PersistentDiskTypes: []bftinput.DiskConfig{
						{
							Name: "z1",
							CloudProperties: map[string]string{
								"steve": "alvin",
								"jack":  "bob",
							},
						},
					},
					Compilation: bftinput.CompilationConfig{
						CloudProperties: map[string]string{
							"anderson": "robinson",
							"hook":     "molimo",
						},
					},
				},
			}))
		})

		It("fuzz Resource Pool", func() {
			rand.Seed(64)
			fakeReuseDecider := &fakebftdecider.FakeDecider{}
			cloudProperties = NewCloudProperties([]int{2}, &fakeNameGenerator, fakeReuseDecider)

			input := bftinput.Input{
				CloudConfig: bftinput.CloudConfig{
					ResourcePools: []bftinput.ResourcePoolConfig{
						{
							Name:            "z1",
							CloudProperties: map[string]string{},
						},
					},
				},
			}
			result := cloudProperties.Apply(input, bftinput.Input{})

			Expect(result).To(Equal(bftinput.Input{
				CloudConfig: bftinput.CloudConfig{
					ResourcePools: []bftinput.ResourcePoolConfig{
						{
							Name: "z1",
							CloudProperties: map[string]string{
								"steve": "alvin",
								"jack":  "bob",
							},
						},
					},
					Compilation: bftinput.CompilationConfig{
						CloudProperties: map[string]string{
							"anderson": "robinson",
							"hook":     "molimo",
						},
					},
				},
			}))
		})

		It("fuzz Compilation", func() {
			rand.Seed(64)
			fakeReuseDecider := &fakebftdecider.FakeDecider{}
			cloudProperties = NewCloudProperties([]int{2}, &fakeNameGenerator, fakeReuseDecider)

			input := bftinput.Input{
				CloudConfig: bftinput.CloudConfig{
					Compilation: bftinput.CompilationConfig{
						CloudProperties: map[string]string{},
					},
				},
			}
			result := cloudProperties.Apply(input, bftinput.Input{})

			Expect(result).To(Equal(bftinput.Input{
				CloudConfig: bftinput.CloudConfig{
					Compilation: bftinput.CompilationConfig{
						CloudProperties: map[string]string{
							"steve": "alvin",
							"jack":  "bob",
						},
					},
				},
			}))
		})
	})

	It("reuses previous cloud properties", func() {
		rand.Seed(64)
		fakeReuseDecider := &fakebftdecider.FakeDecider{true}
		cloudProperties = NewCloudProperties([]int{2}, &fakeNameGenerator, fakeReuseDecider)

		previousInput := bftinput.Input{
			CloudConfig: bftinput.CloudConfig{
				AvailabilityZones: []bftinput.AvailabilityZone{
					{
						Name: "z1",
						CloudProperties: map[string]string{
							"foo":  "bar",
							"blah": "doug",
						},
					},
				},
			},
		}

		input := bftinput.Input{
			CloudConfig: bftinput.CloudConfig{
				AvailabilityZones: []bftinput.AvailabilityZone{
					{
						Name:            "z1",
						CloudProperties: map[string]string{},
					},
				},
			},
		}
		result := cloudProperties.Apply(input, previousInput)

		Expect(result).To(Equal(bftinput.Input{
			CloudConfig: bftinput.CloudConfig{
				AvailabilityZones: []bftinput.AvailabilityZone{
					{
						Name: "z1",
						CloudProperties: map[string]string{
							"foo":  "bar",
							"blah": "doug",
						},
					},
				},
			},
		}))
	})
})
