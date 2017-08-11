package action

import (
	"encoding/json"
	"net/url"

	bltclirunner "github.com/cloudfoundry-incubator/bosh-load-tests/action/clirunner"
)

type prepareConfigServer struct {
	directorInfo DirectorInfo
	uaaRunner    bltclirunner.Runner
}

func NewPrepareConfigServer(
	directorInfo DirectorInfo,
	uaaRunner bltclirunner.Runner,
) *prepareConfigServer {
	return &prepareConfigServer{
		directorInfo: directorInfo,
		uaaRunner:    uaaRunner,
	}
}

func (p *prepareConfigServer) Execute() error {
	// Setup UAA
	targetURL, err := url.Parse(p.directorInfo.URL)
	if nil != err {
		return err
	}
	targetURL.Scheme = "https"
	targetURL.Path = "/uaa"

	target := targetURL.String()
	if err := p.uaaRunner.RunWithArgs("target", target, "--skip-ssl-validation"); nil != err {
		return err
	}

	if err := p.uaaRunner.RunWithArgs("token", "client", "get", "test", "-s", "secret"); nil != err {
		return err
	}

	if err := p.setValue("/num_instances", 2); nil != err {
		return err
	}

	if err := p.setValue("/prop3_value", "this is the value of prop 3!"); nil != err {
		return err
	}

	return nil
}

func (p *prepareConfigServer) setValue(key string, value interface{}) error {
	dataStruct := struct {
		Name  string      `json:"name"`
		Value interface{} `json:"value"`
	}{key, value}

	data, err := json.Marshal(dataStruct)
	if nil != err {
		return err
	}

	if err := p.uaaRunner.RunWithArgs("curl", "--insecure", "--request", "PUT", "--header", "Content-Type:Application/JSON", "--data", string(data), "https://localhost:65005/v1/data"); nil != err {
		return err
	}

	return nil
}
