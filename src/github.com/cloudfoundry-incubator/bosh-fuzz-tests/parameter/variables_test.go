package parameter_test

import (
	"fmt"

	bftinput "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"
	. "github.com/cloudfoundry-incubator/bosh-fuzz-tests/parameter"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("variables", func() {
	var (
		variables     Parameter
		input         bftinput.Input
		previousInput bftinput.Input
		variableTypes []string = []string{"type1", "type2", "type3"}
	)

	BeforeEach(func() {
		input = bftinput.Input{}
		previousInput = bftinput.Input{}
	})

	Context("When number of variables given is 0", func() {
		BeforeEach(func() {
			variables = NewVariables(0, variableTypes)
		})

		It("sets the input to be an empty array", func() {
			output := variables.Apply(input, previousInput)
			Expect(output.Variables).To(BeEmpty())
		})
	})

	Context("When number of variables given is > 0", func() {
		var (
			numVariables int = 5
		)

		BeforeEach(func() {
			variables = NewVariables(numVariables, variableTypes)
		})

		It(fmt.Sprintf("variables section contains %d length", numVariables), func() {
			output := variables.Apply(input, previousInput)
			Expect(output.Variables).To(HaveLen(numVariables))
		})
	})
})
