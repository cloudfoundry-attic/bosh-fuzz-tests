package parameter

import (
	bftinput "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"
)

type stemcell struct {
	definition string
}

func NewStemcell(definition) Parameter {
	return &stemcell{
		definition: definition,
	}
}

func (s *stemcell) Apply(input bftinput.Input) bftinput.Input {

}
