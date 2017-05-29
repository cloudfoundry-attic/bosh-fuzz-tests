package parameter

import (
	"math/rand"

	bftdecider "github.com/cloudfoundry-incubator/bosh-fuzz-tests/decider"
	bftinput "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"
	"github.com/cloudfoundry-incubator/bosh-fuzz-tests/name_generator"
)

type variables struct {
	numVariables  int
	variableTypes []string
	nameGenerator name_generator.NameGenerator
	decider       bftdecider.Decider
}

func NewVariables(numVariables int, variableTypes []string, nameGenerator name_generator.NameGenerator, decider bftdecider.Decider) Parameter {
	return &variables{
		numVariables:  numVariables,
		variableTypes: variableTypes,
		nameGenerator: nameGenerator,
		decider:       decider,
	}
}

func (v *variables) Apply(input bftinput.Input, previousInput bftinput.Input) bftinput.Input {
	variables := []bftinput.Variable{}
	previousVariables := previousInput.Variables
	certsHash := map[string]bool{}

	for _, variable := range previousVariables {
		if v.decider.IsYes() {
			if variable.Type == "certificate" {
				caName, caExists := variable.Options["ca"]
				if caExists {
					if _, ok := certsHash[caName.(string)]; !ok {
						// skip previous certificate unless the dependencies exist in the current chain
						continue
					}
				}
				certsHash[variable.Name] = true
			}
			variables = append(variables, variable)
		}
	}

	missingVariables := v.numVariables - len(variables)
	for i := 0; i < missingVariables; i++ {
		variableType := v.variableTypes[rand.Intn(len(v.variableTypes))]
		variable := bftinput.Variable{
			Name: v.nameGenerator.Generate(20),
			Type: variableType,
		}
		switch variableType {
		case "certificate":
			certificateVariable := v.generateCertificateVariable(variables, variable)
			variables = append(variables, certificateVariable)
		default:
			variables = append(variables, variable)
		}
	}

	input.Variables = variables[:v.numVariables]

	return input
}

func (v *variables) generateCertificateVariable(existingVariables []bftinput.Variable, variable bftinput.Variable) bftinput.Variable {

	variable.Options = make(map[string]interface{})
	variable.Options["common_name"] = v.nameGenerator.Generate(10)
	variable.Options["is_ca"] = v.decider.IsYes()

	caNames := []string{}
	for _, existingVariable := range existingVariables {
		if existingVariable.Type == "certificate" {
			result, ok := existingVariable.Options["is_ca"]
			useVariable := ok && result.(bool)
			if useVariable {
				caNames = append(caNames, existingVariable.Name)
			}
		}
	}

	if len(caNames) == 0 {
		// Force first cert to be the best ever CA
		variable.Options["is_ca"] = true
	} else {
		selectedCA := rand.Intn(len(caNames))
		variable.Options["ca"] = caNames[selectedCA]
	}

	return variable
}
