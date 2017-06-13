package deployment

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshsys "github.com/cloudfoundry/bosh-utils/system"

	bftanalyzer "github.com/cloudfoundry-incubator/bosh-fuzz-tests/analyzer"
	"github.com/cloudfoundry-incubator/bosh-fuzz-tests/variables"

	bltaction "github.com/cloudfoundry-incubator/bosh-load-tests/action"
	bltclirunner "github.com/cloudfoundry-incubator/bosh-load-tests/action/clirunner"
	bltconfig "github.com/cloudfoundry-incubator/bosh-load-tests/config"

	"encoding/json"
	"fmt"
	"net/url"
)

type Deployer interface {
	RunDeploys() error
}

type deployer struct {
	cliRunner            bltclirunner.Runner
	directorInfo         bltaction.DirectorInfo
	renderer             Renderer
	inputGenerator       InputGenerator
	analyzer             bftanalyzer.Analyzer
	sprinkler            variables.Sprinkler
	absoluteSprinkler    variables.Sprinkler
	fs                   boshsys.FileSystem
	logger               boshlog.Logger
	generateManifestOnly bool
}

func NewDeployer(
	cliRunner bltclirunner.Runner,
	directorInfo bltaction.DirectorInfo,
	renderer Renderer,
	inputGenerator InputGenerator,
	analyzer bftanalyzer.Analyzer,
	sprinkler variables.Sprinkler,
	absoluteSprinkler variables.Sprinkler,
	fs boshsys.FileSystem,
	logger boshlog.Logger,
	generateManifestOnly bool,
) Deployer {
	return &deployer{
		cliRunner:            cliRunner,
		directorInfo:         directorInfo,
		renderer:             renderer,
		inputGenerator:       inputGenerator,
		analyzer:             analyzer,
		sprinkler:            sprinkler,
		absoluteSprinkler:    absoluteSprinkler,
		fs:                   fs,
		logger:               logger,
		generateManifestOnly: generateManifestOnly,
	}
}


var badManifestPatterns [][]interface{} = [][]interface{}{
	{"instance_groups", variables.Integer, "env"},
	{"instance_groups", variables.Integer, "jobs", variables.Integer, "consumes", variables.String, "properties"},
	{"instance_groups", variables.Integer, "jobs", variables.Integer, "properties"},
	{"instance_groups", variables.Integer, "properties"},
	{"jobs", variables.Integer, "env"},
	{"jobs", variables.Integer, "properties"},
	{"jobs", variables.Integer, "templates", variables.Integer, "consumes", variables.String, "properties"},
	{"jobs", variables.Integer, "templates", variables.Integer, "properties"},
	{"name"},
	{"properties"},
	{"variables", variables.Integer, variables.String},
	{"variables", variables.Integer},
	{"variables"},
	{"releases", variables.Integer, variables.String}, // should be supported. not working now.
	{"releases", variables.Integer},                   // should be supported. not working now.
	{"releases"},                                      // should be supported. not working now.
	{"stemcells", variables.Integer},                  // should be supported. not working now.
	{"stemcells"},                                     // should be supported. not working now.
}

var badCloudConfigPatterns [][]interface{} = [][]interface{}{
	{"properties"},
	{"resource_pools", variables.Integer},
	{"resource_pools", variables.Integer, "env"},
	{"resource_pools", variables.Integer, "stemcell"},
	{"resource_pools", variables.Integer, "stemcell", variables.String},
	{"resource_pools"},
}


func (d *deployer) RunDeploys() error {
	d.cliRunner.SetEnv(d.directorInfo.URL)

	inputs, err := d.inputGenerator.Generate()
	if err != nil {
		return bosherr.WrapError(err, "Generating input")
	}

	logger := boshlog.NewLogger(boshlog.LevelDebug)
	cmdRunner := boshsys.NewExecCmdRunner(logger)
	fs := boshsys.NewOsFileSystem(logger)

	envConfig := bltconfig.NewConfig(fs)

	cliRunnerFactory := bltclirunner.NewFactory(envConfig.CliCmd, cmdRunner, fs)

	uaaRunner, err := cliRunnerFactory.Create("uaac")
	if err != nil {
		panic(err)
	}

	cases := d.analyzer.Analyze(inputs)

	for _, testCase := range cases {
		input := testCase.Input
		input.DirectorUUID = d.directorInfo.UUID

		manifestPath, err := d.fs.TempFile("manifest")
		if err != nil {
			return bosherr.WrapError(err, "Creating manifest file")
		}
		//defer d.fs.RemoveAll(manifestPath.Name())

		cloudConfigPath, err := d.fs.TempFile("cloud-config")
		if err != nil {
			return bosherr.WrapError(err, "Creating cloud config file")
		}
		//defer d.fs.RemoveAll(cloudConfigPath.Name())

		err = d.renderer.Render(input, manifestPath.Name(), cloudConfigPath.Name())
		if err != nil {
			return bosherr.WrapError(err, "Rendering deployment manifest")
		}

		substitutionMap, err := d.sprinkler.SprinklePlaceholders(manifestPath.Name(), badManifestPatterns)
		if err != nil {
			return bosherr.WrapError(err, "Could not sprinkle placeholders in manifest")
		}

		cloudConfigSubstitutionMap, err := d.absoluteSprinkler.SprinklePlaceholders(cloudConfigPath.Name(), badCloudConfigPatterns)
		if err != nil {
			return bosherr.WrapError(err, "Could not sprinkle placeholders in cloud config")
		}

		// Setup UAA
		targetURL, err := url.Parse(d.directorInfo.URL)
		if nil != err {
			return err
		}
		targetURL.Scheme = "https"
		targetURL.Path = "/uaa"

		target := targetURL.String()
		if err := uaaRunner.RunWithArgs("target", target, "--skip-ssl-validation"); nil != err {
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
			if err := uaaRunner.RunWithArgs("token", "client", "get", "test", "-s", "secret"); nil != err {
				return err
			}
			if err := uaaRunner.RunWithArgs("curl", "--insecure", "--request", "PUT", "--header", "Content-Type:Application/JSON", "--data", string(data), "https://localhost:65005/v1/data"); nil != err {
				return err
			}
		}

		for key, value := range cloudConfigSubstitutionMap {

			stringMapValue := d.convertToStringMap(value)

			dataStruct := struct {
				Name  string      `json:"name"`
				Value interface{} `json:"value"`
			}{key, stringMapValue}
			fmt.Println("\ndataStructName", dataStruct.Name)
			fmt.Println("\ndataStructValue", dataStruct.Value)

			data, err := json.Marshal(dataStruct)
			if nil != err {
				return err
			}
			if err := uaaRunner.RunWithArgs("token", "client", "get", "test", "-s", "secret"); nil != err {
				return err
			}
			if err := uaaRunner.RunWithArgs("curl", "--insecure", "--request", "PUT", "--header", "Content-Type:Application/JSON", "--data", string(data), "https://localhost:65005/v1/data"); nil != err {
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
				errorPrefix := ""
				if testCase.DeploymentWillFail {
					errorPrefix += "\n==========================================================\n"
					errorPrefix += "DEPLOYMENT FAILURE IS EXPECTED DUE TO UNSUPPORTED SCENARIO\n"
					errorPrefix += "==========================================================\n"
					continue
				}
				return bosherr.WrapError(err, errorPrefix+"Running deploy")
			}

			for _, expectation := range testCase.Expectations {
				err := expectation.Run(d.cliRunner, taskId)
				if err != nil {
					return bosherr.WrapError(err, "Running expectation")
				}
			}
		}
	}

	return nil
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
