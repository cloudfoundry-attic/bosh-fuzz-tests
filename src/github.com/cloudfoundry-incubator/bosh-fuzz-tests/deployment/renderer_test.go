package deployment_test

import (
	. "github.com/cloudfoundry-incubator/bosh-fuzz-tests/deployment"
	bftinput "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"
	faksesys "github.com/cloudfoundry/bosh-utils/system/fakes"
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
			InstanceGroups: []bftinput.InstanceGroup{
				{
					Name:               "foo-instance-group",
					Instances:          5,
					Lifecycle:          "service",
					AvailabilityZones:  []string{"z1", "z2"},
					PersistentDiskSize: 100,
					VmType:             "default",
					Stemcell:           "default",
					Jobs: []bftinput.Job{
						{Name: "simple"},
					},
					Networks: []bftinput.InstanceGroupNetworkConfig{
						{
							Name:          "default",
							StaticIps:     []string{"192.168.1.5"},
							DefaultDNSnGW: true,
						},
					},
				},
				{
					Name:               "bar-instance-group",
					Instances:          2,
					Lifecycle:          "errand",
					AvailabilityZones:  []string{"z3", "z4"},
					PersistentDiskPool: "fast-disks",
					VmType:             "default",
					Stemcell:           "default",
					Jobs: []bftinput.Job{
						{Name: "simple"},
					},
					Networks: []bftinput.InstanceGroupNetworkConfig{
						{
							Name:          "default",
							DefaultDNSnGW: true,
						},
					},
					MigratedFrom: []bftinput.MigratedFromConfig{
						{Name: "baz-instance-group", AvailabilityZone: "z5"},
					},
				},
			},
			Update: bftinput.UpdateConfig{
				Canaries:    1,
				MaxInFlight: 3,
				Serial:      "true",
			},
			CloudConfig: bftinput.CloudConfig{
				AvailabilityZones: []bftinput.AvailabilityZone{
					{
						Name: "z1",
						CloudProperties: map[string]string{
							"foo": "bar",
							"baz": "qux",
						},
					},
					{Name: "z2"},
					{Name: "z3"},
					{Name: "z4"},
				},
				PersistentDiskPools: []bftinput.DiskConfig{
					{
						Name: "fast-disks",
						Size: 200,
						CloudProperties: map[string]string{
							"foo": "bar",
							"baz": "qux",
						},
					},
				},
				VmTypes: []bftinput.VmTypeConfig{
					{
						Name: "default",
						CloudProperties: map[string]string{
							"foo": "bar",
							"baz": "qux",
						},
					},
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
								CloudProperties: map[string]string{
									"foo": "bar",
									"baz": "qux",
								},
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
				Compilation: bftinput.CompilationConfig{
					Network:          "default",
					AvailabilityZone: "z1",
					NumberOfWorkers:  3,
					CloudProperties: map[string]string{
						"foo": "bar",
						"baz": "qux",
					},
				},
			},
			Stemcells: []bftinput.StemcellConfig{
				{
					Alias:   "default",
					OS:      "toronto-os",
					Version: "1",
				},
			},
			Variables: []bftinput.Variable{
				{
					Name: "var1",
					Type: "ssh",
				},
				{
					Name: "var2",
					Type: "rsa",
				},
				{
					Name: "var3",
					Type: "password",
				},
				{
					Name:    "var4",
					Type:    "certificate",
					Options: map[string]interface{}{"is_ca": true},
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
  canaries: 1
  canary_watch_time: 4000
  max_in_flight: 3
  update_watch_time: 20
  serial: true

jobs:
- name: foo-instance-group
  instances: 5
  lifecycle: service
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
- name: bar-instance-group
  instances: 2
  lifecycle: errand
  vm_type: default
  persistent_disk_pool: fast-disks
  stemcell: default
  migrated_from:
  - name: baz-instance-group
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

variables:
- name: var1
  type: ssh
- name: var2
  type: rsa
- name: var3
  type: password
- name: var4
  type: certificate
  options:
    is_ca: true
`

		manifestContents, err := fs.ReadFileString(manifestPath)
		Expect(err).ToNot(HaveOccurred())
		Expect(manifestContents).To(Equal(expectedManifestContents))

		expectedCloudConfigContents := `---
azs:
- name: z1
  cloud_properties:
    baz: qux
    foo: bar
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
  - cloud_properties:
      baz: qux
      foo: bar
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
  workers: 3
  network: default
  cloud_properties:
    baz: qux
    foo: bar
  az: z1

vm_types:
- name: default
  cloud_properties:
    baz: qux
    foo: bar

disk_pools:
- name: fast-disks
  disk_size: 200
  cloud_properties:
    baz: qux
    foo: bar
`

		cloudConfigContents, err := fs.ReadFileString(cloudConfigPath)
		Expect(err).ToNot(HaveOccurred())
		Expect(cloudConfigContents).To(Equal(expectedCloudConfigContents))
	})

	Context("when AvailabilityZone is nil", func() {
		It("does not specify az key in manifest", func() {
			input := bftinput.Input{
				DirectorUUID: "d820eb0d-13db-4777-8c9b-7a9bc55e3628",
				InstanceGroups: []bftinput.InstanceGroup{
					{
						Name:      "foo-instance-group",
						Instances: 5,
						Networks:  []bftinput.InstanceGroupNetworkConfig{{Name: "default"}},
						Jobs: []bftinput.Job{
							{Name: "simple"},
						},
					},
				},
				Update: bftinput.UpdateConfig{
					Canaries:    1,
					MaxInFlight: 3,
					Serial:      "not_specified",
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
					Compilation: bftinput.CompilationConfig{
						Network:         "default",
						NumberOfWorkers: 3,
					},
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
  canaries: 1
  canary_watch_time: 4000
  max_in_flight: 3
  update_watch_time: 20

jobs:
- name: foo-instance-group
  instances: 5
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
  workers: 3
  network: default
  cloud_properties: {}
`

			cloudConfigContents, err := fs.ReadFileString(cloudConfigPath)
			Expect(err).ToNot(HaveOccurred())
			Expect(cloudConfigContents).To(Equal(expectedCloudConfigContents))
		})
	})

	It("uses the disk pool specified for instance group", func() {
		input := bftinput.Input{
			DirectorUUID: "d820eb0d-13db-4777-8c9b-7a9bc55e3628",
			InstanceGroups: []bftinput.InstanceGroup{
				{
					Name:               "foo-instance-group",
					Instances:          5,
					PersistentDiskPool: "fast-disks",
					Networks:           []bftinput.InstanceGroupNetworkConfig{{Name: "default"}},
					Jobs: []bftinput.Job{
						{Name: "simple"},
					},
				},
			},
			Update: bftinput.UpdateConfig{
				Canaries:    3,
				MaxInFlight: 5,
				Serial:      "true",
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
				Compilation: bftinput.CompilationConfig{
					Network:         "default",
					NumberOfWorkers: 3,
				},
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
  canaries: 3
  canary_watch_time: 4000
  max_in_flight: 5
  update_watch_time: 20
  serial: true

jobs:
- name: foo-instance-group
  instances: 5
  persistent_disk_pool: fast-disks
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
  workers: 3
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
			InstanceGroups: []bftinput.InstanceGroup{
				{
					Name:               "foo-instance-group",
					Instances:          5,
					PersistentDiskType: "fast-disks",
					Networks:           []bftinput.InstanceGroupNetworkConfig{{Name: "default"}},
					Jobs: []bftinput.Job{
						{Name: "simple"},
					},
				},
			},
			Update: bftinput.UpdateConfig{
				Canaries:    2,
				MaxInFlight: 4,
				Serial:      "false",
			},
			CloudConfig: bftinput.CloudConfig{
				PersistentDiskTypes: []bftinput.DiskConfig{
					{
						Name: "fast-disks",
						Size: 100,
						CloudProperties: map[string]string{
							"foo": "bar",
							"baz": "qux",
						},
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
				Compilation: bftinput.CompilationConfig{
					Network:         "default",
					NumberOfWorkers: 3,
				},
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
  max_in_flight: 4
  update_watch_time: 20
  serial: false

jobs:
- name: foo-instance-group
  instances: 5
  persistent_disk_type: fast-disks
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
  workers: 3
  network: default
  cloud_properties: {}

disk_types:
- name: fast-disks
  disk_size: 100
  cloud_properties:
    baz: qux
    foo: bar
`

		cloudConfigContents, err := fs.ReadFileString(cloudConfigPath)
		Expect(err).ToNot(HaveOccurred())
		Expect(cloudConfigContents).To(Equal(expectedCloudConfigContents))
	})

	It("uses the resource pool specified for instance group", func() {
		input := bftinput.Input{
			DirectorUUID: "d820eb0d-13db-4777-8c9b-7a9bc55e3628",
			InstanceGroups: []bftinput.InstanceGroup{
				{
					Name:               "foo-instance-group",
					Instances:          5,
					ResourcePool:       "foo-pool",
					PersistentDiskPool: "fast-disks",
					Networks:           []bftinput.InstanceGroupNetworkConfig{{Name: "default"}},
					Jobs: []bftinput.Job{
						{Name: "simple"},
					},
				},
			},
			Update: bftinput.UpdateConfig{
				Canaries:    3,
				MaxInFlight: 5,
				Serial:      "true",
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
				ResourcePools: []bftinput.ResourcePoolConfig{
					{
						Name: "foo-pool",
						Stemcell: bftinput.StemcellConfig{
							Name:    "foo",
							Version: "1",
						},
						CloudProperties: map[string]string{
							"foo": "bar",
							"baz": "qux",
						},
					},
				},
				Compilation: bftinput.CompilationConfig{
					Network:         "default",
					NumberOfWorkers: 3,
				},
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
  canaries: 3
  canary_watch_time: 4000
  max_in_flight: 5
  update_watch_time: 20
  serial: true

jobs:
- name: foo-instance-group
  instances: 5
  resource_pool: foo-pool
  persistent_disk_pool: fast-disks
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
  workers: 3
  network: default
  cloud_properties: {}

resource_pools:
- name: foo-pool
  stemcell:
    version: 1
    name: foo
  cloud_properties:
    baz: qux
    foo: bar

disk_pools:
- name: fast-disks
  disk_size: 100
  cloud_properties: {}
`

		cloudConfigContents, err := fs.ReadFileString(cloudConfigPath)
		Expect(err).ToNot(HaveOccurred())
		Expect(cloudConfigContents).To(Equal(expectedCloudConfigContents))
	})
})
