package deployment

var CloudTemplate = `---{{ if .CloudConfig.AvailabilityZones }}
azs:{{ range .CloudConfig.AvailabilityZones }}
- name: {{ .Name }}
  cloud_properties: {{ if .CloudProperties }}{{ range $key, $value := .CloudProperties }}
  {{ $key }}: {{ $value }}{{ end }}{{ else }}{}{{ end }}{{ end }}{{ end }}

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
  az: {{ .CloudConfig.CompilationAvailabilityZone }}{{ end }}{{ if .CloudConfig.VmTypes }}

vm_types:{{ range .CloudConfig.VmTypes }}
- name: {{ .Name }}
  cloud_properties: {}{{ end }}{{ end }}{{ if .CloudConfig.ResourcePools }}

resource_pools:{{ range .CloudConfig.ResourcePools }}
- name: {{ .Name }}
  stemcell:
    version: {{ .Stemcell.Version }}{{ if .Stemcell.Name }}
    name: {{ .Stemcell.Name }}{{ end }}{{ if .Stemcell.Alias }}
    alias: {{ .Stemcell.Alias }}{{ end }}{{ if .Stemcell.OS }}
    os: {{ .Stemcell.OS }}{{ end }}
  cloud_properties: {}{{ end }}{{ end }}{{ if .CloudConfig.PersistentDiskPools }}

disk_pools:{{ range .CloudConfig.PersistentDiskPools }}
- name: {{ .Name }}
  disk_size: {{ .Size }}
  cloud_properties: {}{{ end }}{{ end }}{{ if .CloudConfig.PersistentDiskTypes }}

disk_types:{{ range .CloudConfig.PersistentDiskTypes }}
- name: {{ .Name }}
  disk_size: {{ .Size }}
  cloud_properties: {}{{ end }}{{ end }}
`
