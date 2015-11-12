package deployment

import (
	"math/rand"

	bftconfig "github.com/cloudfoundry-incubator/bosh-fuzz-tests/config"
)

type InputRandomizer interface {
	Generate() ([]Input, error)
}

type inputRandomizer struct {
	parameters                bftconfig.Parameters
	numberOfConsequentDeploys int
	seed                      int64
}

func NewSeededInputRandomizer(parameters bftconfig.Parameters, numberOfConsequentDeploys int, seed int64) InputRandomizer {
	return &inputRandomizer{
		parameters:                parameters,
		numberOfConsequentDeploys: numberOfConsequentDeploys,
		seed: seed,
	}
}

func NewInputRandomizer(parameters bftconfig.Parameters, numberOfConsequentDeploys int) InputRandomizer {
	return &inputRandomizer{
		parameters:                parameters,
		numberOfConsequentDeploys: numberOfConsequentDeploys,
		seed: 1,
	}
}

func (ir *inputRandomizer) Generate() ([]Input, error) {
	rand.Seed(ir.seed)

	inputs := []Input{}

	for i := 0; i < ir.numberOfConsequentDeploys; i++ {
		inputs = append(inputs, Input{
			Name:              ir.parameters.Name[rand.Intn(len(ir.parameters.Name))],
			Instances:         ir.parameters.Instances[rand.Intn(len(ir.parameters.Instances))],
			AvailabilityZones: ir.parameters.AvailabilityZones[rand.Intn(len(ir.parameters.AvailabilityZones))],
		})
	}
	return inputs, nil
}
