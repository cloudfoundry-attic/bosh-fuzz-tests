package deployment

import (
	"fmt"
	"math/rand"

	"github.com/cloudfoundry-incubator/bosh-fuzz-tests/analyzer"
	bftinput "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"
	"github.com/cloudfoundry-incubator/bosh-load-tests/action/clirunner"
)

type ErrandStepGenerator struct{}

type ErrandStep struct {
	Name           string
	InstanceFilter string
	DeploymentName string
}

//go:generate counterfeiter . Step
type Step interface {
	Run(clirunner.Runner) error
}

//go:generate counterfeiter . Step
type StepGenerator interface {
	Steps(analyzer.Case) []Step
}

func NewErrandStepGenerator() ErrandStepGenerator {
	return ErrandStepGenerator{}
}

func (g ErrandStepGenerator) Steps(testCase analyzer.Case) []Step {
	steps := []Step{}

	jobs := testCase.Input.Jobs
	if len(jobs) == 0 {
		return steps
	}

	for i := 0; i < rand.Intn(6); i++ {
		job := jobs[rand.Intn(len(jobs))]

		if len(job.Templates) > 0 && job.Instances > 0 {
			instanceFilters := []string{
				"",
				job.Name,
				fmt.Sprintf("%s/0", job.Name),
			}

			steps = append(steps,
				ErrandStep{
					Name:           getErrandName(job),
					DeploymentName: "foo-deployment",
					InstanceFilter: instanceFilters[rand.Intn(len(instanceFilters))],
				},
			)
		}
	}

	return steps
}

func getErrandName(job bftinput.Job) string {
	possibilities := []string{job.Templates[rand.Intn(len(job.Templates))].Name}

	if job.Lifecycle == "errand" {
		possibilities = append(possibilities, job.Name)
	}

	return possibilities[rand.Intn(len(possibilities))]
}

func (es ErrandStep) Run(runner clirunner.Runner) error {
	args := []string{"run-errand", es.Name, "-d", es.DeploymentName}
	if es.InstanceFilter != "" {
		args = append(args, "--instance", es.InstanceFilter)
	}
	return runner.RunWithArgs(args...)
}
