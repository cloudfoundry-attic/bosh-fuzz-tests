package fakes

import (
	bftinput "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"
)

type FakeTemplate struct {
}

func NewFakeTemplate() *FakeTemplate {
	return &FakeTemplate{}
}

func (s *FakeTemplate) Apply(input bftinput.Input, previousInput bftinput.Input) bftinput.Input {
	for j, _ := range input.Jobs {
		input.Jobs[j].Templates = []bftinput.Template{
			{Name: "simple"},
		}
	}

	return input
}
