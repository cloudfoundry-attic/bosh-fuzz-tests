package action

import (
	"errors"

	bltclirunner "github.com/cloudfoundry-incubator/bosh-load-tests/action/clirunner"
	bltassets "github.com/cloudfoundry-incubator/bosh-load-tests/assets"

	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

type Factory interface {
	Create(name string, flowNumber int, deploymentName string, boshRunner bltclirunner.Runner, uaacRunner bltclirunner.Runner, usingLegacyManifest bool) (Action, error)
}

type factory struct {
	directorInfo        DirectorInfo
	fs                  boshsys.FileSystem
	assetsProvider      bltassets.Provider
	usingLegacyManifest bool
}

func NewFactory(
	directorInfo DirectorInfo,
	fs boshsys.FileSystem,
	assetsProvider bltassets.Provider,
) *factory {
	return &factory{
		directorInfo:   directorInfo,
		fs:             fs,
		assetsProvider: assetsProvider,
	}
}

func (f *factory) Create(
	name string,
	flowNumber int,
	deploymentName string,
	boshRunner bltclirunner.Runner,
	uaacRunner bltclirunner.Runner,
	usingLegacyManifest bool,
) (Action, error) {

	switch name {
	case "prepare":
		return NewPrepare(f.directorInfo, boshRunner, f.fs, f.assetsProvider), nil
	case "prepare_config_server":
		return NewPrepareConfigServer(f.directorInfo, uaacRunner), nil
	case "ignore":
		return NewIgnore(f.directorInfo, deploymentName, boshRunner, f.fs, f.assetsProvider), nil
	case "upload_cloud_config":
		return NewUploadCloudConfig(f.directorInfo, boshRunner, f.assetsProvider), nil
	case "deploy_with_dynamic":
		return NewDeployWithDynamic(f.directorInfo, deploymentName, boshRunner, f.fs, f.assetsProvider, usingLegacyManifest), nil
	case "deploy_with_static":
		return NewDeployWithStatic(f.directorInfo, flowNumber, deploymentName, boshRunner, f.fs, f.assetsProvider, usingLegacyManifest), nil
	case "deploy_with_variables":
		return NewDeployWithVariables(f.directorInfo, deploymentName, boshRunner, f.fs, f.assetsProvider), nil
	case "recreate":
		return NewRecreate(f.directorInfo, deploymentName, boshRunner, f.fs), nil
	case "stop_hard":
		return NewStopHard(f.directorInfo, deploymentName, boshRunner, f.fs), nil
	case "start":
		return NewStart(f.directorInfo, deploymentName, boshRunner, f.fs), nil
	case "delete_deployment":
		return NewDeleteDeployment(f.directorInfo, deploymentName, boshRunner, f.fs), nil
	}

	return nil, errors.New("unknown action")
}
