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
  canaries: 2
  canary_watch_time: 4000
  max_in_flight: 1
  update_watch_time: 20

jobs:{{ range .Jobs }}
- name: {{ .Name }}
  instances: {{ .Instances }}{{ if .VmType }}
  vm_type: {{ .VmType }}{{ end }}{{ if .ResourcePool }}
  resource_pool: {{ .ResourcePool }}{{ end }}{{ if .PersistentDiskPool }}
  persistent_disk_pool: {{ .PersistentDiskPool }}{{ else if .PersistentDiskType }}
  persistent_disk_type: {{ .PersistentDiskType }}{{ else if .PersistentDiskSize }}
  persistent_disk: {{ .PersistentDiskSize }}{{ end }}{{ if .Stemcell }}
  stemcell: {{ .Stemcell }}{{ end }}{{ if .MigratedFrom }}
  migrated_from:{{ range .MigratedFrom }}
  - name: {{ .Name }}{{ if .AvailabilityZone }}
    az: {{ .AvailabilityZone }}{{ end }}{{ end }}{{ end }}{{ if .AvailabilityZones }}
  azs:{{ range .AvailabilityZones }}
  - {{ . }}{{ end }}{{ end }}
  templates:{{ range .Templates }}
  - name: {{ .Name }}
    release: foo-release{{ end }}
  networks:{{ range .Networks }}
  - name: {{ .Name }}{{ if .DefaultDNSnGW }}
    default: [dns, gateway]{{ end }}{{ if .StaticIps }}
    static_ips:{{ range .StaticIps }}
    - {{ . }}{{ end }}{{ end }}{{ end }}{{ end }}
`
