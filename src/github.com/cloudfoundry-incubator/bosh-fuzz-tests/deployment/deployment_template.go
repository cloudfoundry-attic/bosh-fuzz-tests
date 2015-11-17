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

jobs:
- name: {{ .Name }}
  instances: {{ .Instances }}
  vm_type: default{{ if .PersistentDiskSize }}
  persistent_disk: {{ .PersistentDiskSize }}{{ end }}
  stemcell: default{{ if .AvailabilityZones }}
  azs:{{ range .AvailabilityZones }}
  - {{ . }}{{ end }}{{ end }}
  templates:
  - name: simple
    release: foo-release
  networks: [{name: default}]
`
