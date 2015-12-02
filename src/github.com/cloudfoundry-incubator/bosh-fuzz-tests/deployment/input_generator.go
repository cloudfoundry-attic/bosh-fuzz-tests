package deployment

import (
	"math/rand"

	bftconfig "github.com/cloudfoundry-incubator/bosh-fuzz-tests/config"
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
	logger                    boshlog.Logger
}

func NewInputGenerator(
	parameters bftconfig.Parameters,
	parameterProvider bftparam.ParameterProvider,
	numberOfConsequentDeploys int,
	nameGenerator bftnamegen.NameGenerator,
	logger boshlog.Logger,
) InputGenerator {
	return &inputGenerator{
		parameters:                parameters,
		parameterProvider:         parameterProvider,
		numberOfConsequentDeploys: numberOfConsequentDeploys,
		nameGenerator:             nameGenerator,
		logger:                    logger,
	}
}

func (g *inputGenerator) Generate() ([]bftinput.Input, error) {
	inputs := []bftinput.Input{}

	for i := 0; i < g.numberOfConsequentDeploys; i++ {
		jobNames := g.generateJobNames(i, inputs)
		input := g.generateInput(jobNames, false)

		migratingJobs := []string{}
		for _, j := range input.Jobs {
			for _, m := range j.MigratedFrom {
				migratingJobs = append(migratingJobs, m.Name)
			}
		}

		if len(migratingJobs) > 0 {
			migratingInput := g.generateInput(migratingJobs, true)

			g.specifyAzIfMigratingJobDoesNotHaveAz(migratingInput, input)

			inputs = append(inputs, migratingInput)
		}

		inputs = append(inputs, input)
	}

	return inputs, nil
}

func (g *inputGenerator) generateInput(jobNames []string, migratingDeployment bool) bftinput.Input {
	input := &bftinput.Input{
		Jobs: []bftinput.Job{},
	}

	for _, jobName := range jobNames {
		job := &bftinput.Job{
			Name:      jobName,
			Instances: g.parameters.Instances[rand.Intn(len(g.parameters.Instances))],
		}

		if !migratingDeployment {
			migratedFromCount := g.parameters.MigratedFromCount[rand.Intn(len(g.parameters.MigratedFromCount))]
			for i := 0; i < migratedFromCount; i++ {
				migratedFromName := g.nameGenerator.Generate(10)
				job.MigratedFrom = append(job.MigratedFrom, bftinput.MigratedFromConfig{Name: migratedFromName})
			}
		}

		input.Jobs = append(input.Jobs, *job)
	}

	input = g.parameterProvider.Get("availability_zone").Apply(input)
	input = g.parameterProvider.Get("vm_type").Apply(input)
	input = g.parameterProvider.Get("stemcell").Apply(input)
	input = g.parameterProvider.Get("persistent_disk").Apply(input)

	return *input
}

func (g *inputGenerator) generateJobNames(i int, inputs []bftinput.Input) []string {
	numberOfJobs := g.parameters.NumberOfJobs[rand.Intn(len(g.parameters.NumberOfJobs))]
	jobNames := []string{}

	for j := 0; j < numberOfJobs; j++ {
		var jobName string
		if i > 0 && len(inputs[i-1].Jobs) > j {
			jobName = inputs[i-1].Jobs[j].Name
		}

		if jobName == "" {
			jobName = g.nameGenerator.Generate(g.parameters.NameLength[rand.Intn(len(g.parameters.NameLength))])
		}

		jobNames = append(jobNames, jobName)
	}

	return jobNames
}

func (g *inputGenerator) specifyAzIfMigratingJobDoesNotHaveAz(migratingInput bftinput.Input, currentInput bftinput.Input) {
	for _, migratingJob := range migratingInput.Jobs {
		if migratingJob.AvailabilityZones == nil {
			for k, job := range currentInput.Jobs {
				if job.AvailabilityZones == nil {
					continue
				}

				for m, migratedFromConfig := range job.MigratedFrom {
					if migratedFromConfig.Name == migratingJob.Name {
						currentInput.Jobs[k].MigratedFrom[m].AvailabilityZone = currentInput.Jobs[k].AvailabilityZones[0]
					}
				}
			}
		}
	}
}
