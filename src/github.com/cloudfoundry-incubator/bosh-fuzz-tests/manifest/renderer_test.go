package manifest_test

import (
	faksesys "github.com/cloudfoundry/bosh-utils/system/fakes"

	. "github.com/cloudfoundry-incubator/bosh-fuzz-tests/manifest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Manifest/Renderer", func() {
	var (
		renderer Renderer
		fs       *faksesys.FakeFileSystem
	)

	BeforeEach(func() {
		fs = faksesys.NewFakeFileSystem()
		renderer = NewRenderer(fs)
	})

	It("creates manifest based on input values", func() {
		input := Input{
			Name:              "foo-job",
			DirectorUUID:      "d820eb0d-13db-4777-8c9b-7a9bc55e3628",
			Instances:         5,
			AvailabilityZones: []string{"z1", "z2"},
		}
		manifestPath := "manifest-path"

		err := renderer.Render(input, manifestPath)
		Expect(err).ToNot(HaveOccurred())
		expectedManifestContents := `---
name: foo-deployment

director_uuid: d820eb0d-13db-4777-8c9b-7a9bc55e3628

releases:
- name: foo-release
  version: latest

jobs:
- name: foo-job
  instances: 5
  vm_type: default
  stemcell: foo-stemcell
  availability_zones:
  - z1
  - z2
  templates:
  - name: foo-template
    release: foo-release
  networks: [{name: default}]
`

		manifestContents, err := fs.ReadFileString(manifestPath)
		Expect(err).ToNot(HaveOccurred())
		Expect(manifestContents).To(Equal(expectedManifestContents))
	})

	Context("when AvailabilityZone is nil", func() {
		It("does not specify az key in manifest", func() {
			input := Input{
				Name:              "foo-job",
				DirectorUUID:      "d820eb0d-13db-4777-8c9b-7a9bc55e3628",
				Instances:         5,
				AvailabilityZones: nil,
			}

			manifestPath := "manifest-path"

			err := renderer.Render(input, manifestPath)
			Expect(err).ToNot(HaveOccurred())
			expectedManifestContents := `---
name: foo-deployment

director_uuid: d820eb0d-13db-4777-8c9b-7a9bc55e3628

releases:
- name: foo-release
  version: latest

jobs:
- name: foo-job
  instances: 5
  vm_type: default
  stemcell: foo-stemcell
  templates:
  - name: foo-template
    release: foo-release
  networks: [{name: default}]
`

			manifestContents, err := fs.ReadFileString(manifestPath)
			Expect(err).ToNot(HaveOccurred())
			Expect(manifestContents).To(Equal(expectedManifestContents))
		})
	})
})
