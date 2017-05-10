package expectation

import (
	"strings"

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

func (d *debugLog) Run(debugLog string) error {
	if !strings.Contains(debugLog, d.expectedString) {
		return bosherr.Errorf("Task debug logs output does not contain expected string: %s", d.expectedString)
	}

	return nil
}
