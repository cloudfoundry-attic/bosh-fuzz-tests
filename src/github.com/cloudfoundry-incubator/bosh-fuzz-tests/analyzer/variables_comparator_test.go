package analyzer_test

import (
	. "github.com/cloudfoundry-incubator/bosh-fuzz-tests/analyzer"
	bftinput "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("VariablesComparator", func() {
	var (
		variablesComparator Comparator
		previousInputs      []bftinput.Input
		currentInput        bftinput.Input
	)

	BeforeEach(func() {
		variablesComparator = NewVariablesComparator()
	})

	Context("when calling compare", func() {
		It("returns a sincele Expectation", func() {
			expectations := variablesComparator.Compare(previousInputs, currentInput)
			Expect(expectations).To(HaveLen(1))
		})
	})

})
