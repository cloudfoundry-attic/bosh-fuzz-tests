package deployment

var DeploymentTemplate = `---
name: foo-deployment

director_uuid: {{ .DirectorUUID }}{{ if .Stemcells }}

stemcells:{{ range .Stemcells }}
- version: {{ .Version }}{{ if .Name }}
  name: {{ .Name }}{{ end }}{{ if .Alias }}
  alias: {{ .Alias }}{{ end }}{{ if .OS }}
  os: {{ .OS }}{{ end }}{{ end }}{{ end }}

releases:
- name: foo-release
  version: latest

update:
  canaries: {{ .Update.Canaries }}
  canary_watch_time: 4000
  max_in_flight: {{ .Update.MaxInFlight }}
  update_watch_time: 20{{ if ne .Update.Serial "not_specified" }}
  serial: {{ .Update.Serial }}{{ end }}

instance_groups:{{ range .InstanceGroups }}
- name: {{ .Name }}
  instances: {{ .Instances }}{{ if .Lifecycle }}
  lifecycle: {{ .Lifecycle }}{{ end }}{{ if .VmType }}
  vm_type: {{ .VmType }}{{ end }}{{ if .PersistentDiskType }}
  persistent_disk_type: {{ .PersistentDiskType }}{{ else if .PersistentDiskSize }}
  persistent_disk: {{ .PersistentDiskSize }}{{ end }}{{ if .Stemcell }}
  stemcell: {{ .Stemcell }}{{ end }}{{ if .MigratedFrom }}
  migrated_from:{{ range .MigratedFrom }}
  - name: {{ .Name }}{{ if .AvailabilityZone }}
    az: {{ .AvailabilityZone }}{{ end }}{{ end }}{{ end }}{{ if .AvailabilityZones }}
  azs:{{ range .AvailabilityZones }}
  - {{ . }}{{ end }}{{ end }}
  jobs:{{ range .Jobs }}
  - name: {{ .Name }}
    release: foo-release{{ end }}
  networks:{{ range .Networks }}
  - name: {{ .Name }}{{ if .DefaultDNSnGW }}
    default: [dns, gateway]{{ end }}{{ if .StaticIps }}
    static_ips:{{ range .StaticIps }}
    - {{ . }}{{ end }}{{ end }}{{ end }}{{ end }}{{ if .Variables }}

variables:{{ range .Variables }}
- name: {{ .Name }}
  type: {{ .Type }}{{ if .Options }}
  options:{{ range $key, $value := .Options }}
    {{ $key }}: {{ $value }}{{ end }}{{ end }}{{ end }}{{ end }}{{ if not .CloudConfig.AvailabilityZones }}

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
  az: {{ .CloudConfig.Compilation.AvailabilityZone }}{{ end }}{{ end }}
`
