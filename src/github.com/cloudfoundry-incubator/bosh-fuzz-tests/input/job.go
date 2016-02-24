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

func (j Job) HasPersistentDisk() bool {
	return j.PersistentDiskSize != 0 || j.PersistentDiskPool != "" || j.PersistentDiskType != ""
}

func (j Job) FindNetworkByName(networkName string) (JobNetworkConfig, bool) {
	for _, network := range j.Networks {
		if network.Name == networkName {
			return network, true
		}
	}
	return JobNetworkConfig{}, false
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
