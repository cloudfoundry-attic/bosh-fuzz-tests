package fakes

import (
	bftinput "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"
)

type FakeLifecycle struct {
}

func NewFakeLifecycle() *FakeLifecycle {
	return &FakeLifecycle{}
}

func (s *FakeLifecycle) Apply(input bftinput.Input, previousInput bftinput.Input) bftinput.Input {
	for i := range input.Jobs {
		input.Jobs[i].Lifecycle = "mufasa"
	}
	return input
}
