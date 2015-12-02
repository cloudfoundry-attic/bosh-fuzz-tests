package deployment_test

import (
	bftinput "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"
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
		input := bftinput.Input{
			DirectorUUID: "d820eb0d-13db-4777-8c9b-7a9bc55e3628",
			Jobs: []bftinput.Job{
				{
					Name:               "foo-job",
					Instances:          5,
					AvailabilityZones:  []string{"z1", "z2"},
					PersistentDiskSize: 100,
					VmType:             "default",
					Networks: []bftinput.JobNetworkConfig{
						{
							Name:          "default",
							StaticIps:     []string{"192.168.1.5"},
							DefaultDNSnGW: true,
						},
					},
				},
				{
					Name:               "bar-job",
					Instances:          2,
					AvailabilityZones:  []string{"z3", "z4"},
					PersistentDiskPool: "fast-disks",
					VmType:             "default",
					Networks: []bftinput.JobNetworkConfig{
						{
							Name:          "default",
							DefaultDNSnGW: true,
						},
					},
					MigratedFrom: []bftinput.MigratedFromConfig{
						{Name: "baz-job", AvailabilityZone: "z5"},
					},
				},
			},
			CloudConfig: bftinput.CloudConfig{
				AvailabilityZones: []string{"z1", "z2", "z3", "z4"},
				PersistentDiskPools: []bftinput.DiskConfig{
					{
						Name: "fast-disks",
						Size: 200,
					},
				},
				VmTypes: []bftinput.VmTypeConfig{
					{Name: "default"},
				},
				Networks: []bftinput.NetworkConfig{
					{
						Name: "default",
						Type: "manual",
						Subnets: []bftinput.SubnetConfig{
							{
								IpPool: &bftinput.IpPool{
									IpRange: "192.168.1.0/24",
									Gateway: "192.168.1.254",
									Reserved: []string{
										"192.168.1.11",
										"192.168.1.120",
										"192.168.1.186-192.168.1.234",
									},
									Static: []string{
										"192.168.1.5",
									},
								},
								AvailabilityZones: []string{"z1", "z2", "z3", "z4"},
							},
						},
					},
					{
						Name: "no-az",
						Type: "dynamic",
						Subnets: []bftinput.SubnetConfig{
							{},
						},
					},
				},
				CompilationNetwork:          "default",
				CompilationAvailabilityZone: "z1",
			},
			Stemcells: []bftinput.StemcellConfig{
				{
					Alias:   "default",
					OS:      "toronto-os",
					Version: "1",
				},
			},
		}

		err := renderer.Render(input, manifestPath, cloudConfigPath)
		Expect(err).ToNot(HaveOccurred())
		expectedManifestContents := `---
name: foo-deployment

director_uuid: d820eb0d-13db-4777-8c9b-7a9bc55e3628

stemcells:
- version: 1
  alias: default
  os: toronto-os

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
    default: [dns, gateway]
    static_ips:
    - 192.168.1.5
- name: bar-job
  instances: 2
  vm_type: default
  persistent_disk_pool: fast-disks
  stemcell: default
  migrated_from:
  - name: baz-job
    az: z5
  azs:
  - z3
  - z4
  templates:
  - name: simple
    release: foo-release
  networks:
  - name: default
    default: [dns, gateway]
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
  type: manual
  subnets:
  - cloud_properties: {}
    dns: ["8.8.8.8"]
    range: 192.168.1.0/24
    gateway: 192.168.1.254
    static:
    - 192.168.1.5
    reserved:
    - 192.168.1.11
    - 192.168.1.120
    - 192.168.1.186-192.168.1.234
    azs:
    - z1
    - z2
    - z3
    - z4
- name: no-az
  type: dynamic
  subnets:
  - cloud_properties: {}
    dns: ["8.8.8.8"]

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
			input := bftinput.Input{
				DirectorUUID: "d820eb0d-13db-4777-8c9b-7a9bc55e3628",
				Jobs: []bftinput.Job{
					{
						Name:      "foo-job",
						Instances: 5,
						Networks:  []bftinput.JobNetworkConfig{{Name: "default"}},
					},
				},
				CloudConfig: bftinput.CloudConfig{
					Networks: []bftinput.NetworkConfig{
						{
							Name: "default",
							Type: "manual",
							Subnets: []bftinput.SubnetConfig{
								{},
							},
						},
						{
							Name: "no-az",
							Type: "dynamic",
							Subnets: []bftinput.SubnetConfig{
								{},
							},
						},
					},
					CompilationNetwork: "default",
				},
			}

			err := renderer.Render(input, manifestPath, cloudConfigPath)
			Expect(err).ToNot(HaveOccurred())
			expectedManifestContents := `---
name: foo-deployment

director_uuid: d820eb0d-13db-4777-8c9b-7a9bc55e3628

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
  type: manual
  subnets:
  - cloud_properties: {}
    dns: ["8.8.8.8"]
- name: no-az
  type: dynamic
  subnets:
  - cloud_properties: {}
    dns: ["8.8.8.8"]

compilation:
  workers: 1
  network: default
  cloud_properties: {}
`

			cloudConfigContents, err := fs.ReadFileString(cloudConfigPath)
			Expect(err).ToNot(HaveOccurred())
			Expect(cloudConfigContents).To(Equal(expectedCloudConfigContents))
		})
	})

	It("uses the disk pool specified for job", func() {
		input := bftinput.Input{
			DirectorUUID: "d820eb0d-13db-4777-8c9b-7a9bc55e3628",
			Jobs: []bftinput.Job{
				{
					Name:               "foo-job",
					Instances:          5,
					PersistentDiskPool: "fast-disks",
					Networks:           []bftinput.JobNetworkConfig{{Name: "default"}},
				},
			},
			CloudConfig: bftinput.CloudConfig{
				PersistentDiskPools: []bftinput.DiskConfig{
					{
						Name: "fast-disks",
						Size: 100,
					},
				},
				Networks: []bftinput.NetworkConfig{
					{
						Name: "default",
						Type: "manual",
						Subnets: []bftinput.SubnetConfig{
							{},
						},
					},
					{
						Name: "no-az",
						Type: "dynamic",
						Subnets: []bftinput.SubnetConfig{
							{},
						},
					},
				},
				CompilationNetwork: "default",
			},
		}

		err := renderer.Render(input, manifestPath, cloudConfigPath)
		Expect(err).ToNot(HaveOccurred())
		expectedManifestContents := `---
name: foo-deployment

director_uuid: d820eb0d-13db-4777-8c9b-7a9bc55e3628

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
  type: manual
  subnets:
  - cloud_properties: {}
    dns: ["8.8.8.8"]
- name: no-az
  type: dynamic
  subnets:
  - cloud_properties: {}
    dns: ["8.8.8.8"]

compilation:
  workers: 1
  network: default
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
		input := bftinput.Input{
			DirectorUUID: "d820eb0d-13db-4777-8c9b-7a9bc55e3628",
			Jobs: []bftinput.Job{
				{
					Name:               "foo-job",
					Instances:          5,
					PersistentDiskType: "fast-disks",
					Networks:           []bftinput.JobNetworkConfig{{Name: "default"}},
				},
			},
			CloudConfig: bftinput.CloudConfig{
				PersistentDiskTypes: []bftinput.DiskConfig{
					{
						Name: "fast-disks",
						Size: 100,
					},
				},
				Networks: []bftinput.NetworkConfig{
					{
						Name: "default",
						Type: "manual",
						Subnets: []bftinput.SubnetConfig{
							{},
						},
					},
					{
						Name: "no-az",
						Type: "dynamic",
						Subnets: []bftinput.SubnetConfig{
							{},
						},
					},
				},
				CompilationNetwork: "default",
			},
		}

		err := renderer.Render(input, manifestPath, cloudConfigPath)
		Expect(err).ToNot(HaveOccurred())
		expectedManifestContents := `---
name: foo-deployment

director_uuid: d820eb0d-13db-4777-8c9b-7a9bc55e3628

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
  type: manual
  subnets:
  - cloud_properties: {}
    dns: ["8.8.8.8"]
- name: no-az
  type: dynamic
  subnets:
  - cloud_properties: {}
    dns: ["8.8.8.8"]

compilation:
  workers: 1
  network: default
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
