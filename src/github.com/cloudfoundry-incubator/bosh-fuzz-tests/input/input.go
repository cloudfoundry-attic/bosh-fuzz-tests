package input

type Input struct {
	DirectorUUID string
	Jobs         []Job
	CloudConfig  CloudConfig
	Stemcells    []StemcellConfig
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
	Stemcell           string
}

type CloudConfig struct {
	AvailabilityZones           []string
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

type VmTypeConfig struct {
	Name string
}

type ResourcePoolConfig struct {
	Name     string
	Stemcell StemcellConfig
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

type SubnetConfig struct {
	AvailabilityZones []string
	IpPool            *IpPool
}

type JobNetworkConfig struct {
	Name          string
	DefaultDNSnGW bool
	StaticIps     []string
}
