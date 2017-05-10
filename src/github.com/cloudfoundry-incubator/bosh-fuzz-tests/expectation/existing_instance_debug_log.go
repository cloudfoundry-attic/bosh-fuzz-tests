package expectation

import (
	"fmt"
	"regexp"
	"strings"

	bltclirunner "github.com/cloudfoundry-incubator/bosh-load-tests/action/clirunner"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

type existingInstanceDebugLog struct {
	expectedString string
	cliRunner      bltclirunner.Runner
	jobName        string
}

func NewExistingInstanceDebugLog(expectedString string, jobName string) Expectation {
	return &existingInstanceDebugLog{
		expectedString: expectedString,
		jobName:        jobName,
	}
}

func (d *existingInstanceDebugLog) Run(cliRunner bltclirunner.Runner, taskId string) error {

	debugLog, err := cliRunner.RunWithOutput("task", taskId, "--debug")
	if err != nil {
		return bosherr.WrapError(err, "Getting task debug logs")
	}

	regexString := fmt.Sprintf("Existing desired instance '(%s[^']+)'", d.jobName)
	re := regexp.MustCompile(regexString)
	matches := re.FindAllStringSubmatch(debugLog, -1)

	for _, match := range matches {
		if len(match) > 1 {
			instanceName := match[1]
			instanceNameParts := strings.Split(instanceName, "/")
			expectedRe := regexp.MustCompile(fmt.Sprintf("%s.* %s\\/[0-9a-f]{8}-[0-9a-f-]{27} \\(%s\\)", d.expectedString, instanceNameParts[0], instanceNameParts[1]))
			expectedMatches := expectedRe.FindAllStringSubmatch(debugLog, -1)

			if len(expectedMatches) == 0 {
				return bosherr.Errorf("Task debug logs output does not contain expected string: %s for instance %s", d.expectedString, instanceName)
			}
		}
	}

	return nil
}
