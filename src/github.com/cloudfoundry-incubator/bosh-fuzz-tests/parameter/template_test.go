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

	It("assigns randomly picked template to job", func() {
		rand.Seed(64)
		input := bftinput.Input{
			Jobs: []bftinput.Job{
				{
					Name: "fake-job-1",
				},
				{
					Name: "fake-job-2",
				},
			},
		}

		result := template.Apply(input)

		Expect(result).To(Equal(bftinput.Input{
			Jobs: []bftinput.Job{
				{
					Name: "fake-job-1",
					Templates: []bftinput.Template{
						{Name: "simple"},
					},
				},
				{
					Name: "fake-job-2",
					Templates: []bftinput.Template{
						{Name: "foo"},
						{Name: "bar"},
					},
				},
			},
		}))
	})
})
