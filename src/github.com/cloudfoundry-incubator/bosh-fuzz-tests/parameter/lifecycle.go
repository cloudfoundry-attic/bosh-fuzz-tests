package parameter

import (
	"math/rand"

	bftinput "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"
)

type Lifecycle struct{}

func lifecycles(job, previousJob bftinput.Job) []string {
	lifecycles := []string{"service"}

	if job.PersistentDiskPool == "" && previousJob.PersistentDiskPool == "" &&
		job.PersistentDiskType == "" && previousJob.PersistentDiskType == "" &&
		job.PersistentDiskSize == 0 && previousJob.PersistentDiskSize == 0 {
		lifecycles = append(lifecycles, "errand")
	}

	return lifecycles
}

func NewLifecycle() Parameter {
	return Lifecycle{}
}

func (l Lifecycle) Apply(input, previousInput bftinput.Input) bftinput.Input {
	newInput := bftinput.Input{
		DirectorUUID: input.DirectorUUID,
		Jobs:         input.Jobs,
		Update:       input.Update,
		CloudConfig:  input.CloudConfig,
		Stemcells:    input.Stemcells,
		Variables:    input.Variables,
	}

	for i, job := range newInput.Jobs {
		previousJob := findJobFromInput(job, previousInput)

		cycles := lifecycles(job, previousJob)
		newInput.Jobs[i].Lifecycle = cycles[rand.Intn(len(cycles))]
	}

	return newInput
}

func findJobFromInput(desiredJob bftinput.Job, input bftinput.Input) bftinput.Job {
	for _, job := range input.Jobs {
		if job.Name == desiredJob.Name {
			return job
		} else {
			for _, migratedJob := range desiredJob.MigratedFrom {
				if job.Name == migratedJob.Name {
					return job
				}
			}
		}
	}

	return bftinput.Job{}
}
