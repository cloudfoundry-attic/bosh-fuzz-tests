package fakes

import (
	bftinput "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"
)

type FakeVariables struct {
}

func NewFakeVariables() *FakeVariables {
	return &FakeVariables{}
}

func (s *FakeVariables) Apply(input bftinput.Input, previousInput bftinput.Input) bftinput.Input {
	return input
}
