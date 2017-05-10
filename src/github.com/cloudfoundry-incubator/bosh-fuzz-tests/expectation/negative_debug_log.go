package expectation

import (
	"strings"

	bltclirunner "github.com/cloudfoundry-incubator/bosh-load-tests/action/clirunner"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

type negativeDebugLog struct {
	unexpectedString string
}

func NewNegativeDebugLog(unexpectedString string) Expectation {
	return &negativeDebugLog{
		unexpectedString: unexpectedString,
	}
}

func (d *negativeDebugLog) Run(cliRunner bltclirunner.Runner, taskId string) error {
	debugLog, err := cliRunner.RunWithOutput("task", taskId, "--debug")
	if err != nil {
		return bosherr.WrapError(err, "Getting task debug logs")
	}

	if strings.Contains(debugLog, d.unexpectedString) {
		return bosherr.Errorf("Task debug logs output contains unexpected string: %s", d.unexpectedString)
	}

	return nil
}
