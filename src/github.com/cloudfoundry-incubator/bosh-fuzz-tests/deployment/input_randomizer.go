package deployment

import (
	"math/rand"
	"time"

	bftconfig "github.com/cloudfoundry-incubator/bosh-fuzz-tests/config"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
)

type InputRandomizer interface {
	Generate() ([]Input, error)
}

type inputRandomizer struct {
	parameters                bftconfig.Parameters
	numberOfConsequentDeploys int
	seed                      int64
	logger                    boshlog.Logger
}

func NewSeededInputRandomizer(parameters bftconfig.Parameters, numberOfConsequentDeploys int, seed int64, logger boshlog.Logger) InputRandomizer {
	return &inputRandomizer{
		parameters:                parameters,
		numberOfConsequentDeploys: numberOfConsequentDeploys,
		seed:   seed,
		logger: logger,
	}
}

func NewInputRandomizer(parameters bftconfig.Parameters, numberOfConsequentDeploys int, logger boshlog.Logger) InputRandomizer {
	return &inputRandomizer{
		parameters:                parameters,
		numberOfConsequentDeploys: numberOfConsequentDeploys,
		seed:   time.Now().Unix(),
		logger: logger,
	}
}

func (ir *inputRandomizer) Generate() ([]Input, error) {
	ir.logger.Info("inputRandomizer", "Seeding with %d", ir.seed)

	rand.Seed(ir.seed)

	inputs := []Input{}

	nameGenerator := NewNameGenerator()

	for i := 0; i < ir.numberOfConsequentDeploys; i++ {
		inputs = append(inputs, Input{
			Name:              nameGenerator.Generate(ir.parameters.NameLength[rand.Intn(len(ir.parameters.NameLength))]),
			Instances:         ir.parameters.Instances[rand.Intn(len(ir.parameters.Instances))],
			AvailabilityZones: ir.parameters.AvailabilityZones[rand.Intn(len(ir.parameters.AvailabilityZones))],
		})
	}
	return inputs, nil
}
