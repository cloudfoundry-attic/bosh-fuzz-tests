package deployment

type Input struct {
	DirectorUUID string
	Jobs         []Job
	CloudConfig  CloudConfig
}

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
}

type CloudConfig struct {
	AvailabilityZones           []string
	PersistentDiskPools         []DiskConfig
	PersistentDiskTypes         []DiskConfig
	Networks                    []NetworkConfig
	CompilationNetwork          string
	CompilationAvailabilityZone string
	VmTypes                     []VmTypeConfig
	ResourcePools               []VmTypeConfig
}

type DiskConfig struct {
	Name string
	Size int
}

type VmTypeConfig struct {
	Name string
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

type SubnetConfig struct {
	AvailabilityZones []string
	IpPool            *IpPool
}

type JobNetworkConfig struct {
	Name          string
	DefaultDNSnGW bool
	StaticIps     []string
}
