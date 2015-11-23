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
}

type CloudConfig struct {
	AvailabilityZones           []string
	PersistentDiskPools         []DiskConfig
	PersistentDiskTypes         []DiskConfig
	Networks                    []NetworkConfig
	CompilationNetwork          string
	CompilationAvailabilityZone string
}

type DiskConfig struct {
	Name string
	Size int
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
	IpRange           string
	Gateway           string
	AvailabilityZones []string
	Reserved          []string
}

type JobNetworkConfig struct {
	Name          string
	DefaultDNSnGW bool
}
