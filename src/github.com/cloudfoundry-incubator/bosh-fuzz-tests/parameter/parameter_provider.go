package parameter

import (
	"math/rand"

	bftconfig "github.com/cloudfoundry-incubator/bosh-fuzz-tests/config"
)

type ParameterProvider interface {
	Get(name string) Parameter
}

type parameterProvider struct {
	parameters bftconfig.Parameters
}

func NewParameterProvider(parameters bftconfig.Parameters) ParameterProvider {
	return &parameterProvider{
		parameters: parameters,
	}
}

func (p *parameterProvider) Get(name string) Parameter {
	if name == "stemcell" {
		stemcellDefinition := p.parameters.StemcellDefinition[rand.Intn(len(p.parameters.StemcellDefinition))]
		return NewStemcell(stemcellDefinition)
	}

	return nil
}
