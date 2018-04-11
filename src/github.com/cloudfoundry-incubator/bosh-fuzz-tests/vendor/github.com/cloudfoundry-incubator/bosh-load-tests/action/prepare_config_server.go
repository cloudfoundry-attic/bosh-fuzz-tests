package action

import (
	"encoding/json"
	"net/url"

	"errors"
	"fmt"
	"os"
	"strings"

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
	urlWithoutPort := strings.Split(targetURL.Host, ":")[0]
	targetURL.Host = fmt.Sprintf("%s:8443", urlWithoutPort)
	targetURL.Scheme = "https"

	target := targetURL.String()
	if err := p.uaaRunner.RunWithArgs("target", target, "--skip-ssl-validation"); nil != err {
		return err
	}

	if err := p.uaaRunner.RunWithArgs("token", "client", "get", "director_config_server", "-s", os.Getenv("CONFIG_SERVER_PASSWORD")); nil != err {
		return err
	}

	if err := p.setValue("/num_instances", 5); nil != err {
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

	if directorIP, exist := os.LookupEnv("BOSH_DIRECTOR_IP"); exist {
		if err := p.uaaRunner.RunWithArgs("curl", "--insecure", "--request", "PUT", "--header", "Content-Type:Application/JSON", "--data", string(data), fmt.Sprintf("https://%s:8080/v1/data", directorIP)); nil != err {
			return err
		}
	} else {
		return errors.New("could not find environment: BOSH_DIRECTOR_IP")
	}

	return nil
}
