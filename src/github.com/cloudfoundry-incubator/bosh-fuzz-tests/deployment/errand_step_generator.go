package deployment

import (
	"math/rand"

	"github.com/cloudfoundry-incubator/bosh-fuzz-tests/analyzer"
	"github.com/cloudfoundry-incubator/bosh-load-tests/action/clirunner"
)

type ErrandStepGenerator struct{}

type ErrandStep struct {
	Name           string
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
	jobs := testCase.Input.Jobs

	if len(jobs) > 0 {
		job := jobs[rand.Intn(len(jobs))]

		if len(job.Templates) > 0 {
			return []Step{
				ErrandStep{
					Name:           job.Templates[rand.Intn(len(job.Templates))].Name,
					DeploymentName: "foo-deployment",
				},
			}
		}
	}

	return []Step{}
}

func (es ErrandStep) Run(runner clirunner.Runner) error {
	return runner.RunWithArgs("run-errand", es.Name, "-d", es.DeploymentName)
}
