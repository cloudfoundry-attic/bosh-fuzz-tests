package manifest

import (
	bftconfig "github.com/cloudfoundry-incubator/bosh-fuzz-tests/config"
)

type InputRandomizer interface {
	Generate(parameters bftconfig.Parameters, numberOfConsequentDeploys int) ([]Input, error)
}

type inputRandomizer struct {
}

func NewInputRandomizer() InputRandomizer {
	return &inputRandomizer{}
}

func (i *inputRandomizer) Generate(parameters bftconfig.Parameters, numberOfConsequentDeploys int) ([]Input, error) {
	return []Input{}, nil
}
