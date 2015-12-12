package fakes

import (
	bftinput "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"
)

type FakeCompilation struct {
}

func NewFakeCompilation() *FakeCompilation {
	return &FakeCompilation{}
}

func (s *FakeCompilation) Apply(input bftinput.Input) bftinput.Input {
	input.CloudConfig.NumberOfCompilationWorkers = 3
	return input
}
