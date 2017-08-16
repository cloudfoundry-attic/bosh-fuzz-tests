package parameter_test

import (
	"math/rand"

	bftinput "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"

	. "github.com/cloudfoundry-incubator/bosh-fuzz-tests/parameter"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("AvailabilityZone", func() {
	var (
		az Parameter
	)

	Context("when definition is os", func() {
		BeforeEach(func() {
			rand.Seed(64)

			az = NewAvailabilityZone([][]string{[]string{"z1", "z2"}, []string{"z2", "z3"}})
		})

		It("adds azs to the input", func() {
			input := bftinput.Input{
				InstanceGroups: []bftinput.InstanceGroup{
					{Name: "fake-instance-group-1"},
					{Name: "fake-instance-group-2"},
				},
			}

			result := az.Apply(input, bftinput.Input{})
			Expect(result).To(Equal(bftinput.Input{
				InstanceGroups: []bftinput.InstanceGroup{
					{
						Name:              "fake-instance-group-1",
						AvailabilityZones: []string{"z2", "z3"},
					},
					{
						Name:              "fake-instance-group-2",
						AvailabilityZones: []string{"z1", "z2"},
					},
				},
				CloudConfig: bftinput.CloudConfig{
					AvailabilityZones: []bftinput.AvailabilityZone{
						{Name: "z2"},
						{Name: "z3"},
						{Name: "z1"},
					},
				},
			}))
		})
	})
})
