package action

import (
	bltclirunner "github.com/cloudfoundry-incubator/bosh-load-tests/action/clirunner"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"regexp"
)

var taskIDRegex = regexp.MustCompilePOSIX(`^Task ([0-9]+)$`)

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
	taskID := ""

	matches := taskIDRegex.FindStringSubmatch(output)
	if len(matches) > 0 {
		taskID = matches[1]
	} else {
		return "", bosherr.Error("Failed to get task id")
	}

	if err != nil {
		debugErr := d.cliRunner.RunWithArgs("task", taskID, "--debug")
		if debugErr != nil {
			return taskID, debugErr
		}
	}

	return taskID, err
}
