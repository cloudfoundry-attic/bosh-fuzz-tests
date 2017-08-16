package input

import "reflect"

type Input struct {
	DirectorUUID   string
	InstanceGroups []InstanceGroup
	Update         UpdateConfig
	CloudConfig    CloudConfig
	Stemcells      []StemcellConfig
	Variables      []Variable
}

func (i Input) FindInstanceGroupByName(instanceGroupName string) (InstanceGroup, bool) {
	for _, instanceGroup := range i.InstanceGroups {
		if instanceGroup.Name == instanceGroupName {
			return instanceGroup, true
		}
	}
	return InstanceGroup{}, false
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

func (i Input) FindVmTypeByName(vmTypeName string) (VmTypeConfig, bool) {
	for _, vmType := range i.CloudConfig.VmTypes {
		if vmType.Name == vmTypeName {
			return vmType, true
		}
	}
	return VmTypeConfig{}, false
}

func (i Input) FindSubnetByIpRange(ipRange string) (SubnetConfig, bool) {
	for _, network := range i.CloudConfig.Networks {
		for _, subnet := range network.Subnets {
			if subnet.IpPool != nil && subnet.IpPool.IpRange == ipRange {
				return subnet, true
			}
		}
	}

	return SubnetConfig{}, false
}

func (i Input) FindStemcellByName(stemcellName string) (StemcellConfig, bool) {
	for _, stemcell := range i.Stemcells {
		if stemcell.Name == stemcellName {
			return stemcell, true
		}
	}
	return StemcellConfig{}, false
}

type CloudConfig struct {
	AvailabilityZones   []AvailabilityZone
	PersistentDiskPools []DiskConfig
	PersistentDiskTypes []DiskConfig
	Networks            []NetworkConfig
	Compilation         CompilationConfig
	VmTypes             []VmTypeConfig
	ResourcePools       []ResourcePoolConfig
}

type DiskConfig struct {
	Name            string
	Size            int
	CloudProperties map[string]string
}

func (d DiskConfig) IsEqual(other DiskConfig) bool {
	return reflect.DeepEqual(d, other)
}

type CompilationConfig struct {
	Network          string
	AvailabilityZone string
	NumberOfWorkers  int
	CloudProperties  map[string]string
}

type AvailabilityZone struct {
	Name            string
	CloudProperties map[string]string
}

func (a AvailabilityZone) IsEqual(other AvailabilityZone) bool {
	return reflect.DeepEqual(a, other)
}

type VmTypeConfig struct {
	Name            string
	CloudProperties map[string]string
}

func (v VmTypeConfig) IsEqual(other VmTypeConfig) bool {
	return reflect.DeepEqual(v, other)
}

type ResourcePoolConfig struct {
	Name            string
	Stemcell        StemcellConfig
	CloudProperties map[string]string
}

func (r ResourcePoolConfig) IsEqual(other ResourcePoolConfig) bool {
	return reflect.DeepEqual(r, other)
}

type UpdateConfig struct {
	Canaries    int
	MaxInFlight int
	Serial      string
}

type StemcellConfig struct {
	Name    string
	OS      string
	Version string
	Alias   string
}

func (s StemcellConfig) IsEqual(other StemcellConfig) bool {
	return s.Version == other.Version
}

type NetworkConfig struct {
	Name            string
	Type            string
	Subnets         []SubnetConfig
	CloudProperties map[string]string
}

func (n NetworkConfig) IsEqual(other NetworkConfig) bool {
	return reflect.DeepEqual(n, other)
}

type SubnetConfig struct {
	AvailabilityZones []string
	IpPool            *IpPool
	CloudProperties   map[string]string
}
