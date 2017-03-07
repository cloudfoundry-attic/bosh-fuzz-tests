package action

import (
	"bytes"
	"text/template"

	bltclirunner "github.com/cloudfoundry-incubator/bosh-load-tests/action/clirunner"
	bltassets "github.com/cloudfoundry-incubator/bosh-load-tests/assets"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

type deployWithDynamic struct {
	directorInfo        DirectorInfo
	deploymentName      string
	cliRunner           bltclirunner.Runner
	fs                  boshsys.FileSystem
	assetsProvider      bltassets.Provider
	usingLegacyManifest bool
}

func NewDeployWithDynamic(
	directorInfo DirectorInfo,
	deploymentName string,
	cliRunner bltclirunner.Runner,
	fs boshsys.FileSystem,
	assetsProvider bltassets.Provider,
	usingLegacyManifest bool,
) *deployWithDynamic {
	return &deployWithDynamic{
		directorInfo:        directorInfo,
		deploymentName:      deploymentName,
		cliRunner:           cliRunner,
		fs:                  fs,
		assetsProvider:      assetsProvider,
		usingLegacyManifest: usingLegacyManifest,
	}
}

func (d *deployWithDynamic) Execute() error {
	d.cliRunner.SetEnv(d.directorInfo.URL)

	manifestFilename := "manifest.yml"
	if d.usingLegacyManifest == true {
		manifestFilename = "legacy_manifest.yml"
	}

	manifestTemplatePath, err := d.assetsProvider.FullPath(manifestFilename)
	if err != nil {
		return err
	}

	manifestPath, err := d.fs.TempFile("manifest-test")
	if err != nil {
		return err
	}
	defer d.fs.RemoveAll(manifestPath.Name())

	t := template.Must(template.ParseFiles(manifestTemplatePath))
	buffer := bytes.NewBuffer([]byte{})
	data := manifestData{
		DeploymentName: d.deploymentName,
		DirectorUUID:   d.directorInfo.UUID,
	}
	err = t.Execute(buffer, data)
	if err != nil {
		return err
	}
	err = d.fs.WriteFile(manifestPath.Name(), buffer.Bytes())
	if err != nil {
		return err
	}

	deployWrapper := NewDeployWrapper(d.cliRunner)
	_, err = deployWrapper.RunWithDebug("-d", d.deploymentName, "deploy", manifestPath.Name())
	if err != nil {
		return err
	}

	return nil
}
