package deployment

import (
	"math/rand"

	bftconfig "github.com/cloudfoundry-incubator/bosh-fuzz-tests/config"
	bftinput "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
)

type JobsRandomizer interface {
	Generate() ([]bftinput.Input, error)
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

func (ir *jobsRandomizer) Generate() ([]bftinput.Input, error) {
	inputs := []bftinput.Input{}

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

func (ir *jobsRandomizer) generateInput(jobNames []string, migratingDeployment bool) bftinput.Input {
	input := &bftinput.Input{
		Jobs: []bftinput.Job{},
	}

	azs := map[string]bool{}
	persistentDiskDefinition := ir.parameters.PersistentDiskDefinition[rand.Intn(len(ir.parameters.PersistentDiskDefinition))]
	vmTypeDefinition := ir.parameters.VmTypeDefinition[rand.Intn(len(ir.parameters.VmTypeDefinition))]
	stemcellDefinition := ir.parameters.StemcellDefinition[rand.Intn(len(ir.parameters.StemcellDefinition))]

	for _, jobName := range jobNames {
		job := &bftinput.Job{
			Name:              jobName,
			Instances:         ir.parameters.Instances[rand.Intn(len(ir.parameters.Instances))],
			AvailabilityZones: ir.parameters.AvailabilityZones[rand.Intn(len(ir.parameters.AvailabilityZones))],
		}

		ir.assignPersistentDisk(persistentDiskDefinition, job, input)
		ir.assignVmType(vmTypeDefinition, stemcellDefinition, job, input)

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
				job.MigratedFrom = append(job.MigratedFrom, bftinput.MigratedFromConfig{Name: migratedFromName})
			}
		}

		input.Jobs = append(input.Jobs, *job)
	}

	return *input
}

func (ir *jobsRandomizer) generateJobNames(i int, inputs []bftinput.Input) []string {
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

func (ir *jobsRandomizer) specifyAzIfMigratingJobDoesNotHaveAz(migratingInput bftinput.Input, currentInput bftinput.Input) {
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

func (ir *jobsRandomizer) assignPersistentDisk(persistentDiskDefinition string, job *bftinput.Job, input *bftinput.Input) {
	persistentDiskSize := ir.parameters.PersistentDiskSize[rand.Intn(len(ir.parameters.PersistentDiskSize))]
	if persistentDiskSize != 0 {
		if persistentDiskDefinition == "disk_pool" {
			job.PersistentDiskPool = ir.nameGenerator.Generate(10)
			input.CloudConfig.PersistentDiskPools = append(
				input.CloudConfig.PersistentDiskPools,
				bftinput.DiskConfig{Name: job.PersistentDiskPool, Size: persistentDiskSize},
			)
		} else if persistentDiskDefinition == "disk_type" {
			job.PersistentDiskType = ir.nameGenerator.Generate(10)
			input.CloudConfig.PersistentDiskTypes = append(
				input.CloudConfig.PersistentDiskTypes,
				bftinput.DiskConfig{Name: job.PersistentDiskType, Size: persistentDiskSize},
			)
		} else {
			job.PersistentDiskSize = persistentDiskSize
		}
	}
}

func (ir *jobsRandomizer) assignVmType(vmTypeDefinition string, stemcellDefinition string, job *bftinput.Job, input *bftinput.Input) {
	var stemcellConfig bftinput.StemcellConfig
	if stemcellDefinition == "os_version" {
		stemcellConfig = bftinput.StemcellConfig{
			OS:      "toronto-os",
			Version: "1",
		}
	} else {
		stemcellConfig = bftinput.StemcellConfig{
			Name:    "ubuntu-stemcell",
			Version: "1",
		}
	}

	if vmTypeDefinition == "vm_type" {
		job.VmType = ir.nameGenerator.Generate(10)
		input.CloudConfig.VmTypes = append(
			input.CloudConfig.VmTypes,
			bftinput.VmTypeConfig{Name: job.VmType},
		)
		stemcellConfig.Alias = "default"
		input.Stemcells = []bftinput.StemcellConfig{stemcellConfig}
	} else if vmTypeDefinition == "resource_pool" {
		job.ResourcePool = ir.nameGenerator.Generate(10)
		input.CloudConfig.ResourcePools = append(
			input.CloudConfig.ResourcePools,
			bftinput.ResourcePoolConfig{
				Name:     job.ResourcePool,
				Stemcell: stemcellConfig,
			},
		)
	}
}
