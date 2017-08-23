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
		input := bftinput.Input{InstanceGroups: []bftinput.InstanceGroup{{}}}

		Eventually(func() string {
			newInput := lifecycle.Apply(input, bftinput.Input{})
			return newInput.InstanceGroups[0].Lifecycle
		}).Should(Equal("errand"))

		Eventually(func() string {
			newInput := lifecycle.Apply(input, bftinput.Input{})
			return newInput.InstanceGroups[0].Lifecycle
		}).Should(Equal("service"))
	})

	It("does not change underlying slice", func() {
		instanceGroups := []bftinput.InstanceGroup{{}}
		input := bftinput.Input{InstanceGroups: instanceGroups}

		lifecycle.Apply(input, bftinput.Input{})
		Expect(instanceGroups[0].Lifecycle).To(Equal(""))
	})

	It("adds a lifecycle to all instance groups", func() {
		input := bftinput.Input{InstanceGroups: []bftinput.InstanceGroup{{}, {}, {}, {}}}
		newInput := lifecycle.Apply(input, bftinput.Input{})
		for _, instanceGroup := range newInput.InstanceGroups {
			Expect([]string{"errand", "service"}).To(ContainElement(instanceGroup.Lifecycle))
		}
	})

	Context("when instance group has a persistent disk pool", func() {
		It("always sets lifecycle to service, never to errand", func() {
			input := bftinput.Input{InstanceGroups: []bftinput.InstanceGroup{{PersistentDiskPool: "diskpool"}}}
			Consistently(func() string {
				return lifecycle.Apply(input, bftinput.Input{}).InstanceGroups[0].Lifecycle
			}).Should(Equal("service"))
		})
	})

	Context("when instance group has a persistent disk type", func() {
		It("always sets lifecycle to service, never to errand", func() {
			input := bftinput.Input{InstanceGroups: []bftinput.InstanceGroup{{PersistentDiskType: "diskpool"}}}
			Consistently(func() string {
				return lifecycle.Apply(input, bftinput.Input{}).InstanceGroups[0].Lifecycle
			}).Should(Equal("service"))
		})
	})

	Context("when instance group has a persistent disk size", func() {
		It("always sets lifecycle to service, never to errand", func() {
			input := bftinput.Input{InstanceGroups: []bftinput.InstanceGroup{{PersistentDiskSize: 100}}}
			Consistently(func() string {
				return lifecycle.Apply(input, bftinput.Input{}).InstanceGroups[0].Lifecycle
			}).Should(Equal("service"))
		})
	})

	Context("when previous instance group has a persistent disk pool", func() {
		It("always sets lifecycle to service, never to errand", func() {
			input := bftinput.Input{InstanceGroups: []bftinput.InstanceGroup{{Name: "foo"}}}
			previousInput := bftinput.Input{InstanceGroups: []bftinput.InstanceGroup{{Name: "foo", PersistentDiskPool: "diskpool"}}}
			Consistently(func() string {
				return lifecycle.Apply(input, previousInput).InstanceGroups[0].Lifecycle
			}).Should(Equal("service"))
		})
	})

	Context("when previous instance group has a persistent disk type", func() {
		It("always sets lifecycle to service, never to errand", func() {
			input := bftinput.Input{InstanceGroups: []bftinput.InstanceGroup{{Name: "foo"}}}
			previousInput := bftinput.Input{InstanceGroups: []bftinput.InstanceGroup{{Name: "foo", PersistentDiskType: "diskpool"}}}
			Consistently(func() string {
				return lifecycle.Apply(input, previousInput).InstanceGroups[0].Lifecycle
			}).Should(Equal("service"))
		})
	})

	Context("when previous instance group has a persistent disk size", func() {
		It("always sets lifecycle to service, never to errand", func() {
			input := bftinput.Input{InstanceGroups: []bftinput.InstanceGroup{{Name: "foo"}}}
			previousInput := bftinput.Input{InstanceGroups: []bftinput.InstanceGroup{{Name: "foo", PersistentDiskSize: 100}}}
			Consistently(func() string {
				return lifecycle.Apply(input, previousInput).InstanceGroups[0].Lifecycle
			}).Should(Equal("service"))
		})
	})

	Context("when previous migrated from instance group has a persistent disk pool", func() {
		It("always sets lifecycle to service, never to errand", func() {
			input := bftinput.Input{
				InstanceGroups: []bftinput.InstanceGroup{
					{MigratedFrom: []bftinput.MigratedFromConfig{{Name: "foo"}}},
				},
			}
			previousInput := bftinput.Input{InstanceGroups: []bftinput.InstanceGroup{{Name: "foo", PersistentDiskPool: "diskpool"}}}
			Consistently(func() string {
				return lifecycle.Apply(input, previousInput).InstanceGroups[0].Lifecycle
			}).Should(Equal("service"))
		})
	})

	Context("when previous migrated from instance group has a persistent disk type", func() {
		It("always sets lifecycle to service, never to errand", func() {
			input := bftinput.Input{
				InstanceGroups: []bftinput.InstanceGroup{
					{MigratedFrom: []bftinput.MigratedFromConfig{{Name: "foo"}}},
				},
			}
			previousInput := bftinput.Input{InstanceGroups: []bftinput.InstanceGroup{{Name: "foo", PersistentDiskType: "diskpool"}}}
			Consistently(func() string {
				return lifecycle.Apply(input, previousInput).InstanceGroups[0].Lifecycle
			}).Should(Equal("service"))
		})
	})

	Context("when previous migrated from instance group has a persistent disk size", func() {
		It("always sets lifecycle to service, never to errand", func() {
			input := bftinput.Input{
				InstanceGroups: []bftinput.InstanceGroup{
					{MigratedFrom: []bftinput.MigratedFromConfig{{Name: "foo"}}},
				},
			}
			previousInput := bftinput.Input{InstanceGroups: []bftinput.InstanceGroup{{Name: "foo", PersistentDiskSize: 100}}}
			Consistently(func() string {
				return lifecycle.Apply(input, previousInput).InstanceGroups[0].Lifecycle
			}).Should(Equal("service"))
		})
	})
})
