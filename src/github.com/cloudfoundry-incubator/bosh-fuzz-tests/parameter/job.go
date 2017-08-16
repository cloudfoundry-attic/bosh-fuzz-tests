package parameter

import (
	"math/rand"

	bftinput "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"
)

type job struct {
	jobs [][]string
}

func NewJob(jobs [][]string) Parameter {
	return &job{
		jobs: jobs,
	}
}

func (t *job) Apply(input bftinput.Input, previousInput bftinput.Input) bftinput.Input {
	for j, _ := range input.InstanceGroups {
		pickedJobs := t.jobs[rand.Intn(len(t.jobs))]
		input.InstanceGroups[j].Jobs = []bftinput.Job{}

		for _, pickedJobName := range pickedJobs {
			input.InstanceGroups[j].Jobs = append(input.InstanceGroups[j].Jobs, bftinput.Job{
				Name: pickedJobName,
			})
		}
	}

	return input
}
