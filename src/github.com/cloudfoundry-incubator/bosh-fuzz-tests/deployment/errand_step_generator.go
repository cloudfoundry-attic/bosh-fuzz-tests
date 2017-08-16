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

	instanceGroups := testCase.Input.InstanceGroups
	if len(instanceGroups) == 0 {
		return steps
	}

	for i := 0; i < rand.Intn(6); i++ {
		instanceGroup := instanceGroups[rand.Intn(len(instanceGroups))]

		if len(instanceGroup.Jobs) > 0 && instanceGroup.Instances > 0 {
			instanceFilters := []string{
				"",
				instanceGroup.Name,
				fmt.Sprintf("%s/0", instanceGroup.Name),
			}

			steps = append(steps,
				ErrandStep{
					Name:           getErrandName(instanceGroup),
					DeploymentName: "foo-deployment",
					InstanceFilter: instanceFilters[rand.Intn(len(instanceFilters))],
				},
			)
		}
	}

	return steps
}

func getErrandName(instanceGroup bftinput.InstanceGroup) string {
	possibilities := []string{instanceGroup.Jobs[rand.Intn(len(instanceGroup.Jobs))].Name}

	if instanceGroup.Lifecycle == "errand" {
		possibilities = append(possibilities, instanceGroup.Name)
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
