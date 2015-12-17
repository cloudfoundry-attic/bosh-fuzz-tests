package deployment

var CloudTemplate = `---{{ if .CloudConfig.AvailabilityZones }}
azs:{{ range .CloudConfig.AvailabilityZones }}
- name: {{ .Name }}
  cloud_properties:{{ if .CloudProperties }}{{ range $key, $value := .CloudProperties }}
    {{ $key }}: {{ $value }}{{ end }}{{ else }} {}{{ end }}{{ end }}{{ end }}

networks:{{ range .CloudConfig.Networks }}
- name: {{ .Name }}
  type: {{ .Type }}{{ if .CloudProperties }}
  cloud_properties:{{ range $key, $value := .CloudProperties }}
    {{ $key }}: {{ $value }}{{ end }}{{ end }}{{ if .Subnets }}
  subnets:{{ range .Subnets }}
  - cloud_properties:{{ if .CloudProperties }}{{ range $key, $value := .CloudProperties }}
      {{ $key }}: {{ $value }}{{ end }}{{ else }} {}{{ end }}
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
  workers: {{ .CloudConfig.Compilation.NumberOfWorkers }}
  network: {{ .CloudConfig.Compilation.Network }}
  cloud_properties:{{ if .CloudConfig.Compilation.CloudProperties }}{{ range $key, $value := .CloudConfig.Compilation.CloudProperties }}
    {{ $key }}: {{ $value }}{{ end }}{{ else }} {}{{ end }}{{ if .CloudConfig.Compilation.AvailabilityZone }}
  az: {{ .CloudConfig.Compilation.AvailabilityZone }}{{ end }}{{ if .CloudConfig.VmTypes }}

vm_types:{{ range .CloudConfig.VmTypes }}
- name: {{ .Name }}
  cloud_properties:{{ if .CloudProperties }}{{ range $key, $value := .CloudProperties }}
    {{ $key }}: {{ $value }}{{ end }}{{ else }} {}{{ end }}{{ end }}{{ end }}{{ if .CloudConfig.ResourcePools }}

resource_pools:{{ range .CloudConfig.ResourcePools }}
- name: {{ .Name }}
  stemcell:
    version: {{ .Stemcell.Version }}{{ if .Stemcell.Name }}
    name: {{ .Stemcell.Name }}{{ end }}{{ if .Stemcell.Alias }}
    alias: {{ .Stemcell.Alias }}{{ end }}{{ if .Stemcell.OS }}
    os: {{ .Stemcell.OS }}{{ end }}
  cloud_properties:{{ if .CloudProperties }}{{ range $key, $value := .CloudProperties }}
    {{ $key }}: {{ $value }}{{ end }}{{ else }} {}{{ end }}{{ end }}{{ end }}{{ if .CloudConfig.PersistentDiskPools }}

disk_pools:{{ range .CloudConfig.PersistentDiskPools }}
- name: {{ .Name }}
  disk_size: {{ .Size }}
  cloud_properties:{{ if .CloudProperties }}{{ range $key, $value := .CloudProperties }}
    {{ $key }}: {{ $value }}{{ end }}{{ else }} {}{{ end }}{{ end }}{{ end }}{{ if .CloudConfig.PersistentDiskTypes }}

disk_types:{{ range .CloudConfig.PersistentDiskTypes }}
- name: {{ .Name }}
  disk_size: {{ .Size }}
  cloud_properties:{{ if .CloudProperties }}{{ range $key, $value := .CloudProperties }}
    {{ $key }}: {{ $value }}{{ end }}{{ else }} {}{{ end }}{{ end }}{{ end }}
`
