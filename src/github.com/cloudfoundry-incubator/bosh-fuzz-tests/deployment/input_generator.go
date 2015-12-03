package deployment

import (
	"fmt"
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

	jobNames := g.generateJobNames()
	previousInput := g.generateInputWithJobNames(jobNames)

	for i := 0; i < g.numberOfConsequentDeploys; i++ {
		input := g.fuzzInput(previousInput, false)

		migratingJobNames := []string{}
		for _, j := range input.Jobs {
			for _, m := range j.MigratedFrom {
				migratingJobNames = append(migratingJobNames, m.Name)
			}
		}

		if len(migratingJobNames) > 0 {
			migratingInput := g.generateInputWithJobNames(migratingJobNames)
			migratingInput = g.fuzzInput(migratingInput, true)

			g.specifyAzIfMigratingJobDoesNotHaveAz(migratingInput, input)

			inputs = append(inputs, migratingInput)
		}

		inputs = append(inputs, input)

		previousInput = input
	}

	g.logger.Debug("input_generator", fmt.Sprintf("Generated inputs: %#v", inputs))

	return inputs, nil
}

func (g *inputGenerator) fuzzInput(previousInput bftinput.Input, migratingDeployment bool) bftinput.Input {
	input := bftinput.Input{
		CloudConfig: previousInput.CloudConfig,
		Stemcells:   previousInput.Stemcells,
	}
	for _, job := range previousInput.Jobs {
		input.Jobs = append(input.Jobs, job)
	}
	// input.Jobs = g.randomizeJobs(input.Jobs)

	for j, job := range input.Jobs {
		input.Jobs[j].Instances = g.parameters.Instances[rand.Intn(len(g.parameters.Instances))]
		input.Jobs[j].MigratedFrom = nil

		if !migratingDeployment {
			migratedFromCount := g.parameters.MigratedFromCount[rand.Intn(len(g.parameters.MigratedFromCount))]
			for i := 0; i < migratedFromCount; i++ {
				migratedFromName := g.nameGenerator.Generate(10)
				input.Jobs[j].MigratedFrom = append(job.MigratedFrom, bftinput.MigratedFromConfig{Name: migratedFromName})
			}
		}
	}

	input = g.parameterProvider.Get("availability_zone").Apply(input)
	input = g.parameterProvider.Get("vm_type").Apply(input)
	input = g.parameterProvider.Get("stemcell").Apply(input)
	input = g.parameterProvider.Get("persistent_disk").Apply(input)

	return input
}

// func (g *inputGenerator) randomizeJobs(jobs []bftinput.Job) []bftinput.Job {
// 	numberOfJobs := g.parameters.NumberOfJobs[rand.Intn(len(g.parameters.NumberOfJobs))]
// 	if numberOfJobs > len(jobs) {
// 		for i := 0; i < numberOfJobs-len(jobs); i++ {
// 			jobName := g.nameGenerator.Generate(g.parameters.NameLength[rand.Intn(len(g.parameters.NameLength))])
// 			jobs = append(jobs, bftinput.Job{
// 				Name: jobName,
// 			})
// 		}
// 	} else if numberOfJobs < len(jobs) {
// 		for i := 0; i < len(jobs)-numberOfJobs; i++ {
// 			jobIdxToRemove := rand.Intn(len(jobs))
// 			jobs = append(jobs[:jobIdxToRemove], jobs[jobIdxToRemove+1:]...)
// 		}
// 	}

// 	return jobs
// }

func (g *inputGenerator) generateInputWithJobNames(jobNames []string) bftinput.Input {
	input := bftinput.Input{
		Jobs: []bftinput.Job{},
	}
	for _, jobName := range jobNames {
		input.Jobs = append(input.Jobs, bftinput.Job{
			Name: jobName,
		})
	}

	return input
}

func (g *inputGenerator) generateJobNames() []string {
	numberOfJobs := g.parameters.NumberOfJobs[rand.Intn(len(g.parameters.NumberOfJobs))]
	jobNames := []string{}

	for j := 0; j < numberOfJobs; j++ {
		jobName := g.nameGenerator.Generate(g.parameters.NameLength[rand.Intn(len(g.parameters.NameLength))])
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
