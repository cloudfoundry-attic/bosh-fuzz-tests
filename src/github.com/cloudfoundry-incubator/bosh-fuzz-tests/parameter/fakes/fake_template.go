package fakes

import (
	bftinput "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"
)

type FakeJob struct {
}

func NewFakeJob() *FakeJob {
	return &FakeJob{}
}

func (s *FakeJob) Apply(input bftinput.Input, previousInput bftinput.Input) bftinput.Input {
	for j, _ := range input.InstanceGroups {
		input.InstanceGroups[j].Jobs = []bftinput.Job{
			{Name: "simple"},
		}
	}

	return input
}
