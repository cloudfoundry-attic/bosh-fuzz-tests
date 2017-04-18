package action

import (
	"fmt"
	bltclirunner "github.com/cloudfoundry-incubator/bosh-load-tests/action/clirunner"
	bltassets "github.com/cloudfoundry-incubator/bosh-load-tests/assets"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

type ignore struct {
	directorInfo   DirectorInfo
	deploymentName string
	cliRunner      bltclirunner.Runner
	fs             boshsys.FileSystem
	assetsProvider bltassets.Provider
}

func NewIgnore(
	directorInfo DirectorInfo,
	deploymentName string,
	cliRunner bltclirunner.Runner,
	fs boshsys.FileSystem,
	assetsProvider bltassets.Provider,
) *deployWithDynamic {
	return &deployWithDynamic{
		directorInfo:   directorInfo,
		deploymentName: deploymentName,
		cliRunner:      cliRunner,
		fs:             fs,
		assetsProvider: assetsProvider,
	}
}

func (d *ignore) Execute() error {
	output, err := d.cliRunner.RunWithOutput("-d", d.deploymentName, "instances")
	if err != nil {
		panic(err)
	}
	fmt.Println(output)

	d.cliRunner.SetEnv(d.directorInfo.URL)

	deployWrapper := NewDeployWrapper(d.cliRunner)
	_, err = deployWrapper.RunWithDebug("-d", d.deploymentName, "ignore", "simple/0")
	if err != nil {
		return err
	}

	return nil
}
