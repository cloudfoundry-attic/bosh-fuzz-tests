package input_test

import (
	. "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("InstanceGroup", func() {
	Describe("HasPersistentDisk", func() {
		Context("when InstanceGroup has persistent disk type", func() {
			It("should return true", func() {
				instanceGroup := InstanceGroup{
					Name:               "foo",
					PersistentDiskType: "disk-type",
				}
				Expect(instanceGroup.HasPersistentDisk()).To(BeTrue())
			})
		})

		Context("when InstanceGroup has persistent disk size", func() {
			It("should return true", func() {
				instanceGroup := InstanceGroup{
					Name:               "foo",
					PersistentDiskSize: 1024,
				}
				Expect(instanceGroup.HasPersistentDisk()).To(BeTrue())
			})
		})

		Context("when InstanceGroup has no mention of persistent disk", func() {
			It("should return false", func() {
				instanceGroup := InstanceGroup{
					Name: "foo",
				}
				Expect(instanceGroup.HasPersistentDisk()).To(BeFalse())
			})
		})
	})
})
