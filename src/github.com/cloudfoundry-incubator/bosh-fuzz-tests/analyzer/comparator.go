package analyzer

import (
	bftexpectation "github.com/cloudfoundry-incubator/bosh-fuzz-tests/expectation"
	bftinput "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"
)

type Comparator interface {
	Compare(previousInput bftinput.Input, currentInput bftinput.Input) []bftexpectation.Expectation
}
