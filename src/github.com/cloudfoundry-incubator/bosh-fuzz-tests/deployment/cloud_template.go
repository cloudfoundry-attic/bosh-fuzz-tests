package deployment

var CloudTemplate = `---{{ if .CloudConfig.AvailabilityZones }}
azs:{{ range .CloudConfig.AvailabilityZones }}
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
    cloud_properties: {}{{ if .CloudConfig.AvailabilityZones }}
    azs:{{ range .CloudConfig.AvailabilityZones }}
    - {{ . }}{{ end }}{{ end }}
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
  cloud_properties: {}{{ if .CloudConfig.AvailabilityZones }}
  az: {{ index .CloudConfig.AvailabilityZones 0 }}{{ end }}

vm_types:
- name: default
  cloud_properties: {}{{ if .CloudConfig.PersistentDiskPools }}

disk_pools:{{ range .CloudConfig.PersistentDiskPools }}
- name: {{ .Name }}
  disk_size: {{ .Size }}
  cloud_properties: {}{{ end }}{{ end }}{{ if .CloudConfig.PersistentDiskTypes }}

disk_types:{{ range .CloudConfig.PersistentDiskTypes }}
- name: {{ .Name }}
  disk_size: {{ .Size }}
  cloud_properties: {}{{ end }}{{ end }}
`
