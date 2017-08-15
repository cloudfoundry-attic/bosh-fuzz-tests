package parameter

import (
	"math/rand"

	bftinput "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"
)

type Lifecycle struct{}

func lifecycles(job bftinput.Job) []string {
	lifecycles := []string{"service"}

	if job.PersistentDiskPool == "" && job.PersistentDiskSize == 0 && job.PersistentDiskType == "" {
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
		cycles := lifecycles(job)
		newInput.Jobs[i].Lifecycle = cycles[rand.Intn(len(cycles))]
	}

	return newInput
}
