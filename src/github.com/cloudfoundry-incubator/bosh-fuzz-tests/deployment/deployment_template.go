package deployment

var DeploymentTemplate = `---
name: foo-deployment

director_uuid: {{ .DirectorUUID }}

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

jobs:{{ range .Jobs }}
- name: {{ .Name }}
  instances: {{ .Instances }}
  vm_type: default{{ if .PersistentDiskPool }}
  persistent_disk_pool: {{ .PersistentDiskPool }}{{ else if .PersistentDiskType }}
  persistent_disk_type: {{ .PersistentDiskType }}{{ else if .PersistentDiskSize }}
  persistent_disk: {{ .PersistentDiskSize }}{{ end }}
  stemcell: default{{ if .MigratedFrom }}
  migrated_from:{{ range .MigratedFrom }}
  - name: {{ .Name }}{{ if .AvailabilityZone }}
    az: {{ .AvailabilityZone }}{{ end }}{{ end }}{{ end }}{{ if .AvailabilityZones }}
  azs:{{ range .AvailabilityZones }}
  - {{ . }}{{ end }}{{ end }}
  templates:
  - name: simple
    release: foo-release
  networks:{{ range .Networks }}
  - name: {{ .Name }}{{ if .DefaultDNSnGW }}
    default: [dns, gateway]{{ end }}{{ end }}{{ end }}
`
