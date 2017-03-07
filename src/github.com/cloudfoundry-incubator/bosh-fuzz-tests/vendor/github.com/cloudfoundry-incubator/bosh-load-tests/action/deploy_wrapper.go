package action

import (
	"fmt"

	"encoding/json"
	bltclirunner "github.com/cloudfoundry-incubator/bosh-load-tests/action/clirunner"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

type deployWrapper struct {
	cliRunner bltclirunner.Runner
}

func NewDeployWrapper(cliRunner bltclirunner.Runner) *deployWrapper {
	return &deployWrapper{
		cliRunner: cliRunner,
	}
}

func (d *deployWrapper) RunWithDebug(args ...string) (string, error) {
	output, err := d.cliRunner.RunWithOutput(args...)
	taskId := ""

	output, err = d.cliRunner.RunWithOutput("tasks", "--recent=1", "--json")
	if err != nil {
		return taskId, err
	}

	var outputStruct Output
	json.Unmarshal([]byte(output), &outputStruct)

	if outputStruct.Tables != nil {
		for _, row := range outputStruct.Tables[0].Rows {
			if val, found := row["0"]; found {
				taskId = val.(string)
				break
			}
		}
	} else {
		fmt.Println(fmt.Sprintf("OUTPUT: %s", output))
	}

	if err != nil {
		debugErr := d.cliRunner.RunWithArgs("task", taskId, "--debug")
		if debugErr != nil {
			return taskId, debugErr
		}
	}

	if taskId == "" {
		return "", bosherr.Error("Failed to get task id")
	}

	return taskId, err
}
