package expectation

import (
	bltclirunner "github.com/cloudfoundry-incubator/bosh-load-tests/action/clirunner"
)

type Expectation interface {
	Run(cliRunner bltclirunner.Runner, taskId string) error
}
