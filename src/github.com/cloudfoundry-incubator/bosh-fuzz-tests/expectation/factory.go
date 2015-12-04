package expectation

import (
	bltclirunner "github.com/cloudfoundry-incubator/bosh-load-tests/action/clirunner"
)

type Factory interface {
	CreateDebugLog(expectedString string) Expectation
}

type factory struct {
	cliRunner bltclirunner.Runner
}

func NewFactory(cliRunner bltclirunner.Runner) Factory {
	return &factory{
		cliRunner: cliRunner,
	}
}

func (f *factory) CreateDebugLog(expectedString string) Expectation {
	return NewDebugLog(expectedString, f.cliRunner)
}
