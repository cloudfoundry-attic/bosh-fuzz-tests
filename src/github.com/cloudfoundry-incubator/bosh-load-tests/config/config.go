package config

import (
	"encoding/json"

	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

type Config struct {
	AssetsPath               string             `json:"assets_path"`
	Environment              string             `json:"environment"`
	DirectorMigrationCommand string             `json:"director_migration_cmd"`
	DirectorStartCommand     string             `json:"director_start_cmd"`
	WorkerStartCommand       string             `json:"worker_start_cmd"`
	NginxStartCommand        string             `json:"nginx_start_cmd"`
	VerifyMultidigest        string             `json:"verify_multidigest"`
	NatsStartCommand         string             `json:"nats_start_cmd"`
	UAAConfig                UAAConfig          `json:"uaa"`
	ConfigServerConfig       ConfigServerConfig `json:"config_server"`
	DummyCPIPath             string             `json:"dummy_cpi_path"`
	RubyVersion              string             `json:"ruby_version"`
	CliCmd                   string             `json:"cli_cmd"`
	Flows                    [][]string         `json:"flows"`
	NumberOfWorkers          int                `json:"number_of_workers"`
	NumberOfDeployments      int                `json:"number_of_deployments"`
	UsingLegacyManifest      bool               `json:"using_legacy_manifest"`

	fs boshsys.FileSystem
}

type UAAConfig struct {
	Enabled    bool   `json:"enabled"`
	TomcatPath string `json:"tomcat_path"`
}

type ConfigServerConfig struct {
	Enabled                  bool   `json:"enabled"`
	ConfigServerStartCommand string `json:"config_server_start_cmd"`
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
