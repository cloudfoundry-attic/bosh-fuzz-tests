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

	jobNames := g.generateJobNames()
	previousInput := g.generateInputWithJobNames(jobNames)

	for i := 0; i < g.numberOfConsequentDeploys; i++ {
		reusePreviousInput := g.decider.IsYes()
		var input bftinput.Input

		if i > 0 && reusePreviousInput {
			input = previousInput
		} else {
			input = g.createInputFromPrevious(previousInput)
			migratingJobNames := []string{}
			for _, j := range input.Jobs {
				for _, m := range j.MigratedFrom {
					migratingJobNames = append(migratingJobNames, m.Name)
				}
			}

			if len(migratingJobNames) == 0 {
				input = g.fuzzInput(input, previousInput)
			} else {
				migratingInput := g.generateInputWithJobNames(migratingJobNames)
				migratingInput = g.fuzzInput(migratingInput, previousInput)
				input = g.fuzzInput(input, migratingInput)

				g.specifyAzIfMigratingJobDoesNotHaveAz(migratingInput, input)

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

	for _, job := range previousInput.Jobs {
		job.Instances = g.parameters.Instances[rand.Intn(len(g.parameters.Instances))]
		job.MigratedFrom = nil

		input.Jobs = append(input.Jobs, job)
	}

	input.Jobs = g.randomizeJobs(input.Jobs)

	for j := range input.Jobs {
		migratedFromCount := g.parameters.MigratedFromCount[rand.Intn(len(g.parameters.MigratedFromCount))]
		for i := 0; i < migratedFromCount; i++ {
			migratedFromName := g.nameGenerator.Generate(10)
			input.Jobs[j].MigratedFrom = append(input.Jobs[j].MigratedFrom, bftinput.MigratedFromConfig{Name: migratedFromName})
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
	input = g.parameterProvider.Get("template").Apply(input, previousInput)
	input = g.parameterProvider.Get("compilation").Apply(input, previousInput)
	input = g.parameterProvider.Get("update").Apply(input, previousInput)
	input = g.parameterProvider.Get("cloud_properties").Apply(input, previousInput)
	input = g.parameterProvider.Get("fixed_migrated_from").Apply(input, previousInput)
	input = g.parameterProvider.Get("lifecycle").Apply(input, previousInput)

	return input
}

func (g *inputGenerator) randomizeJobs(jobs []bftinput.Job) []bftinput.Job {
	numberOfJobs := g.parameters.NumberOfJobs[rand.Intn(len(g.parameters.NumberOfJobs))]
	jobsSize := len(jobs)
	if numberOfJobs > jobsSize {
		for i := 0; i < numberOfJobs-jobsSize; i++ {
			jobName := g.nameGenerator.Generate(g.parameters.NameLength[rand.Intn(len(g.parameters.NameLength))])
			jobs = append(jobs, bftinput.Job{
				Name: jobName,
			})
		}
	} else if numberOfJobs < jobsSize {
		for i := 0; i < jobsSize-numberOfJobs; i++ {
			jobIdxToRemove := rand.Intn(len(jobs))
			if jobIdxToRemove == len(jobs)-1 {
				jobs = jobs[:jobIdxToRemove]
			} else {
				jobs = append(jobs[:jobIdxToRemove], jobs[jobIdxToRemove+1:]...)
			}
		}
	}

	shuffledJobs := []bftinput.Job{}
	shuffledJobsIndeces := rand.Perm(numberOfJobs)
	for _, jobIndex := range shuffledJobsIndeces {
		shuffledJobs = append(shuffledJobs, jobs[jobIndex])
	}

	return shuffledJobs
}

func (g *inputGenerator) generateInputWithJobNames(jobNames []string) bftinput.Input {
	input := bftinput.Input{
		Jobs: []bftinput.Job{},
	}

	for _, jobName := range jobNames {
		input.Jobs = append(input.Jobs, bftinput.Job{
			Name:      jobName,
			Instances: g.parameters.Instances[rand.Intn(len(g.parameters.Instances))],
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
