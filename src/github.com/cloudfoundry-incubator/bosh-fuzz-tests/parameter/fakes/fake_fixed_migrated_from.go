package fakes

import (
	bftinput "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"
)

type FakeFixedMigratedFrom struct {
}

func NewFakeFixedMigratedFrom() *FakeFixedMigratedFrom {
	return &FakeFixedMigratedFrom{}
}

func (s *FakeFixedMigratedFrom) Apply(input bftinput.Input, previousInput bftinput.Input) bftinput.Input {
	return input
}
