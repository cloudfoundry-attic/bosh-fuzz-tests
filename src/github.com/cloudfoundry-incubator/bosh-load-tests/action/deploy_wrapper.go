package action

import (
	"regexp"

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

	if err != nil {
		re := regexp.MustCompile("bosh task ([0-9]+) --debug")
		matches := re.FindAllStringSubmatch(output, -1)

		if len(matches) > 0 && len(matches[0]) > 1 {
			taskId = matches[0][1]
			debugErr := d.cliRunner.RunWithArgs("task", taskId, "--debug")
			if debugErr != nil {
				return taskId, debugErr
			}
		}
	} else {
		re := regexp.MustCompile("Task ([0-9]+) done")
		matches := re.FindAllStringSubmatch(output, -1)

		if len(matches) > 0 && len(matches[0]) > 1 {
			taskId = matches[0][1]
		}
	}

	if taskId == "" {
		return "", bosherr.Error("Failed to get task id")
	}

	return taskId, err
}
