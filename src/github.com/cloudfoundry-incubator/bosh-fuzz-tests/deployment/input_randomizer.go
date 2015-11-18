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
		input := Input{
			Jobs: []Job{},
		}

		numberOfJobs := ir.parameters.NumberOfJobs[rand.Intn(len(ir.parameters.NumberOfJobs))]
		azs := map[string]bool{}
		persistentDiskDefinition := ir.parameters.PersistentDiskDefinition[rand.Intn(len(ir.parameters.PersistentDiskDefinition))]

		for jobNumber := 0; jobNumber < numberOfJobs; jobNumber++ {
			var jobName string
			if i > 0 && len(inputs[i-1].Jobs) > jobNumber {
				jobName = inputs[i-1].Jobs[jobNumber].Name
			}

			if jobName == "" {
				jobName = nameGenerator.Generate(ir.parameters.NameLength[rand.Intn(len(ir.parameters.NameLength))])
			}

			job := Job{
				Name:              jobName,
				Instances:         ir.parameters.Instances[rand.Intn(len(ir.parameters.Instances))],
				AvailabilityZones: ir.parameters.AvailabilityZones[rand.Intn(len(ir.parameters.AvailabilityZones))],
			}

			persistentDiskSize := ir.parameters.PersistentDiskSize[rand.Intn(len(ir.parameters.PersistentDiskSize))]

			if persistentDiskDefinition == "disk_pool" {
				job.PersistentDiskPool = nameGenerator.Generate(10)
				input.CloudConfig.PersistentDiskPools = append(
					input.CloudConfig.PersistentDiskPools,
					DiskConfig{Name: job.PersistentDiskPool, Size: persistentDiskSize},
				)
			} else if persistentDiskDefinition == "disk_type" {
				job.PersistentDiskType = nameGenerator.Generate(10)
				input.CloudConfig.PersistentDiskTypes = append(
					input.CloudConfig.PersistentDiskTypes,
					DiskConfig{Name: job.PersistentDiskType, Size: persistentDiskSize},
				)
			} else {
				job.PersistentDiskSize = persistentDiskSize
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

			input.Jobs = append(input.Jobs, job)
		}

		inputs = append(inputs, input)
	}
	return inputs, nil
}
