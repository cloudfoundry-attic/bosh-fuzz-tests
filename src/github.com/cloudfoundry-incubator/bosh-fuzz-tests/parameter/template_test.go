package parameter_test

import (
	"math/rand"

	bftinput "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"

	. "github.com/cloudfoundry-incubator/bosh-fuzz-tests/parameter"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Template", func() {
	var (
		template Parameter
	)

	BeforeEach(func() {
		template = NewTemplate([][]string{[]string{"foo", "bar"}, []string{"simple"}})
	})

	It("assigns randomly picked template to instance group", func() {
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

		result := template.Apply(input, bftinput.Input{})

		Expect(result).To(Equal(bftinput.Input{
			InstanceGroups: []bftinput.InstanceGroup{
				{
					Name: "fake-instance-group-1",
					Templates: []bftinput.Template{
						{Name: "simple"},
					},
				},
				{
					Name: "fake-instance-group-2",
					Templates: []bftinput.Template{
						{Name: "foo"},
						{Name: "bar"},
					},
				},
			},
		}))
	})
})
