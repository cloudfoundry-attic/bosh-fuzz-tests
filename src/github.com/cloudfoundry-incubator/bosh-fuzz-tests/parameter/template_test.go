package parameter_test

import (
	"math/rand"

	bftinput "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"

	. "github.com/cloudfoundry-incubator/bosh-fuzz-tests/parameter"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Job", func() {
	var (
		job Parameter
	)

	BeforeEach(func() {
		job = NewJob([][]string{[]string{"foo", "bar"}, []string{"simple"}})
	})

	It("assigns randomly picked job to instance group", func() {
		rand.Seed(64)
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

		result := job.Apply(input, bftinput.Input{})

		Expect(result).To(Equal(bftinput.Input{
			InstanceGroups: []bftinput.InstanceGroup{
				{
					Name: "fake-instance-group-1",
					Jobs: []bftinput.Job{
						{Name: "simple"},
					},
				},
				{
					Name: "fake-instance-group-2",
					Jobs: []bftinput.Job{
						{Name: "foo"},
						{Name: "bar"},
					},
				},
			},
		}))
	})
})
