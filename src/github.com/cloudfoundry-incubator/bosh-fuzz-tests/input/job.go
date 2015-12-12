package input

import "reflect"

type Job struct {
	Name               string
	Instances          int
	AvailabilityZones  []string
	PersistentDiskSize int
	PersistentDiskPool string
	PersistentDiskType string
	Networks           []JobNetworkConfig
	MigratedFrom       []MigratedFromConfig
	VmType             string
	ResourcePool       string
	Stemcell           string
	Templates          []Template
}

func (j Job) IsEqual(other Job) bool {
	return reflect.DeepEqual(j, other)
}

type Template struct {
	Name string
}

type JobNetworkConfig struct {
	Name          string
	DefaultDNSnGW bool
	StaticIps     []string
}

type MigratedFromConfig struct {
	Name             string
	AvailabilityZone string
}
