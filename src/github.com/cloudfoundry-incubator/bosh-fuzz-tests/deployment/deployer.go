package deployment

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshsys "github.com/cloudfoundry/bosh-utils/system"

	bftanalyzer "github.com/cloudfoundry-incubator/bosh-fuzz-tests/analyzer"
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
	inputGenerator       InputGenerator
	networksAssigner     NetworksAssigner
	analyzer             bftanalyzer.Analyzer
	fs                   boshsys.FileSystem
	logger               boshlog.Logger
	generateManifestOnly bool
}

func NewDeployer(
	cliRunner bltclirunner.Runner,
	directorInfo bltaction.DirectorInfo,
	renderer Renderer,
	inputGenerator InputGenerator,
	networksAssigner NetworksAssigner,
	analyzer bftanalyzer.Analyzer,
	fs boshsys.FileSystem,
	logger boshlog.Logger,
	generateManifestOnly bool,
) Deployer {
	return &deployer{
		cliRunner:            cliRunner,
		directorInfo:         directorInfo,
		renderer:             renderer,
		inputGenerator:       inputGenerator,
		networksAssigner:     networksAssigner,
		analyzer:             analyzer,
		fs:                   fs,
		logger:               logger,
		generateManifestOnly: generateManifestOnly,
	}
}

func (d *deployer) RunDeploys() error {
	manifestPath, err := d.fs.TempFile("manifest")
	if err != nil {
		return bosherr.WrapError(err, "Creating manifest file")
	}
	defer d.fs.RemoveAll(manifestPath.Name())

	cloudConfigPath, err := d.fs.TempFile("cloud-config")
	if err != nil {
		return bosherr.WrapError(err, "Creating cloud config file")
	}
	defer d.fs.RemoveAll(cloudConfigPath.Name())

	inputs, err := d.inputGenerator.Generate()
	if err != nil {
		return bosherr.WrapError(err, "Generating input")
	}

	d.networksAssigner.Assign(inputs)

	cases := d.analyzer.Analyze(inputs)

	for _, testCase := range cases {
		input := testCase.Input
		input.DirectorUUID = d.directorInfo.UUID

		err = d.renderer.Render(input, manifestPath.Name(), cloudConfigPath.Name())
		if err != nil {
			return bosherr.WrapError(err, "Rendering deployment manifest")
		}

		if !d.generateManifestOnly {
			err = d.cliRunner.RunWithArgs("update", "cloud-config", cloudConfigPath.Name())
			if err != nil {
				return bosherr.WrapError(err, "Updating cloud config")
			}

			err = d.cliRunner.RunWithArgs("deployment", manifestPath.Name())
			if err != nil {
				return bosherr.WrapError(err, "Setting deployment manifest")
			}

			deployWrapper := bltaction.NewDeployWrapper(d.cliRunner)
			taskId, err := deployWrapper.RunWithDebug("deploy")
			if err != nil {
				return bosherr.WrapError(err, "Running deploy")
			}

			debugLog, err := d.cliRunner.RunWithOutput("task", taskId, "--debug")
			if err != nil {
				return bosherr.WrapError(err, "Getting task debug logs")
			}

			for _, expectation := range testCase.Expectations {
				err := expectation.Run(debugLog)
				if err != nil {
					return bosherr.WrapError(err, "Running expectation")
				}
			}
		}
	}

	return nil
}
