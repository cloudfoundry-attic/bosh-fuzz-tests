package deployment

import (
	boshsys "github.com/cloudfoundry/bosh-utils/system"

	bltaction "github.com/cloudfoundry-incubator/bosh-load-tests/action"
	bltclirunner "github.com/cloudfoundry-incubator/bosh-load-tests/action/clirunner"
)

type Deployer interface {
	RunDeploys() error
}

type deployer struct {
	cliRunner            bltclirunner.Runner
	directorInfo         bltaction.DirectorInfo
	renderer             Renderer
	jobsRandomizer       JobsRandomizer
	networksAssigner     NetworksAssigner
	fs                   boshsys.FileSystem
	generateManifestOnly bool
}

func NewDeployer(
	cliRunner bltclirunner.Runner,
	directorInfo bltaction.DirectorInfo,
	renderer Renderer,
	jobsRandomizer JobsRandomizer,
	networksAssigner NetworksAssigner,
	fs boshsys.FileSystem,
	generateManifestOnly bool,
) Deployer {
	return &deployer{
		cliRunner:            cliRunner,
		directorInfo:         directorInfo,
		renderer:             renderer,
		jobsRandomizer:       jobsRandomizer,
		networksAssigner:     networksAssigner,
		fs:                   fs,
		generateManifestOnly: generateManifestOnly,
	}
}

func (d *deployer) RunDeploys() error {
	manifestPath, err := d.fs.TempFile("manifest")
	if err != nil {
		return err
	}
	defer d.fs.RemoveAll(manifestPath.Name())

	cloudConfigPath, err := d.fs.TempFile("cloud-config")
	if err != nil {
		return err
	}
	defer d.fs.RemoveAll(cloudConfigPath.Name())

	inputs, err := d.jobsRandomizer.Generate()
	if err != nil {
		return err
	}

	d.networksAssigner.Assign(inputs)

	for _, input := range inputs {
		input.DirectorUUID = d.directorInfo.UUID

		err = d.renderer.Render(input, manifestPath.Name(), cloudConfigPath.Name())
		if err != nil {
			return err
		}

		if !d.generateManifestOnly {
			err = d.cliRunner.RunWithArgs("update", "cloud-config", cloudConfigPath.Name())
			if err != nil {
				return err
			}

			err = d.cliRunner.RunWithArgs("deployment", manifestPath.Name())
			if err != nil {
				return err
			}

			deployWrapper := bltaction.NewDeployWrapper(d.cliRunner)
			err = deployWrapper.RunWithDebug("deploy")
			if err != nil {
				return err
			}
		}
	}

	return nil
}
