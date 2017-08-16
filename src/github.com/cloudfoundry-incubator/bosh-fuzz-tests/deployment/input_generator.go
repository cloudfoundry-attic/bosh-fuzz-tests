package deployment

import (
	"math/rand"

	bftconfig "github.com/cloudfoundry-incubator/bosh-fuzz-tests/config"
	bftdecider "github.com/cloudfoundry-incubator/bosh-fuzz-tests/decider"
	bftinput "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"
	bftnamegen "github.com/cloudfoundry-incubator/bosh-fuzz-tests/name_generator"
	bftparam "github.com/cloudfoundry-incubator/bosh-fuzz-tests/parameter"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
)

type InputGenerator interface {
	Generate() ([]bftinput.Input, error)
}

type inputGenerator struct {
	parameters                bftconfig.Parameters
	parameterProvider         bftparam.ParameterProvider
	numberOfConsequentDeploys int
	nameGenerator             bftnamegen.NameGenerator
	decider                   bftdecider.Decider
	logger                    boshlog.Logger
}

func NewInputGenerator(
	parameters bftconfig.Parameters,
	parameterProvider bftparam.ParameterProvider,
	numberOfConsequentDeploys int,
	nameGenerator bftnamegen.NameGenerator,
	decider bftdecider.Decider,
	logger boshlog.Logger,
) InputGenerator {
	return &inputGenerator{
		parameters:                parameters,
		parameterProvider:         parameterProvider,
		numberOfConsequentDeploys: numberOfConsequentDeploys,
		nameGenerator:             nameGenerator,
		decider:                   decider,
		logger:                    logger,
	}
}

func (g *inputGenerator) Generate() ([]bftinput.Input, error) {
	inputs := []bftinput.Input{}

	instanceGroupNames := g.generateInstanceGroupNames()
	previousInput := g.generateInputWithInstanceGroupNames(instanceGroupNames)

	for i := 0; i < g.numberOfConsequentDeploys; i++ {
		reusePreviousInput := g.decider.IsYes()
		var input bftinput.Input

		if i > 0 && reusePreviousInput {
			input = previousInput
		} else {
			input = g.createInputFromPrevious(previousInput)
			migratingInstanceGroupNames := []string{}
			for _, j := range input.InstanceGroups {
				for _, m := range j.MigratedFrom {
					migratingInstanceGroupNames = append(migratingInstanceGroupNames, m.Name)
				}
			}

			if len(migratingInstanceGroupNames) == 0 {
				input = g.fuzzInput(input, previousInput)
			} else {
				migratingInput := g.generateInputWithInstanceGroupNames(migratingInstanceGroupNames)
				migratingInput = g.fuzzInput(migratingInput, previousInput)
				input = g.fuzzInput(input, migratingInput)

				g.specifyAzIfMigratingInstanceGroupDoesNotHaveAz(migratingInput, input)

				inputs = append(inputs, migratingInput)
			}
		}

		inputs = append(inputs, input)

		previousInput = input
	}

	return inputs, nil
}

func (g *inputGenerator) createInputFromPrevious(previousInput bftinput.Input) bftinput.Input {
	input := bftinput.Input{}

	for _, instanceGroup := range previousInput.InstanceGroups {
		instanceGroup.Instances = g.parameters.Instances[rand.Intn(len(g.parameters.Instances))]
		instanceGroup.MigratedFrom = nil

		input.InstanceGroups = append(input.InstanceGroups, instanceGroup)
	}

	input.InstanceGroups = g.randomizeInstanceGroups(input.InstanceGroups)

	for j := range input.InstanceGroups {
		migratedFromCount := g.parameters.MigratedFromCount[rand.Intn(len(g.parameters.MigratedFromCount))]
		for i := 0; i < migratedFromCount; i++ {
			migratedFromName := g.nameGenerator.Generate(10)
			input.InstanceGroups[j].MigratedFrom = append(input.InstanceGroups[j].MigratedFrom, bftinput.MigratedFromConfig{Name: migratedFromName})
		}
	}

	return input
}

