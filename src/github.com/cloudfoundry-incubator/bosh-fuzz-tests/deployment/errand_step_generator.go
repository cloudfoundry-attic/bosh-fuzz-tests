package deployment

import (
	"fmt"
	"math/rand"

	"github.com/cloudfoundry-incubator/bosh-fuzz-tests/analyzer"
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

		if len(job.Templates) > 0 {
			instanceFilters := []string{
				"",
				job.Name,
				fmt.Sprintf("%s/0", job.Name),
				fmt.Sprintf("%s/first", job.Name),
				fmt.Sprintf("%s/any", job.Name),
			}

			steps = append(steps,
				ErrandStep{
					Name:           job.Templates[rand.Intn(len(job.Templates))].Name,
					DeploymentName: "foo-deployment",
					InstanceFilter: instanceFilters[rand.Intn(len(instanceFilters))],
				},
			)
		}
	}

	return steps
}

func (es ErrandStep) Run(runner clirunner.Runner) error {
	args := []string{"run-errand", es.Name, "-d", es.DeploymentName}
	if es.InstanceFilter != "" {
		args = append(args, "--instance", es.InstanceFilter)
	}
	return runner.RunWithArgs(args...)
}
