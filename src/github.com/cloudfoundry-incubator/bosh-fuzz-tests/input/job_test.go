package input_test

import (
	. "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Job", func() {
	Describe("HasPersistentDisk", func() {
		Context("when Job has persistent disk pool", func() {
			It("should return true", func() {
				job := Job{
					Name:               "foo",
					PersistentDiskPool: "disk-pool",
				}
				Expect(job.HasPersistentDisk()).To(BeTrue())
			})
		})

		Context("when Job has persistent disk type", func() {
			It("should return true", func() {
				job := Job{
					Name:               "foo",
					PersistentDiskType: "disk-type",
				}
				Expect(job.HasPersistentDisk()).To(BeTrue())
			})
		})

		Context("when Job has persistent disk size", func() {
			It("should return true", func() {
				job := Job{
					Name:               "foo",
					PersistentDiskSize: 1024,
				}
				Expect(job.HasPersistentDisk()).To(BeTrue())
			})
		})

		Context("when Job has no mention of persistent disk", func() {
			It("should return false", func() {
				job := Job{
					Name: "foo",
				}
				Expect(job.HasPersistentDisk()).To(BeFalse())
			})
		})
	})
})
