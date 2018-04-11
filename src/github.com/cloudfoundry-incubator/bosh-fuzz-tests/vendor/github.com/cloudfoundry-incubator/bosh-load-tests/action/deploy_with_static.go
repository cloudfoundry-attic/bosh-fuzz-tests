package action

import (
	"bytes"
	"fmt"
	"net"
	"text/template"

	bltclirunner "github.com/cloudfoundry-incubator/bosh-load-tests/action/clirunner"
	bltassets "github.com/cloudfoundry-incubator/bosh-load-tests/assets"

	"strings"

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

var numInstancesPerFlow = 5

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

	var staticIPs []string
	for i := 0; i < numInstancesPerFlow; i++ {
		staticIPs = append(staticIPs, d.GetNextIP(i))
	}

	data := manifestData{
		DeploymentName:  d.deploymentName,
		DirectorUUID:    d.directorInfo.UUID,
		DeploymentIndex: d.flowNumber,
		StaticIPs:       strings.Join(staticIPs, ","),
		NumInstances:    numInstancesPerFlow,
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

func (d *deployWithStatic) GetNextIP(i int) string {
	ip := net.ParseIP(fmt.Sprintf("10.245.%d.0", d.flowNumber))
	b := ip.To4()
	reservedRange := 11 // reserve 10.245.0.0 to 10.245.0.10. bosh director lives at 10.245.0.3
	instanceIndex := reservedRange + numInstancesPerFlow + i
	b[3] = b[3] + byte(1+(instanceIndex%253))
	return net.IPv4(b[0], b[1], b[2], b[3]).String()
}
