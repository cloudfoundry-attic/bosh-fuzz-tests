package expectation

import (
	"github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"
	"github.com/cloudfoundry-incubator/bosh-fuzz-tests/parser"
	bltclirunner "github.com/cloudfoundry-incubator/bosh-load-tests/action/clirunner"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

type variablesExpectation struct {
	variables []input.Variable
}

func NewVariablesExpectation(variables []input.Variable) Expectation {
	return &variablesExpectation{
		variables: variables,
	}
}

func (v *variablesExpectation) Run(cliRunner bltclirunner.Runner, taskId string) error {
	eventLog, err := cliRunner.RunWithOutput("events", "--task="+taskId, "--object-type=variable", "--json")
	if err != nil {
		return bosherr.WrapError(err, "Getting event logs")
	}

	events, err := parser.ParseEventLog(eventLog)
	if err != nil {
		return bosherr.WrapError(err, "Parsing event logs")
	}

	if len(v.variables) != len(events) {
		return bosherr.Errorf("Expected %d variables to be created but found %d", len(v.variables), len(events))
	}

	for _, variable := range v.variables {
		if _, err := events.FindById(variable.Name); err != nil {
			return bosherr.Errorf("Variable '%s' was not created", variable.Name)
		}
	}

	return nil
}
