package deployment

var CloudTemplate = `---{{ if .CloudConfig.AvailabilityZones }}
azs:{{ range .CloudConfig.AvailabilityZones }}
- name: {{ . }}
  cloud_properties: {}{{ end }}{{ end }}

networks:{{ range .CloudConfig.Networks }}
- name: {{ .Name }}
  type: {{ .Type }}{{ if .Subnets }}
  subnets:{{ range .Subnets }}
  - cloud_properties: {}
    dns: ["8.8.8.8"]{{ with .IpPool }}{{ if .IpRange }}
    range: {{ .IpRange }}{{ end }}{{ if .Gateway }}
    gateway: {{ .Gateway }}{{ end }}{{ if .Static }}
    static:{{ range .Static }}
    - {{ . }}{{ end }}{{ end }}{{ if .Reserved }}
    reserved:{{ range .Reserved }}
    - {{ . }}{{ end }}{{ end }}{{ end }}{{ if .AvailabilityZones }}
    azs:{{ range .AvailabilityZones }}
    - {{ . }}{{ end }}{{ end }}{{ end }}{{ end }}{{ end }}

compilation:
  workers: 1
  network: {{ .CloudConfig.CompilationNetwork }}
  cloud_properties: {}{{ if .CloudConfig.CompilationAvailabilityZone }}
  az: {{ .CloudConfig.CompilationAvailabilityZone }}{{ end }}

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