func (g *inputGenerator) fuzzInput(input bftinput.Input, previousInput bftinput.Input) bftinput.Input {
	input.CloudConfig = previousInput.CloudConfig
	input.Stemcells = previousInput.Stemcells

	input = g.parameterProvider.Get("variables").Apply(input, previousInput)
	input = g.parameterProvider.Get("availability_zone").Apply(input, previousInput)
	input = g.parameterProvider.Get("vm_type").Apply(input, previousInput)
	input = g.parameterProvider.Get("stemcell").Apply(input, previousInput)
	input = g.parameterProvider.Get("persistent_disk").Apply(input, previousInput)
	input = g.parameterProvider.Get("network").Apply(input, previousInput)
	input = g.parameterProvider.Get("job").Apply(input, previousInput)
	input = g.parameterProvider.Get("compilation").Apply(input, previousInput)
	input = g.parameterProvider.Get("update").Apply(input, previousInput)
	input = g.parameterProvider.Get("cloud_properties").Apply(input, previousInput)
	input = g.parameterProvider.Get("fixed_migrated_from").Apply(input, previousInput)
	input = g.parameterProvider.Get("lifecycle").Apply(input, previousInput)

	return input
}

func (g *inputGenerator) randomizeInstanceGroups(instanceGroups []bftinput.InstanceGroup) []bftinput.InstanceGroup {
	numberOfInstanceGroups := g.parameters.NumberOfInstanceGroups[rand.Intn(len(g.parameters.NumberOfInstanceGroups))]
	instanceGroupsSize := len(instanceGroups)
	if numberOfInstanceGroups > instanceGroupsSize {
		for i := 0; i < numberOfInstanceGroups-instanceGroupsSize; i++ {
			instanceGroupName := g.nameGenerator.Generate(g.parameters.NameLength[rand.Intn(len(g.parameters.NameLength))])
			instanceGroups = append(instanceGroups, bftinput.InstanceGroup{
				Name: instanceGroupName,
			})
		}
	} else if numberOfInstanceGroups < instanceGroupsSize {
		for i := 0; i < instanceGroupsSize-numberOfInstanceGroups; i++ {
			instanceGroupIdxToRemove := rand.Intn(len(instanceGroups))
			if instanceGroupIdxToRemove == len(instanceGroups)-1 {
				instanceGroups = instanceGroups[:instanceGroupIdxToRemove]
			} else {
				instanceGroups = append(instanceGroups[:instanceGroupIdxToRemove], instanceGroups[instanceGroupIdxToRemove+1:]...)
			}
		}
	}

	shuffledInstanceGroups := []bftinput.InstanceGroup{}
	shuffledInstanceGroupsIndeces := rand.Perm(numberOfInstanceGroups)
	for _, instanceGroupIndex := range shuffledInstanceGroupsIndeces {
		shuffledInstanceGroups = append(shuffledInstanceGroups, instanceGroups[instanceGroupIndex])
	}

	return shuffledInstanceGroups
}

func (g *inputGenerator) generateInputWithInstanceGroupNames(instanceGroupNames []string) bftinput.Input {
	input := bftinput.Input{
		InstanceGroups: []bftinput.InstanceGroup{},
	}

	for _, instanceGroupName := range instanceGroupNames {
		input.InstanceGroups = append(input.InstanceGroups, bftinput.InstanceGroup{
			Name:      instanceGroupName,
			Instances: g.parameters.Instances[rand.Intn(len(g.parameters.Instances))],
		})
	}

	return input
}

func (g *inputGenerator) generateInstanceGroupNames() []string {
	numberOfInstanceGroups := g.parameters.NumberOfInstanceGroups[rand.Intn(len(g.parameters.NumberOfInstanceGroups))]
	instanceGroupNames := []string{}

	for j := 0; j < numberOfInstanceGroups; j++ {
		instanceGroupName := g.nameGenerator.Generate(g.parameters.NameLength[rand.Intn(len(g.parameters.NameLength))])
		instanceGroupNames = append(instanceGroupNames, instanceGroupName)
	}

	return instanceGroupNames
}

func (g *inputGenerator) specifyAzIfMigratingInstanceGroupDoesNotHaveAz(migratingInput bftinput.Input, currentInput bftinput.Input) {
	for _, migratingInstanceGroup := range migratingInput.InstanceGroups {
		if migratingInstanceGroup.AvailabilityZones == nil {
			for k, instanceGroup := range currentInput.InstanceGroups {
				if instanceGroup.AvailabilityZones == nil {
					continue
				}

				for m, migratedFromConfig := range instanceGroup.MigratedFrom {
					if migratedFromConfig.Name == migratingInstanceGroup.Name {
						currentInput.InstanceGroups[k].MigratedFrom[m].AvailabilityZone = currentInput.InstanceGroups[k].AvailabilityZones[0]
					}
				}
			}
		}
	}
}
