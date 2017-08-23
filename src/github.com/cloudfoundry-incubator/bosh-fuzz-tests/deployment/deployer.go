package deployment

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshsys "github.com/cloudfoundry/bosh-utils/system"

	bftanalyzer "github.com/cloudfoundry-incubator/bosh-fuzz-tests/analyzer"
	"github.com/cloudfoundry-incubator/bosh-fuzz-tests/variables"

	bltaction "github.com/cloudfoundry-incubator/bosh-load-tests/action"
	bltclirunner "github.com/cloudfoundry-incubator/bosh-load-tests/action/clirunner"

	"encoding/json"
	"net/url"
)

type Deployer interface {
	RunDeploys() error
	CasesRun() []bftanalyzer.Case
}

type deployer struct {
	cliRunner            bltclirunner.Runner
	uaaRunner            bltclirunner.Runner
	directorInfo         bltaction.DirectorInfo
	renderer             Renderer
	inputGenerator       InputGenerator
	stepGenerators       []StepGenerator
	analyzer             bftanalyzer.Analyzer
	sprinkler            variables.Sprinkler
	fs                   boshsys.FileSystem
	logger               boshlog.Logger
	generateManifestOnly bool
	casesRun						 []bftanalyzer.Case
}

func NewDeployer(
	cliRunner bltclirunner.Runner,
	uaaRunner bltclirunner.Runner,
	directorInfo bltaction.DirectorInfo,
	renderer Renderer,
	inputGenerator InputGenerator,
	stepGenerators []StepGenerator,
	analyzer bftanalyzer.Analyzer,
	sprinkler variables.Sprinkler,
	fs boshsys.FileSystem,
	logger boshlog.Logger,
	generateManifestOnly bool,
) Deployer {
	return &deployer{
		cliRunner:            cliRunner,
		uaaRunner:            uaaRunner,
		directorInfo:         directorInfo,
		renderer:             renderer,
		inputGenerator:       inputGenerator,
		stepGenerators:       stepGenerators,
		analyzer:             analyzer,
		sprinkler:            sprinkler,
		fs:                   fs,
		logger:               logger,
		generateManifestOnly: generateManifestOnly,
	}
}

func (d *deployer) RunDeploys() error {
	d.cliRunner.SetEnv(d.directorInfo.URL)

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

	cases := d.analyzer.Analyze(inputs)
	d.casesRun = []bftanalyzer.Case{}

	for _, testCase := range cases {
		input := testCase.Input
		input.DirectorUUID = d.directorInfo.UUID

		err = d.renderer.Render(input, manifestPath.Name(), cloudConfigPath.Name())
		if err != nil {
			return bosherr.WrapError(err, "Rendering deployment manifest")
		}

		substitutionMap, err := d.sprinkler.SprinklePlaceholders(manifestPath.Name())
		if err != nil {
			return bosherr.WrapError(err, "Could not sprinkle placholders in manifest")
		}

		// Setup UAA
		targetURL, err := url.Parse(d.directorInfo.URL)
		if nil != err {
			return err
		}
		targetURL.Scheme = "https"
		targetURL.Path = "/uaa"

		target := targetURL.String()
		if err := d.uaaRunner.RunWithArgs("target", target, "--skip-ssl-validation"); nil != err {
			return err
		}

		for key, value := range substitutionMap {

			stringMapValue := d.convertToStringMap(value)

			dataStruct := struct {
				Name  string      `json:"name"`
				Value interface{} `json:"value"`
			}{"/TestDirector/foo-deployment/" + key, stringMapValue}

			data, err := json.Marshal(dataStruct)
			if nil != err {
				return err
			}
			if err := d.uaaRunner.RunWithArgs("token", "client", "get", "test", "-s", "secret"); nil != err {
				return err
			}
			if err := d.uaaRunner.RunWithArgs("curl", "--insecure", "--request", "PUT", "--header", "Content-Type:Application/JSON", "--data", string(data), "https://localhost:65005/v1/data"); nil != err {
				return err
			}
		}

		if err != nil {
			return bosherr.WrapError(err, "Populating config server key values")
		}

		if !d.generateManifestOnly {
			err = d.cliRunner.RunWithArgs("update-cloud-config", cloudConfigPath.Name())
			if err != nil {
				return bosherr.WrapError(err, "Updating cloud config")
			}

			deployWrapper := bltaction.NewDeployWrapper(d.cliRunner)
			taskId, err := deployWrapper.RunWithDebug("-d", "foo-deployment", "deploy", manifestPath.Name())
			if err != nil {
				if testCase.DeploymentWillFail {
					continue
				}
				return bosherr.WrapError(err, "Running deploy")
			}

			instancesCaller := bltaction.NewInstances(d.directorInfo, "foo-deployment", d.cliRunner)
			instances, err := instancesCaller.GetInstances()
			if err != nil {
				return bosherr.WrapError(err, "Listing instances")
			}
			testCase.InstancesAfterDeploy = instances

			d.casesRun = append(d.casesRun, testCase)

			for _, expectation := range testCase.Expectations {
				err := expectation.Run(d.cliRunner, taskId)
				if err != nil {
					return bosherr.WrapError(err, "Running expectation")
				}
			}

			steps := []Step{}
			for _, stepGenerator := range d.stepGenerators {
				steps = append(steps, stepGenerator.Steps(testCase)...)
			}

			for _, step := range steps {
				err = step.Run(d.cliRunner)
				if err != nil {
					return bosherr.WrapError(err, "Running step")
				}
			}
		}
	}

	return nil
}

func (d deployer) CasesRun() []bftanalyzer.Case {
	return d.casesRun
}

func (d deployer) convertToStringMap(obj interface{}) interface{} {

	switch obj.(type) {
	case []interface{}:
		outputArray := []interface{}{}

		for _, item := range obj.([]interface{}) {
			outputArray = append(outputArray, d.convertToStringMap(item))
		}
		obj = outputArray
	case map[interface{}]interface{}:
		outputMap := map[string]interface{}{}

		for key, value := range obj.(map[interface{}]interface{}) {
			outputMap[key.(string)] = d.convertToStringMap(value)
		}
		obj = outputMap
	default:
		return obj
	}
	return obj
}
