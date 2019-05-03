package config

import (
	"encoding/json"

	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

type Config struct {
	fs boshsys.FileSystem

	GenerateManifestOnly      bool       `json:"generate_manifest_only"`
	Parameters                Parameters `json:"parameters"`
	NumberOfConsequentDeploys int        `json:"number_of_consequent_deploys"`
}

type Parameters struct {
	NameLength                 []int      `json:"name_length"`
	Instances                  []int      `json:"instances"`
	AvailabilityZones          [][]string `json:"availability_zones"`
	PersistentDiskDefinition   []string   `json:"persistent_disk_definition"`
	PersistentDiskSize         []int      `json:"persistent_disk_size"`
	NumberOfInstanceGroups     []int      `json:"number_of_instance_groups"`
	MigratedFromCount          []int      `json:"migrated_from_count"`
	Networks                   [][]string `json:"networks"`
	VmTypeDefinition           []string   `json:"vm_type_definition"`
	StemcellDefinition         []string   `json:"stemcell_definition"`
	StemcellVersions           []string   `json:"stemcell_versions"`
	Jobs                       [][]string `json:"jobs"`
	NumberOfCompilationWorkers []int      `json:"number_of_compilation_workers"`
	Canaries                   []int      `json:"canaries"`
	MaxInFlight                []int      `json:"max_in_flight"`
	Serial                     []string   `json:"serial"`
	NumOfCloudProperties       []int      `json:"num_of_cloud_properties"`
	NumOfVariables             []int      `json:"num_of_variables"`
	VariableTypes              []string   `json:"variable_types"`
	NumOfSubstitutions         []int      `json:"num_of_substitutions"`
	CpiAPIVersion              []int      `json:"preferred_cpi_api_version"`
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
