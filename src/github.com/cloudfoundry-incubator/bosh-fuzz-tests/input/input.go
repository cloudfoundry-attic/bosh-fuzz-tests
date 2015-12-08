package input

import "reflect"

type Input struct {
	DirectorUUID string
	Jobs         []Job
	CloudConfig  CloudConfig
	Stemcells    []StemcellConfig
}

func (i Input) FindJobByName(jobName string) (Job, bool) {
	for _, job := range i.Jobs {
		if job.Name == jobName {
			return job, true
		}
	}
	return Job{}, false
}

func (i Input) FindAzByName(azName string) (AvailabilityZone, bool) {
	for _, az := range i.CloudConfig.AvailabilityZones {
		if az.Name == azName {
			return az, true
		}
	}
	return AvailabilityZone{}, false
}

func (i Input) FindDiskPoolByName(diskName string) (DiskConfig, bool) {
	for _, disk := range i.CloudConfig.PersistentDiskPools {
		if disk.Name == diskName {
			return disk, true
		}
	}
	return DiskConfig{}, false
}

func (i Input) FindDiskTypeByName(diskName string) (DiskConfig, bool) {
	for _, disk := range i.CloudConfig.PersistentDiskTypes {
		if disk.Name == diskName {
			return disk, true
		}
	}
	return DiskConfig{}, false
}

func (i Input) FindNetworkByName(networkName string) (NetworkConfig, bool) {
	for _, network := range i.CloudConfig.Networks {
		if network.Name == networkName {
			return network, true
		}
	}
	return NetworkConfig{}, false
}

func (i Input) FindResourcePoolByName(resourcePoolName string) (ResourcePoolConfig, bool) {
	for _, resourcePool := range i.CloudConfig.ResourcePools {
		if resourcePool.Name == resourcePoolName {
			return resourcePool, true
		}
	}
	return ResourcePoolConfig{}, false
}

type CloudConfig struct {
	AvailabilityZones           []AvailabilityZone
	PersistentDiskPools         []DiskConfig
	PersistentDiskTypes         []DiskConfig
	Networks                    []NetworkConfig
	CompilationNetwork          string
	CompilationAvailabilityZone string
	VmTypes                     []VmTypeConfig
	ResourcePools               []ResourcePoolConfig
}

type DiskConfig struct {
	Name string
	Size int
}

func (d DiskConfig) IsEqual(other DiskConfig) bool {
	return d == other
}

type AvailabilityZone struct {
	Name            string
	CloudProperties map[string]interface{}
}

func (a AvailabilityZone) IsEqual(other AvailabilityZone) bool {
	return reflect.DeepEqual(a, other)
}

type VmTypeConfig struct {
	Name string
}

type ResourcePoolConfig struct {
	Name     string
	Stemcell StemcellConfig
}

func (r ResourcePoolConfig) IsEqual(other ResourcePoolConfig) bool {
	return reflect.DeepEqual(r, other)
}

type StemcellConfig struct {
	Name    string
	OS      string
	Version string
	Alias   string
}

type MigratedFromConfig struct {
	Name             string
	AvailabilityZone string
}

type NetworkConfig struct {
	Name    string
	Type    string
	Subnets []SubnetConfig
}

func (n NetworkConfig) IsEqual(other NetworkConfig) bool {
	return reflect.DeepEqual(n, other)
}

type SubnetConfig struct {
	AvailabilityZones []string
	IpPool            *IpPool
}

type JobNetworkConfig struct {
	Name          string
	DefaultDNSnGW bool
	StaticIps     []string
}
