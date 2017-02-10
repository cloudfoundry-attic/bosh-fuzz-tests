package expectation

import (
	"strings"

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

func (d *negativeDebugLog) Run(debugLog string) error {
	if strings.Contains(debugLog, d.unexpectedString) {
		return bosherr.Errorf("Task debug logs output contains unexpected string: %s", d.unexpectedString)
	}

	return nil
}
