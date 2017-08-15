package parameter_test

import (
	bftinput "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"
	. "github.com/cloudfoundry-incubator/bosh-fuzz-tests/parameter"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Lifecycle", func() {
	var (
		lifecycle Parameter
	)

	BeforeEach(func() {
		lifecycle = NewLifecycle()
	})

	It("randomly chooses errand or service", func() {
		input := bftinput.Input{Jobs: []bftinput.Job{{}}}

		Eventually(func() string {
			newInput := lifecycle.Apply(input, bftinput.Input{})
			return newInput.Jobs[0].Lifecycle
		}).Should(Equal("errand"))

		Eventually(func() string {
			newInput := lifecycle.Apply(input, bftinput.Input{})
			return newInput.Jobs[0].Lifecycle
		}).Should(Equal("service"))
	})

	It("adds a lifecycle to all jobs", func() {
		input := bftinput.Input{Jobs: []bftinput.Job{{}, {}, {}, {}}}
		newInput := lifecycle.Apply(input, bftinput.Input{})
		for _, job := range newInput.Jobs {
			Expect([]string{"errand", "service"}).To(ContainElement(job.Lifecycle))
		}
	})

	Context("when job has a persistent disk pool", func() {
		It("always sets lifecycle to service, never to errand", func() {
			input := bftinput.Input{Jobs: []bftinput.Job{{PersistentDiskPool: "diskpool"}}}
			Consistently(func() string {
				return lifecycle.Apply(input, bftinput.Input{}).Jobs[0].Lifecycle
			}).Should(Equal("service"))
		})
	})

	Context("when job has a persistent disk type", func() {
		It("always sets lifecycle to service, never to errand", func() {
			input := bftinput.Input{Jobs: []bftinput.Job{{PersistentDiskType: "diskpool"}}}
			Consistently(func() string {
				return lifecycle.Apply(input, bftinput.Input{}).Jobs[0].Lifecycle
			}).Should(Equal("service"))
		})
	})

	Context("when job has a persistent disk size", func() {
		It("always sets lifecycle to service, never to errand", func() {
			input := bftinput.Input{Jobs: []bftinput.Job{{PersistentDiskSize: 100}}}
			Consistently(func() string {
				return lifecycle.Apply(input, bftinput.Input{}).Jobs[0].Lifecycle
			}).Should(Equal("service"))
		})
	})

	Context("when previous job has a persistent disk pool", func() {
		It("always sets lifecycle to service, never to errand", func() {
			input := bftinput.Input{Jobs: []bftinput.Job{{Name: "foo"}}}
			previousInput := bftinput.Input{Jobs: []bftinput.Job{{Name: "foo", PersistentDiskPool: "diskpool"}}}
			Consistently(func() string {
				return lifecycle.Apply(input, previousInput).Jobs[0].Lifecycle
			}).Should(Equal("service"))
		})
	})

	Context("when previous job has a persistent disk type", func() {
		It("always sets lifecycle to service, never to errand", func() {
			input := bftinput.Input{Jobs: []bftinput.Job{{Name: "foo"}}}
			previousInput := bftinput.Input{Jobs: []bftinput.Job{{Name: "foo", PersistentDiskType: "diskpool"}}}
			Consistently(func() string {
				return lifecycle.Apply(input, previousInput).Jobs[0].Lifecycle
			}).Should(Equal("service"))
		})
	})

	Context("when previous job has a persistent disk size", func() {
		It("always sets lifecycle to service, never to errand", func() {
			input := bftinput.Input{Jobs: []bftinput.Job{{Name: "foo"}}}
			previousInput := bftinput.Input{Jobs: []bftinput.Job{{Name: "foo", PersistentDiskSize: 100}}}
			Consistently(func() string {
				return lifecycle.Apply(input, previousInput).Jobs[0].Lifecycle
			}).Should(Equal("service"))
		})
	})

	Context("when previous migrated from job has a persistent disk pool", func() {
		It("always sets lifecycle to service, never to errand", func() {
			input := bftinput.Input{
				Jobs: []bftinput.Job{
					{MigratedFrom: []bftinput.MigratedFromConfig{{Name: "foo"}}},
				},
			}
			previousInput := bftinput.Input{Jobs: []bftinput.Job{{Name: "foo", PersistentDiskPool: "diskpool"}}}
			Consistently(func() string {
				return lifecycle.Apply(input, previousInput).Jobs[0].Lifecycle
			}).Should(Equal("service"))
		})
	})

	Context("when previous migrated from job has a persistent disk type", func() {
		It("always sets lifecycle to service, never to errand", func() {
			input := bftinput.Input{
				Jobs: []bftinput.Job{
					{MigratedFrom: []bftinput.MigratedFromConfig{{Name: "foo"}}},
				},
			}
			previousInput := bftinput.Input{Jobs: []bftinput.Job{{Name: "foo", PersistentDiskType: "diskpool"}}}
			Consistently(func() string {
				return lifecycle.Apply(input, previousInput).Jobs[0].Lifecycle
			}).Should(Equal("service"))
		})
	})

	Context("when previous migrated from job has a persistent disk size", func() {
		It("always sets lifecycle to service, never to errand", func() {
			input := bftinput.Input{
				Jobs: []bftinput.Job{
					{MigratedFrom: []bftinput.MigratedFromConfig{{Name: "foo"}}},
				},
			}
			previousInput := bftinput.Input{Jobs: []bftinput.Job{{Name: "foo", PersistentDiskSize: 100}}}
			Consistently(func() string {
				return lifecycle.Apply(input, previousInput).Jobs[0].Lifecycle
			}).Should(Equal("service"))
		})
	})
})
