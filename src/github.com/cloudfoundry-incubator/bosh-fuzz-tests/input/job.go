package input

import "reflect"

type InstanceGroup struct {
	Name               string
	Instances          int
	AvailabilityZones  []string
	PersistentDiskSize int
	PersistentDiskPool string
	PersistentDiskType string
	Networks           []InstanceGroupNetworkConfig
	MigratedFrom       []MigratedFromConfig
	VmType             string
	ResourcePool       string
	Stemcell           string
	Jobs               []Job
	Lifecycle          string
}

func (j InstanceGroup) IsEqual(other InstanceGroup) bool {
	return reflect.DeepEqual(j, other)
}

func (j InstanceGroup) HasPersistentDisk() bool {
	return j.PersistentDiskSize != 0 || j.PersistentDiskPool != "" || j.PersistentDiskType != ""
}

func (j InstanceGroup) FindNetworkByName(networkName string) (InstanceGroupNetworkConfig, bool) {
	for _, network := range j.Networks {
		if network.Name == networkName {
			return network, true
		}
	}
	return InstanceGroupNetworkConfig{}, false
}

type Job struct {
	Name string
}

type InstanceGroupNetworkConfig struct {
	Name          string
	DefaultDNSnGW bool
	StaticIps     []string
}

type MigratedFromConfig struct {
	Name             string
	AvailabilityZone string
}
