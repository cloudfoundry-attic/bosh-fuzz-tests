package action

import (
	bltclirunner "github.com/cloudfoundry-incubator/bosh-load-tests/action/clirunner"
	bltassets "github.com/cloudfoundry-incubator/bosh-load-tests/assets"
)

type uploadCloudConfig struct {
	directorInfo   DirectorInfo
	cliRunner      bltclirunner.Runner
	assetsProvider bltassets.Provider
}

func NewUploadCloudConfig(
	directorInfo DirectorInfo,
	cliRunner bltclirunner.Runner,
	assetsProvider bltassets.Provider,
) *uploadCloudConfig {
	return &uploadCloudConfig{
		directorInfo:   directorInfo,
		cliRunner:      cliRunner,
		assetsProvider: assetsProvider,
	}
}

func (p *uploadCloudConfig) Execute() error {
	err := p.cliRunner.TargetAndLogin(p.directorInfo.URL)
	if err != nil {
		return err
	}

	cloudConfigPath, err := p.assetsProvider.FullPath("cloud_config.yml")
	if err != nil {
		return err
	}

	err = p.cliRunner.RunWithArgs("update", "cloud-config", cloudConfigPath)
	if err != nil {
		return err
	}

	return nil
}
