package action

import (
	"bytes"
	"net"
	"text/template"

	bltclirunner "github.com/cloudfoundry-incubator/bosh-load-tests/action/clirunner"
	bltassets "github.com/cloudfoundry-incubator/bosh-load-tests/assets"

	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

type deployWithStatic struct {
	directorInfo        DirectorInfo
	flowNumber          int
	deploymentName      string
	cliRunner           bltclirunner.Runner
	fs                  boshsys.FileSystem
	assetsProvider      bltassets.Provider
	usingLegacyManifest bool
}

func NewDeployWithStatic(
	directorInfo DirectorInfo,
	flowNumber int,
	deploymentName string,
	cliRunner bltclirunner.Runner,
	fs boshsys.FileSystem,
	assetsProvider bltassets.Provider,
	usingLegacyManifest bool,
) *deployWithStatic {
	return &deployWithStatic{
		directorInfo:        directorInfo,
		flowNumber:          flowNumber,
		deploymentName:      deploymentName,
		cliRunner:           cliRunner,
		fs:                  fs,
		assetsProvider:      assetsProvider,
		usingLegacyManifest: usingLegacyManifest,
	}
}

func (d *deployWithStatic) Execute() error {
	d.cliRunner.SetEnv(d.directorInfo.URL)

	manifestFilename := "manifest_with_static.yml"
	if d.usingLegacyManifest == true {
		manifestFilename = "legacy_manifest_with_static.yml"
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
		StaticIP:       d.getNextIP(),
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

func (d *deployWithStatic) getNextIP() string {
	ip := net.ParseIP("192.168.1.10")
	b := ip.To4()
	b[3] = b[3] + byte(d.flowNumber)
	return net.IPv4(b[0], b[1], b[2], b[3]).String()
}
