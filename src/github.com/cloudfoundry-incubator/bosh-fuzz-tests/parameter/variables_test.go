package parameter_test

import (
	"fmt"

	bftdecider "github.com/cloudfoundry-incubator/bosh-fuzz-tests/decider/fakes"
	bftinput "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"
	bftnamegen "github.com/cloudfoundry-incubator/bosh-fuzz-tests/name_generator/fakes"
	. "github.com/cloudfoundry-incubator/bosh-fuzz-tests/parameter"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("variables", func() {
	var (
		variables     Parameter
		input         bftinput.Input
		previousInput bftinput.Input
		variableTypes []string                      = []string{"type1", "type2", "type3"}
		nameGenerator *bftnamegen.FakeNameGenerator = &bftnamegen.FakeNameGenerator{}
		decider       *bftdecider.FakeDeciderMulti  = &bftdecider.FakeDeciderMulti{}
	)

	BeforeEach(func() {
		input = bftinput.Input{}
		previousInput = bftinput.Input{}
	})

	Context("When number of variables given is 0", func() {
		BeforeEach(func() {
			variables = NewVariables(0, variableTypes, nameGenerator, decider)
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
			nameGenerator.Names = []string{"NewVar1", "NewVar2", "NewVar3", "NewVar4", "NewVar5"}
			variables = NewVariables(numVariables, variableTypes, nameGenerator, decider)
		})

		It(fmt.Sprintf("variables section contains %d length", numVariables), func() {
			output := variables.Apply(input, previousInput)
			Expect(output.Variables).To(HaveLen(numVariables))
		})
	})

	Context("when generating certificates", func() {
		var (
			numVariables      int = 2
			previousVariables []bftinput.Variable
		)

		BeforeEach(func() {
			nameGenerator.Names = []string{"NewVar1", "NewVar2", "NewVar3", "NewVar4", "NewVar5"}
			variables = NewVariables(numVariables, []string{"certificate"}, nameGenerator, decider)
		})

		Context("when previous variables are used", func() {
			It("returns the correct number of variables", func() {
				decider.YesResults = []bool{true}
				previousVariables = []bftinput.Variable{
					{
						Name:    "CertVarName",
						Type:    "certificate",
						Options: map[string]interface{}{"is_ca": true},
					},
				}

				input := variables.Apply(bftinput.Input{}, bftinput.Input{Variables: previousVariables})
				Expect(len(input.Variables)).To(Equal(numVariables))
			})

			Context("if a previous certificate's dependency isn't present", func() {
				It("does not get added", func() {
					decider.YesResults = []bool{false, true, true}
					previousVariables = []bftinput.Variable{
						{
							Name:    "cert1",
							Type:    "certificate",
							Options: map[string]interface{}{"is_ca": true},
						},
						{
							Name:    "cert2",
							Type:    "certificate",
							Options: map[string]interface{}{"is_ca": false, "ca": "cert1"},
						},
					}

					input := variables.Apply(bftinput.Input{}, bftinput.Input{Variables: previousVariables})

					Expect(len(input.Variables)).To(Equal(numVariables))
					for _, variable := range input.Variables {
						Expect(variable.Name).ToNot(Equal("cert2"))
					}
				})
			})
		})
		Context("when previous variables are NOT used", func() {
			BeforeEach(func() {
				previousVariables = []bftinput.Variable{}
				decider.YesResults = []bool{false}
			})

			It("returns the correct number of variables", func() {
				input := variables.Apply(bftinput.Input{}, bftinput.Input{Variables: previousVariables})
				Expect(len(input.Variables)).To(Equal(numVariables))
			})

			It("always returns atleast one CA cert", func() {
				input := variables.Apply(bftinput.Input{}, bftinput.Input{Variables: previousVariables})
				Expect(len(input.Variables)).To(Equal(numVariables))

				results := []bool{}
				for _, variable := range input.Variables {
					result, ok := variable.Options["is_ca"]
					if ok {
						results = append(results, result.(bool))
					}
				}

				Expect(len(results) > 0).To(BeTrue())
				Expect(results[0]).To(BeTrue())
			})
		})
	})
})
