package analyzer

import (
	bftexpectation "github.com/cloudfoundry-incubator/bosh-fuzz-tests/expectation"
	bftinput "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
)

type variablesComparator struct {
	logger boshlog.Logger
}

func NewVariablesComparator() Comparator {
	return &variablesComparator{}
}

func (c *variablesComparator) Compare(previousInputs []bftinput.Input, currentInput bftinput.Input) []bftexpectation.Expectation {
	return []bftexpectation.Expectation{bftexpectation.NewVariablesExpectation(currentInput.Variables)}
}
