package expectation

import (
	"strings"

	bltclirunner "github.com/cloudfoundry-incubator/bosh-load-tests/action/clirunner"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

type debugLog struct {
	expectedString string
	cliRunner      bltclirunner.Runner
}

func NewDebugLog(expectedString string, cliRunner bltclirunner.Runner) Expectation {
	return &debugLog{
		expectedString: expectedString,
		cliRunner:      cliRunner,
	}
}

func (d *debugLog) Run(taskId string) error {
	output, err := d.cliRunner.RunWithOutput("task", taskId, "--debug")
	if err != nil {
		return bosherr.WrapError(err, "Getting task debug logs")
	}

	if !strings.Contains(output, d.expectedString) {
		return bosherr.Errorf("Task debug logs output does not contain expected string: %s", d.expectedString)
	}

	return nil
}
