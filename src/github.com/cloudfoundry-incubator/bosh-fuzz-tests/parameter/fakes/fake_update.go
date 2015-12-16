package fakes

import (
	bftinput "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"
)

type FakeUpdate struct {
}

func NewFakeUpdate() *FakeUpdate {
	return &FakeUpdate{}
}

func (u *FakeUpdate) Apply(input bftinput.Input, previousInput bftinput.Input) bftinput.Input {
	input.Update.Canaries = 3
	input.Update.MaxInFlight = 5
	input.Update.Serial = "true"
	return input
}
