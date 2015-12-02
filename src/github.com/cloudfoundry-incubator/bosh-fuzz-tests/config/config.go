package config

import (
	"encoding/json"

	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

type Config struct {
	fs boshsys.FileSystem

	Parameters                Parameters `json:"parameters"`
	NumberOfConsequentDeploys int        `json:"number_of_consequent_deploys"`
}

type Parameters struct {
	NameLength               []int      `json:"name_length"`
	Instances                []int      `json:"instances"`
	AvailabilityZones        [][]string `json:"availability_zones"`
	PersistentDiskDefinition []string   `json:"persistent_disk_definition"`
	PersistentDiskSize       []int      `json:"persistent_disk_size"`
	NumberOfJobs             []int      `json:"number_of_jobs"`
	MigratedFromCount        []int      `json:"migrated_from_count"`
	Networks                 [][]string `json:"networks"`
	VmTypeDefinition         []string   `json:"vm_type_definition"`
	StemcellDefinition       []string   `json:"stemcell_definition"`
}

func NewConfig(fs boshsys.FileSystem) *Config {
	return &Config{
		fs: fs,
	}
}

func (c *Config) Load(configPath string) error {
	contents, err := c.fs.ReadFile(configPath)
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(contents), &c)
	if err != nil {
		return err
	}

	return nil
}
