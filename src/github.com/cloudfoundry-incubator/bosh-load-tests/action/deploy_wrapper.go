package action

import (
	"regexp"

	bltclirunner "github.com/cloudfoundry-incubator/bosh-load-tests/action/clirunner"
)

type deployWrapper struct {
	cliRunner bltclirunner.Runner
}

func NewDeployWrapper(cliRunner bltclirunner.Runner) *deployWrapper {
	return &deployWrapper{
		cliRunner: cliRunner,
	}
}

func (d *deployWrapper) RunWithDebug(args ...string) error {
	output, err := d.cliRunner.RunWithOutput(args...)
	if err != nil {
		re := regexp.MustCompile("bosh task ([0-9]+) --debug")
		matches := re.FindAllStringSubmatch(output, -1)
		if len(matches) > 0 && len(matches[0][1]) > 1 {
			taskId := matches[0][1]
			err := d.cliRunner.RunWithArgs("task", taskId, "--debug")
			if err != nil {
				return err
			}
		}
	}

	return nil
}
