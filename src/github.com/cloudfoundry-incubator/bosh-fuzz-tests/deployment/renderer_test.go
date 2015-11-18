package deployment_test

import (
	faksesys "github.com/cloudfoundry/bosh-utils/system/fakes"

	. "github.com/cloudfoundry-incubator/bosh-fuzz-tests/deployment"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Manifest/Renderer", func() {
	var (
		renderer        Renderer
		fs              *faksesys.FakeFileSystem
		manifestPath    string
		cloudConfigPath string
	)

	BeforeEach(func() {
		fs = faksesys.NewFakeFileSystem()
		renderer = NewRenderer(fs)
		manifestPath = "manifest-path"
		cloudConfigPath = "cloud-config-path"
	})

	It("creates manifest based on input values", func() {
		input := Input{
			DirectorUUID: "d820eb0d-13db-4777-8c9b-7a9bc55e3628",
			Jobs: []Job{
				{
					Name:               "foo-job",
					Instances:          5,
					AvailabilityZones:  []string{"z1", "z2"},
					PersistentDiskSize: 100,
					Network:            "default",
				},
				{
					Name:               "bar-job",
					Instances:          2,
					AvailabilityZones:  []string{"z3", "z4"},
					PersistentDiskPool: "fast-disks",
					Network:            "default",
					MigratedFrom:       []string{"baz-job"},
				},
			},
			CloudConfig: CloudConfig{
				AvailabilityZones: []string{"z1", "z2", "z3", "z4"},
				PersistentDiskPools: []DiskConfig{
					{
						Name: "fast-disks",
						Size: 200,
					},
				},
			},
		}

		err := renderer.Render(input, manifestPath, cloudConfigPath)
		Expect(err).ToNot(HaveOccurred())
		expectedManifestContents := `---
name: foo-deployment

director_uuid: d820eb0d-13db-4777-8c9b-7a9bc55e3628

stemcells:
- alias: default
  os: toronto-os
  version: 1

releases:
- name: foo-release
  version: latest

update:
  canaries: 2
  canary_watch_time: 4000
  max_in_flight: 1
  update_watch_time: 20

jobs:
- name: foo-job
  instances: 5
  vm_type: default
  persistent_disk: 100
  stemcell: default
  azs:
  - z1
  - z2
  templates:
  - name: simple
    release: foo-release
  networks:
  - name: default
- name: bar-job
  instances: 2
  vm_type: default
  persistent_disk_pool: fast-disks
  stemcell: default
  migrated_from:
  - name: baz-job
  azs:
  - z3
  - z4
  templates:
  - name: simple
    release: foo-release
  networks:
  - name: default
`

		manifestContents, err := fs.ReadFileString(manifestPath)
		Expect(err).ToNot(HaveOccurred())
		Expect(manifestContents).To(Equal(expectedManifestContents))

		expectedCloudConfigContents := `---
azs:
- name: z1
  cloud_properties: {}
- name: z2
  cloud_properties: {}
- name: z3
  cloud_properties: {}
- name: z4
  cloud_properties: {}

networks:
- name: default
  subnets:
  - range: "192.168.1.0/24"
    gateway: "192.168.1.1"
    dns: ["192.168.1.1", "192.168.1.2"]
    static: ["192.168.1.10-192.168.1.30"]
    reserved: []
    cloud_properties: {}
    azs:
    - z1
    - z2
    - z3
    - z4
- name: no-az
  subnets:
  - range: "192.168.1.0/24"
    gateway: "192.168.1.1"
    dns: ["192.168.1.1", "192.168.1.2"]
    static: ["192.168.1.10-192.168.1.30"]
    reserved: []
    cloud_properties: {}

compilation:
  workers: 1
  network: default
  cloud_properties: {}
  az: z1

vm_types:
- name: default
  cloud_properties: {}

disk_pools:
- name: fast-disks
  disk_size: 200
  cloud_properties: {}
`

		cloudConfigContents, err := fs.ReadFileString(cloudConfigPath)
		Expect(err).ToNot(HaveOccurred())
		Expect(cloudConfigContents).To(Equal(expectedCloudConfigContents))
	})

	Context("when AvailabilityZone is nil", func() {
		It("does not specify az key in manifest", func() {
			input := Input{
				DirectorUUID: "d820eb0d-13db-4777-8c9b-7a9bc55e3628",
				Jobs: []Job{
					{
						Name:              "foo-job",
						Instances:         5,
						AvailabilityZones: nil,
						Network:           "default",
					},
				},
			}

			err := renderer.Render(input, manifestPath, cloudConfigPath)
			Expect(err).ToNot(HaveOccurred())
			expectedManifestContents := `---
name: foo-deployment

director_uuid: d820eb0d-13db-4777-8c9b-7a9bc55e3628

stemcells:
- alias: default
  os: toronto-os
  version: 1

releases:
- name: foo-release
  version: latest

update:
  canaries: 2
  canary_watch_time: 4000
  max_in_flight: 1
  update_watch_time: 20

jobs:
- name: foo-job
  instances: 5
  vm_type: default
  stemcell: default
  templates:
  - name: simple
    release: foo-release
  networks:
  - name: default
`

			manifestContents, err := fs.ReadFileString(manifestPath)
			Expect(err).ToNot(HaveOccurred())
			Expect(manifestContents).To(Equal(expectedManifestContents))

			expectedCloudConfigContents := `---

networks:
- name: default
  subnets:
  - range: "192.168.1.0/24"
    gateway: "192.168.1.1"
    dns: ["192.168.1.1", "192.168.1.2"]
    static: ["192.168.1.10-192.168.1.30"]
    reserved: []
    cloud_properties: {}
- name: no-az
  subnets:
  - range: "192.168.1.0/24"
    gateway: "192.168.1.1"
    dns: ["192.168.1.1", "192.168.1.2"]
    static: ["192.168.1.10-192.168.1.30"]
    reserved: []
    cloud_properties: {}

compilation:
  workers: 1
  network: default
  cloud_properties: {}

vm_types:
- name: default
  cloud_properties: {}
`

			cloudConfigContents, err := fs.ReadFileString(cloudConfigPath)
			Expect(err).ToNot(HaveOccurred())
			Expect(cloudConfigContents).To(Equal(expectedCloudConfigContents))
		})
	})

	It("uses the disk pool specified for job", func() {
		input := Input{
			DirectorUUID: "d820eb0d-13db-4777-8c9b-7a9bc55e3628",
			Jobs: []Job{
				{
					Name:               "foo-job",
					Instances:          5,
					AvailabilityZones:  nil,
					PersistentDiskPool: "fast-disks",
					Network:            "default",
				},
			},
			CloudConfig: CloudConfig{
				PersistentDiskPools: []DiskConfig{
					{
						Name: "fast-disks",
						Size: 100,
					},
				},
			},
		}

		err := renderer.Render(input, manifestPath, cloudConfigPath)
		Expect(err).ToNot(HaveOccurred())
		expectedManifestContents := `---
name: foo-deployment

director_uuid: d820eb0d-13db-4777-8c9b-7a9bc55e3628

stemcells:
- alias: default
  os: toronto-os
  version: 1

releases:
- name: foo-release
  version: latest

update:
  canaries: 2
  canary_watch_time: 4000
  max_in_flight: 1
  update_watch_time: 20

jobs:
- name: foo-job
  instances: 5
  vm_type: default
  persistent_disk_pool: fast-disks
  stemcell: default
  templates:
  - name: simple
    release: foo-release
  networks:
  - name: default
`

		manifestContents, err := fs.ReadFileString(manifestPath)
		Expect(err).ToNot(HaveOccurred())
		Expect(manifestContents).To(Equal(expectedManifestContents))

		expectedCloudConfigContents := `---

networks:
- name: default
  subnets:
  - range: "192.168.1.0/24"
    gateway: "192.168.1.1"
    dns: ["192.168.1.1", "192.168.1.2"]
    static: ["192.168.1.10-192.168.1.30"]
    reserved: []
    cloud_properties: {}
- name: no-az
  subnets:
  - range: "192.168.1.0/24"
    gateway: "192.168.1.1"
    dns: ["192.168.1.1", "192.168.1.2"]
    static: ["192.168.1.10-192.168.1.30"]
    reserved: []
    cloud_properties: {}

compilation:
  workers: 1
  network: default
  cloud_properties: {}

vm_types:
- name: default
  cloud_properties: {}

disk_pools:
- name: fast-disks
  disk_size: 100
  cloud_properties: {}
`

		cloudConfigContents, err := fs.ReadFileString(cloudConfigPath)
		Expect(err).ToNot(HaveOccurred())
		Expect(cloudConfigContents).To(Equal(expectedCloudConfigContents))
	})

	It("uses the disk type", func() {
		input := Input{
			DirectorUUID: "d820eb0d-13db-4777-8c9b-7a9bc55e3628",
			Jobs: []Job{
				{
					Name:               "foo-job",
					Instances:          5,
					AvailabilityZones:  nil,
					PersistentDiskType: "fast-disks",
					Network:            "default",
				},
			},
			CloudConfig: CloudConfig{
				PersistentDiskTypes: []DiskConfig{
					{
						Name: "fast-disks",
						Size: 100,
					},
				},
			},
		}

		err := renderer.Render(input, manifestPath, cloudConfigPath)
		Expect(err).ToNot(HaveOccurred())
		expectedManifestContents := `---
name: foo-deployment

director_uuid: d820eb0d-13db-4777-8c9b-7a9bc55e3628

stemcells:
- alias: default
  os: toronto-os
  version: 1

releases:
- name: foo-release
  version: latest

update:
  canaries: 2
  canary_watch_time: 4000
  max_in_flight: 1
  update_watch_time: 20

jobs:
- name: foo-job
  instances: 5
  vm_type: default
  persistent_disk_type: fast-disks
  stemcell: default
  templates:
  - name: simple
    release: foo-release
  networks:
  - name: default
`

		manifestContents, err := fs.ReadFileString(manifestPath)
		Expect(err).ToNot(HaveOccurred())
		Expect(manifestContents).To(Equal(expectedManifestContents))

		expectedCloudConfigContents := `---

networks:
- name: default
  subnets:
  - range: "192.168.1.0/24"
    gateway: "192.168.1.1"
    dns: ["192.168.1.1", "192.168.1.2"]
    static: ["192.168.1.10-192.168.1.30"]
    reserved: []
    cloud_properties: {}
- name: no-az
  subnets:
  - range: "192.168.1.0/24"
    gateway: "192.168.1.1"
    dns: ["192.168.1.1", "192.168.1.2"]
    static: ["192.168.1.10-192.168.1.30"]
    reserved: []
    cloud_properties: {}

compilation:
  workers: 1
  network: default
  cloud_properties: {}

vm_types:
- name: default
  cloud_properties: {}

disk_types:
- name: fast-disks
  disk_size: 100
  cloud_properties: {}
`

		cloudConfigContents, err := fs.ReadFileString(cloudConfigPath)
		Expect(err).ToNot(HaveOccurred())
		Expect(cloudConfigContents).To(Equal(expectedCloudConfigContents))
	})
})
