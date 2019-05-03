package parameter_test

import (
	bftinput "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"
	fakebftnamegen "github.com/cloudfoundry-incubator/bosh-fuzz-tests/name_generator/fakes"
	. "github.com/cloudfoundry-incubator/bosh-fuzz-tests/parameter"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("PersistentDisk", func() {
	var (
		persistentDisk Parameter
	)

	Context("when definition is disk_type", func() {
		BeforeEach(func() {
			fakeNameGenerator := &fakebftnamegen.FakeNameGenerator{
				Names: []string{"fake-disk-config"},
			}
			persistentDisk = NewPersistentDisk("disk_type", []int{100}, fakeNameGenerator)
		})

		It("adds disk_types to the input", func() {
			input := bftinput.Input{
				InstanceGroups: []bftinput.InstanceGroup{
					{
						Name: "fake-instance-group",
					},
				},
			}

			result := persistentDisk.Apply(input, bftinput.Input{})

			Expect(result).To(Equal(bftinput.Input{
				InstanceGroups: []bftinput.InstanceGroup{
					{
						Name:               "fake-instance-group",
						PersistentDiskType: "fake-disk-config",
					},
				},
				CloudConfig: bftinput.CloudConfig{
					PersistentDiskTypes: []bftinput.DiskConfig{
						{
							Name: "fake-disk-config",
							Size: 100,
						},
					},
				},
			}))
		})
	})
})
