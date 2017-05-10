package parameter

import (
	"math/rand"

	bftinput "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"
	"github.com/cloudfoundry-incubator/bosh-fuzz-tests/name_generator"
)

type variables struct {
	numVariables  int
	variableTypes []string
	nameGenerator name_generator.NameGenerator
}

func NewVariables(numVariables int, variableTypes []string) Parameter {
	return &variables{
		numVariables:  numVariables,
		variableTypes: variableTypes,
		nameGenerator: name_generator.NewNameGenerator(),
	}
}

func (s *variables) Apply(input bftinput.Input, previousInput bftinput.Input) bftinput.Input {
	variables := []bftinput.Variable{}

	for i := 0; i < s.numVariables; i++ {
		variableType := s.variableTypes[rand.Intn(len(s.variableTypes))]
		variables = append(variables, bftinput.Variable{
			Name: s.nameGenerator.Generate(20),
			Type: variableType,
		})
	}

	input.Variables = variables

	return input
}
