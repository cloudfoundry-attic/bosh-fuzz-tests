package action

import (
	"errors"

	bltclirunner "github.com/cloudfoundry-incubator/bosh-load-tests/action/clirunner"
	bltassets "github.com/cloudfoundry-incubator/bosh-load-tests/assets"

	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

type Factory interface {
	Create(name string, flowNumber int, deploymentName string, cliRunner bltclirunner.Runner, usingLegacyManifest bool) (Action, error)
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
	cliRunner bltclirunner.Runner,
	usingLegacyManifest bool,
) (Action, error) {
	switch name {
	case "prepare":
		return NewPrepare(f.directorInfo, cliRunner, f.fs, f.assetsProvider), nil
	case "ignore":
		return NewIgnore(f.directorInfo, deploymentName, cliRunner, f.fs, f.assetsProvider), nil
	case "upload_cloud_config":
		return NewUploadCloudConfig(f.directorInfo, cliRunner, f.assetsProvider), nil
	case "deploy_with_dynamic":
		return NewDeployWithDynamic(f.directorInfo, deploymentName, cliRunner, f.fs, f.assetsProvider, usingLegacyManifest), nil
	case "deploy_with_static":
		return NewDeployWithStatic(f.directorInfo, flowNumber, deploymentName, cliRunner, f.fs, f.assetsProvider, usingLegacyManifest), nil
	case "recreate":
		return NewRecreate(f.directorInfo, deploymentName, cliRunner, f.fs), nil
	case "stop_hard":
		return NewStopHard(f.directorInfo, deploymentName, cliRunner, f.fs), nil
	case "start":
		return NewStart(f.directorInfo, deploymentName, cliRunner, f.fs), nil
	case "delete_deployment":
		return NewDeleteDeployment(f.directorInfo, deploymentName, cliRunner, f.fs), nil
	}

	return nil, errors.New("unknown action")
}
