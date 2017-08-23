package action

import (
	"encoding/json"
	"errors"
	"strings"

	bltclirunner "github.com/cloudfoundry-incubator/bosh-load-tests/action/clirunner"
)

type InstancesInfo struct {
	directorInfo   DirectorInfo
	deploymentName string
	cliRunner      bltclirunner.Runner
}

type Instance struct {
	Name string
	ID   string
}

type instancesOutput struct {
	Tables []struct {
		Rows []map[string]string
	}
}

func NewInstances(directorInfo DirectorInfo, deploymentName string, cliRunner bltclirunner.Runner) *InstancesInfo {
	return &InstancesInfo{
		directorInfo:   directorInfo,
		deploymentName: deploymentName,
		cliRunner:      cliRunner,
	}
}

func (i *InstancesInfo) GetInstances() (map[string][]Instance, error) {
	instances := map[string][]Instance{}

	output, err := i.cliRunner.RunWithOutput("-d", i.deploymentName, "instances", "--json")
	if err != nil {
		return instances, err
	}

	outputStruct := instancesOutput{}
	err = json.Unmarshal([]byte(output), &outputStruct)
	if err != nil {
		return instances, errors.New("error unmarshalling JSON")
	}

	if len(outputStruct.Tables) > 0 {
		for _, row := range outputStruct.Tables[0].Rows {
			instanceSlug := row["instance"]
			instanceParts := strings.Split(instanceSlug, "/")
			if len(instanceParts) > 1 {
				instanceName := instanceParts[0]

				instances[instanceName] = append(instances[instanceName], Instance{Name: instanceName, ID: instanceParts[1]})
			}
		}
	}

	return instances, err
}
