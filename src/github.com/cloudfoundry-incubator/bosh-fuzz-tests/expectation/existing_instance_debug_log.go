package expectation

import (
	"fmt"
	"regexp"

	bltclirunner "github.com/cloudfoundry-incubator/bosh-load-tests/action/clirunner"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

type existingInstanceDebugLog struct {
	expectedString string
	cliRunner      bltclirunner.Runner
}

func NewExistingInstanceDebugLog(expectedString string) Expectation {
	return &existingInstanceDebugLog{
		expectedString: expectedString,
	}
}

func (d *existingInstanceDebugLog) Run(debugLog string) error {
	re := regexp.MustCompile("Existing desired instance '(.*)'")
	matches := re.FindAllStringSubmatch(debugLog, -1)

	for _, match := range matches {
		if len(match) > 1 {
			instanceName := match[1]
			expectedRe := regexp.MustCompile(fmt.Sprintf("%s.* %s", d.expectedString, instanceName))
			expectedMatches := expectedRe.FindAllStringSubmatch(debugLog, -1)

			if len(expectedMatches) == 0 {
				return bosherr.Errorf("Task debug logs output does not contain expected string: %s", d.expectedString)
			}
		}
	}

	return nil
}
