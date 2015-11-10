package manifest

import (
	"bytes"
	"text/template"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

type Input struct {
	Name              string
	DirectorUUID      string
	Instances         int
	AvailabilityZones []string
}

type Renderer interface {
	Render(input Input, manifestPath string) error
}

type renderer struct {
	fs boshsys.FileSystem
}

func NewRenderer(fs boshsys.FileSystem) Renderer {
	return &renderer{
		fs: fs,
	}
}

func (g *renderer) Render(input Input, manifestPath string) error {
	deploymentTemplate := template.Must(template.New("deployment").Parse(DeploymentTemplate))

	buffer := bytes.NewBuffer([]byte{})

	err := deploymentTemplate.Execute(buffer, input)
	if err != nil {
		return bosherr.WrapErrorf(err, "Generating deployment manifest")
	}

	err = g.fs.WriteFile(manifestPath, buffer.Bytes())
	if err != nil {
		return bosherr.WrapErrorf(err, "Saving generated manifest")
	}

	return nil
}
