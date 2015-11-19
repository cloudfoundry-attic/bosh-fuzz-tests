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
	nameGenerator             NameGenerator
	logger                    boshlog.Logger
}

func NewSeededInputRandomizer(parameters bftconfig.Parameters, numberOfConsequentDeploys int, seed int64, logger boshlog.Logger) InputRandomizer {
	return &inputRandomizer{
		parameters:                parameters,
		numberOfConsequentDeploys: numberOfConsequentDeploys,
		seed:          seed,
		nameGenerator: NewNameGenerator(),
		logger:        logger,
	}
}

func NewInputRandomizer(parameters bftconfig.Parameters, numberOfConsequentDeploys int, logger boshlog.Logger) InputRandomizer {
	return &inputRandomizer{
		parameters:                parameters,
		numberOfConsequentDeploys: numberOfConsequentDeploys,
		seed:          time.Now().Unix(),
		nameGenerator: NewNameGenerator(),
		logger:        logger,
	}
}

func (ir *inputRandomizer) Generate() ([]Input, error) {
	ir.logger.Info("inputRandomizer", "Seeding with %d", ir.seed)

	rand.Seed(ir.seed)

	inputs := []Input{}

	for i := 0; i < ir.numberOfConsequentDeploys; i++ {
		jobNames := ir.generateJobNames(i, inputs)
		input := ir.generateInput(jobNames, false)

		migratingJobs := []string{}
		for _, j := range input.Jobs {
			for _, m := range j.MigratedFrom {
				migratingJobs = append(migratingJobs, m.Name)
			}
		}

		if len(migratingJobs) > 0 {
			migratingInput := ir.generateInput(migratingJobs, true)

			ir.specifyAzIfMigratingJobDoesNotHaveAz(migratingInput, input)

			inputs = append(inputs, migratingInput)
		}

		inputs = append(inputs, input)
	}

	return inputs, nil
}

func (ir *inputRandomizer) generateInput(jobNames []string, migratingDeployment bool) Input {
	input := Input{
		Jobs: []Job{},
	}

	azs := map[string]bool{}
	persistentDiskDefinition := ir.parameters.PersistentDiskDefinition[rand.Intn(len(ir.parameters.PersistentDiskDefinition))]

	for _, jobName := range jobNames {
		job := Job{
			Name:              jobName,
			Instances:         ir.parameters.Instances[rand.Intn(len(ir.parameters.Instances))],
			AvailabilityZones: ir.parameters.AvailabilityZones[rand.Intn(len(ir.parameters.AvailabilityZones))],
		}

		// Workaround for this bug #108499370, migrating job and destination job cannot have nil az
		if migratingDeployment {
			for job.AvailabilityZones == nil {
				job.AvailabilityZones = ir.parameters.AvailabilityZones[rand.Intn(len(ir.parameters.AvailabilityZones))]
			}
		}

		persistentDiskSize := ir.parameters.PersistentDiskSize[rand.Intn(len(ir.parameters.PersistentDiskSize))]
		if persistentDiskSize != 0 {
			if persistentDiskDefinition == "disk_pool" {
				job.PersistentDiskPool = ir.nameGenerator.Generate(10)
				input.CloudConfig.PersistentDiskPools = append(
					input.CloudConfig.PersistentDiskPools,
					DiskConfig{Name: job.PersistentDiskPool, Size: persistentDiskSize},
				)
			} else if persistentDiskDefinition == "disk_type" {
				job.PersistentDiskType = ir.nameGenerator.Generate(10)
				input.CloudConfig.PersistentDiskTypes = append(
					input.CloudConfig.PersistentDiskTypes,
					DiskConfig{Name: job.PersistentDiskType, Size: persistentDiskSize},
				)
			} else {
				job.PersistentDiskSize = persistentDiskSize
			}
		}

		for _, az := range job.AvailabilityZones {
			if azs[az] != true {
				input.CloudConfig.AvailabilityZones = append(input.CloudConfig.AvailabilityZones, az)
			}
			azs[az] = true
		}

		if job.AvailabilityZones == nil {
			job.Network = "no-az"
		} else {
			job.Network = "default"
		}

		if !migratingDeployment {
			migratedFromCount := ir.parameters.MigratedFromCount[rand.Intn(len(ir.parameters.MigratedFromCount))]
			for i := 0; i < migratedFromCount; i++ {
				migratedFromName := ir.nameGenerator.Generate(10)
				job.MigratedFrom = append(job.MigratedFrom, MigratedFromConfig{Name: migratedFromName})
			}
		}

		input.Jobs = append(input.Jobs, job)
	}

	return input
}

func (ir *inputRandomizer) generateJobNames(i int, inputs []Input) []string {
	numberOfJobs := ir.parameters.NumberOfJobs[rand.Intn(len(ir.parameters.NumberOfJobs))]
	jobNames := []string{}

	for j := 0; j < numberOfJobs; j++ {
		var jobName string
		if i > 0 && len(inputs[i-1].Jobs) > j {
			jobName = inputs[i-1].Jobs[j].Name
		}

		if jobName == "" {
			jobName = ir.nameGenerator.Generate(ir.parameters.NameLength[rand.Intn(len(ir.parameters.NameLength))])
		}

		jobNames = append(jobNames, jobName)
	}

	return jobNames
}

func (ir *inputRandomizer) specifyAzIfMigratingJobDoesNotHaveAz(migratingInput Input, currentInput Input) {
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
