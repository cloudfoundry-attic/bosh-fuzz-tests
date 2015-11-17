package deployment

var CloudTemplate = `---{{ if .AvailabilityZones }}
azs:{{ range .AvailabilityZones }}
- name: {{ . }}
  cloud_properties: {}{{ end }}{{ end }}

networks:
- name: default
  subnets:
  - range: "192.168.1.0/24"
    gateway: "192.168.1.1"
    dns: ["192.168.1.1", "192.168.1.2"]
    static: ["192.168.1.10-192.168.1.30"]
    reserved: []
    cloud_properties: {}{{ if .AvailabilityZones }}
    azs:{{ range .AvailabilityZones }}
    - {{ . }}{{ end }}{{ end }}

compilation:
  workers: 1
  network: default
  cloud_properties: {}{{ if .AvailabilityZones }}
  az: {{ index .AvailabilityZones 0 }}{{ end }}

vm_types:
- name: default
  cloud_properties: {}{{ if eq .PersistentDiskDefinition "disk_pool" }}

disk_pools:
- name: fast-disks
  disk_size: {{ .PersistentDiskSize }}
  cloud_properties: {}{{ end }}{{ if eq .PersistentDiskDefinition "disk_type" }}

disk_types:
- name: fast-disks
  disk_size: {{ .PersistentDiskSize }}
  cloud_properties: {}{{ end }}
`
