package analyzer

import (
	bftexpectation "github.com/cloudfoundry-incubator/bosh-fuzz-tests/expectation"
	bftinput "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"
)

type Analyzer interface {
	Analyze(inputs []bftinput.Input) []Case
}

type Case struct {
	Input        bftinput.Input
	Expectations []bftexpectation.Expectation
}

type analyzer struct {
	expectationFactory bftexpectation.Factory
	stemcellComparator Comparator
}

func NewAnalyzer(expectationFactory bftexpectation.Factory) Analyzer {
	return &analyzer{
		expectationFactory: expectationFactory,
		stemcellComparator: NewStemcellComparator(expectationFactory),
	}
}

func (a *analyzer) Analyze(inputs []bftinput.Input) []Case {
	cases := []Case{}
	for i := range inputs {
		expectations := []bftexpectation.Expectation{}

		if i != 0 {
			expectations = append(expectations, a.stemcellComparator.Compare(inputs[i-1], inputs[i])...)
		}

		cases = append(cases, Case{
			Input:        inputs[i],
			Expectations: expectations,
		})
	}

	return cases
}
