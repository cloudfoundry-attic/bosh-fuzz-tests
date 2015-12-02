package deployment

import (
	"math/rand"

	bftconfig "github.com/cloudfoundry-incubator/bosh-fuzz-tests/config"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
)

type JobsRandomizer interface {
	Generate() ([]Input, error)
}

type jobsRandomizer struct {
	parameters                bftconfig.Parameters
	numberOfConsequentDeploys int
	nameGenerator             NameGenerator
	logger                    boshlog.Logger
}

func NewJobsRandomizer(parameters bftconfig.Parameters, numberOfConsequentDeploys int, nameGenerator NameGenerator, logger boshlog.Logger) JobsRandomizer {
	return &jobsRandomizer{
		parameters:                parameters,
		numberOfConsequentDeploys: numberOfConsequentDeploys,
		nameGenerator:             nameGenerator,
		logger:                    logger,
	}
}

func (ir *jobsRandomizer) Generate() ([]Input, error) {
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

func (ir *jobsRandomizer) generateInput(jobNames []string, migratingDeployment bool) Input {
	input := &Input{
		Jobs: []Job{},
	}

	azs := map[string]bool{}
	persistentDiskDefinition := ir.parameters.PersistentDiskDefinition[rand.Intn(len(ir.parameters.PersistentDiskDefinition))]
	vmTypeDefinition := ir.parameters.VmTypeDefinition[rand.Intn(len(ir.parameters.VmTypeDefinition))]

	for _, jobName := range jobNames {
		job := &Job{
			Name:              jobName,
			Instances:         ir.parameters.Instances[rand.Intn(len(ir.parameters.Instances))],
			AvailabilityZones: ir.parameters.AvailabilityZones[rand.Intn(len(ir.parameters.AvailabilityZones))],
		}

		ir.assignPersistentDisk(persistentDiskDefinition, job, input)
		ir.assignVmType(vmTypeDefinition, job, input)

		for _, az := range job.AvailabilityZones {
			if azs[az] != true {
				input.CloudConfig.AvailabilityZones = append(input.CloudConfig.AvailabilityZones, az)
			}
			azs[az] = true
		}

		if !migratingDeployment {
			migratedFromCount := ir.parameters.MigratedFromCount[rand.Intn(len(ir.parameters.MigratedFromCount))]
			for i := 0; i < migratedFromCount; i++ {
				migratedFromName := ir.nameGenerator.Generate(10)
				job.MigratedFrom = append(job.MigratedFrom, MigratedFromConfig{Name: migratedFromName})
			}
		}

		input.Jobs = append(input.Jobs, *job)
	}

	return *input
}

func (ir *jobsRandomizer) generateJobNames(i int, inputs []Input) []string {
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

func (ir *jobsRandomizer) specifyAzIfMigratingJobDoesNotHaveAz(migratingInput Input, currentInput Input) {
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

func (ir *jobsRandomizer) assignPersistentDisk(persistentDiskDefinition string, job *Job, input *Input) {
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
}

func (ir *jobsRandomizer) assignVmType(vmTypeDefinition string, job *Job, input *Input) {
	if vmTypeDefinition == "vm_type" {
		job.VmType = ir.nameGenerator.Generate(10)
		input.CloudConfig.VmTypes = append(
			input.CloudConfig.VmTypes,
			VmTypeConfig{Name: job.VmType},
		)
	} else if vmTypeDefinition == "resource_pool" {
		job.ResourcePool = ir.nameGenerator.Generate(10)
		input.CloudConfig.ResourcePools = append(
			input.CloudConfig.ResourcePools,
			VmTypeConfig{Name: job.ResourcePool},
		)
	}
}
