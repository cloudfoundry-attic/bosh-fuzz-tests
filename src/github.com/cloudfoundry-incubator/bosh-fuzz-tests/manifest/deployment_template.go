package manifest

var DeploymentTemplate = `---
name: foo-deployment

director_uuid: {{ .DirectorUUID }}

releases:
- name: foo-release
  version: latest

jobs:
- name: {{ .Name }}
  instances: {{ .Instances }}
  vm_type: default
  stemcell: foo-stemcell{{ if .AvailabilityZones }}
  availability_zones:{{ range .AvailabilityZones }}
  - {{ . }}{{ end }}{{ end }}
  templates:
  - name: foo-template
    release: foo-release
  networks: [{name: default}]
`
