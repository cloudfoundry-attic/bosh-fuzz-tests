package expectation

import (
	"strings"

	bltclirunner "github.com/cloudfoundry-incubator/bosh-load-tests/action/clirunner"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

type debugLog struct {
	expectedString string
}

func NewDebugLog(expectedString string) Expectation {
	return &debugLog{
		expectedString: expectedString,
	}
}

func (d *debugLog) Run(cliRunner bltclirunner.Runner, taskId string) error {
	debugLog, err := cliRunner.RunWithOutput("task", taskId, "--debug")
	if err != nil {
		return bosherr.WrapError(err, "Getting task debug logs")
	}

	if !strings.Contains(debugLog, d.expectedString) {
		return bosherr.Errorf("Task debug logs output does not contain expected string: %s", d.expectedString)
	}

	return nil
}
